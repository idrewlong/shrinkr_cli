package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/idrewlong/shrinkr_cli/internal/compressor"
	"github.com/idrewlong/shrinkr_cli/internal/scanner"
	"github.com/idrewlong/shrinkr_cli/internal/ui"
	"github.com/idrewlong/shrinkr_cli/internal/worker"
	"github.com/spf13/cobra"
)

// ShrinkConfig holds all parameters needed to run a compression job.
// Both the flag-based CLI and the interactive wizard populate this struct.
type ShrinkConfig struct {
	InputFolder string
	Output      string
	Format      string
	Size        int
	Quality     int
	MinQuality  int
	MaxQuality  int
	Workers     int
	Recursive   bool
}

var (
	output     string
	size       int
	format     string
	quality    int
	workers    int
	recursive  bool
	minQuality int
	maxQuality int
)

var rootCmd = &cobra.Command{
	Use:   "shrinkr [input-folder]",
	Short: "Shrinkr — Fast image compression & format conversion",
	Long: `Shrinkr is a blazing-fast CLI tool for compressing images and converting
them to modern formats. It uses concurrent processing to handle batches
of images quickly, with smart quality optimization to hit target file sizes.

Run with no arguments to launch the interactive wizard.

Supported input formats: JPG, PNG, TIFF, WebP, GIF, AVIF, HEIF/HEIC
Supported output formats: WebP, AVIF, JPEG, PNG`,
	Args: cobra.RangeArgs(0, 1),
	RunE: run,
}

func init() {
	rootCmd.Flags().StringVarP(&output, "output", "o", "compressed", "Output folder for compressed images")
	rootCmd.Flags().IntVarP(&size, "size", "s", 500, "Target file size in KB")
	rootCmd.Flags().StringVarP(&format, "format", "f", "webp", "Output format: webp, avif, jpeg, png")
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 85, "Initial quality (1-100)")
	rootCmd.Flags().IntVarP(&workers, "workers", "w", runtime.NumCPU(), "Number of concurrent workers")
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Process images in subfolders")
	rootCmd.Flags().IntVar(&minQuality, "min-quality", 60, "Minimum quality threshold")
	rootCmd.Flags().IntVar(&maxQuality, "max-quality", 90, "Maximum quality threshold")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return runWizard()
	}

	ui.PrintLogo()

	config := ShrinkConfig{
		InputFolder: args[0],
		Output:      output,
		Format:      format,
		Size:        size,
		Quality:     quality,
		MinQuality:  minQuality,
		MaxQuality:  maxQuality,
		Workers:     workers,
		Recursive:   recursive,
	}

	return execute(config)
}

func execute(cfg ShrinkConfig) error {
	// Validate input folder
	info, err := os.Stat(cfg.InputFolder)
	if err != nil {
		return fmt.Errorf("cannot access folder %s: %v", cfg.InputFolder, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", cfg.InputFolder)
	}

	// Parse and validate output format
	outputFormat, err := compressor.ParseFormat(cfg.Format)
	if err != nil {
		return err
	}

	// Validate options
	if cfg.Size <= 0 {
		return fmt.Errorf("target size must be a positive number")
	}
	if cfg.Quality < 1 || cfg.Quality > 100 {
		return fmt.Errorf("quality must be between 1 and 100")
	}
	if cfg.MinQuality < 1 || cfg.MinQuality > 100 {
		return fmt.Errorf("min-quality must be between 1 and 100")
	}
	if cfg.MaxQuality < 1 || cfg.MaxQuality > 100 {
		return fmt.Errorf("max-quality must be between 1 and 100")
	}
	if cfg.MinQuality > cfg.MaxQuality {
		return fmt.Errorf("min-quality cannot be greater than max-quality")
	}
	if cfg.Workers < 1 {
		return fmt.Errorf("workers must be at least 1")
	}

	// Find all image files
	files, err := scanner.FindImages(cfg.InputFolder, cfg.Recursive)
	if err != nil {
		return fmt.Errorf("error scanning folder: %v", err)
	}

	if len(files) == 0 {
		fmt.Println(ui.WarningStyle.Render("  No supported image files found."))
		fmt.Println(ui.DimStyle.Render("  Supported: JPG, PNG, TIFF, WebP, GIF, AVIF, HEIF/HEIC"))
		return nil
	}

	// Print config summary
	fmt.Printf("  %s %s\n",
		ui.LabelStyle.Render("Found"),
		ui.TitleStyle.Render(fmt.Sprintf("%d image(s)", len(files))))
	fmt.Printf("  %s %d KB  %s  %s %s  %s  %s %d\n\n",
		ui.LabelStyle.Render("Target:"),
		cfg.Size,
		ui.DimStyle.Render("|"),
		ui.LabelStyle.Render("Format:"),
		ui.AccentStyle.Render(string(outputFormat)),
		ui.DimStyle.Render("|"),
		ui.LabelStyle.Render("Workers:"),
		cfg.Workers)

	// Create output folder
	if err := os.MkdirAll(cfg.Output, 0755); err != nil {
		return fmt.Errorf("cannot create output folder: %v", err)
	}

	// Start timing
	startTime := time.Now()

	// Create and start worker pool
	pool := worker.NewPool(cfg.Workers, len(files))
	pool.Start()

	// Submit all jobs in a goroutine
	go func() {
		for _, file := range files {
			relPath, _ := filepath.Rel(cfg.InputFolder, file)
			dir := filepath.Dir(relPath)
			baseName := filepath.Base(relPath)
			ext := filepath.Ext(baseName)
			nameNoExt := baseName[:len(baseName)-len(ext)]
			outputPath := filepath.Join(cfg.Output, dir, nameNoExt+outputFormat.FileExtension())

			pool.Submit(compressor.Job{
				InputPath:      file,
				OutputPath:     outputPath,
				TargetSizeKB:   cfg.Size,
				Format:         outputFormat,
				MinQuality:     cfg.MinQuality,
				MaxQuality:     cfg.MaxQuality,
				InitialQuality: cfg.Quality,
			})
		}
		pool.Done()
	}()

	// Collect all results with a live progress bar
	progressBar := ui.NewProgressBar(len(files))
	progressBar.Start()
	var results []compressor.Result
	for result := range pool.Results() {
		results = append(results, result)
		progressBar.Increment()
	}
	progressBar.Finish()

	// Print all file results after progress completes
	for _, result := range results {
		ui.PrintFileResult(result, cfg.Size)
	}
	fmt.Println()

	// Print summary
	elapsed := time.Since(startTime)
	ui.PrintSummary(results, elapsed, cfg.Size)

	return nil
}
