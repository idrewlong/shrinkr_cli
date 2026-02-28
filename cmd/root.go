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
	Use:   "shrinkr <input-folder>",
	Short: "Shrinkr — Fast image compression & format conversion",
	Long: `Shrinkr is a blazing-fast CLI tool for compressing images and converting
them to modern formats. It uses concurrent processing to handle batches
of images quickly, with smart quality optimization to hit target file sizes.

Supported input formats: JPG, PNG, TIFF, WebP, GIF, AVIF, HEIF/HEIC
Supported output formats: WebP, AVIF, JPEG, PNG`,
	Args: cobra.ExactArgs(1),
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
	inputFolder := args[0]

	// Print the ASCII logo
	ui.PrintLogo()

	// Validate input folder
	info, err := os.Stat(inputFolder)
	if err != nil {
		return fmt.Errorf("cannot access folder %s: %v", inputFolder, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", inputFolder)
	}

	// Parse and validate output format
	outputFormat, err := compressor.ParseFormat(format)
	if err != nil {
		return err
	}

	// Validate options
	if size <= 0 {
		return fmt.Errorf("target size must be a positive number")
	}
	if quality < 1 || quality > 100 {
		return fmt.Errorf("quality must be between 1 and 100")
	}
	if minQuality < 1 || minQuality > 100 {
		return fmt.Errorf("min-quality must be between 1 and 100")
	}
	if maxQuality < 1 || maxQuality > 100 {
		return fmt.Errorf("max-quality must be between 1 and 100")
	}
	if minQuality > maxQuality {
		return fmt.Errorf("min-quality cannot be greater than max-quality")
	}
	if workers < 1 {
		return fmt.Errorf("workers must be at least 1")
	}

	// Find all image files
	files, err := scanner.FindImages(inputFolder, recursive)
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
		size,
		ui.DimStyle.Render("|"),
		ui.LabelStyle.Render("Format:"),
		ui.AccentStyle.Render(string(outputFormat)),
		ui.DimStyle.Render("|"),
		ui.LabelStyle.Render("Workers:"),
		workers)

	// Create output folder
	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("cannot create output folder: %v", err)
	}

	// Start timing
	startTime := time.Now()

	// Create and start worker pool
	pool := worker.NewPool(workers, len(files))
	pool.Start()

	// Submit all jobs in a goroutine
	go func() {
		for _, file := range files {
			// Build output path maintaining folder structure
			relPath, _ := filepath.Rel(inputFolder, file)
			dir := filepath.Dir(relPath)
			baseName := filepath.Base(relPath)
			ext := filepath.Ext(baseName)
			nameNoExt := baseName[:len(baseName)-len(ext)]
			outputPath := filepath.Join(output, dir, nameNoExt+outputFormat.FileExtension())

			pool.Submit(compressor.Job{
				InputPath:      file,
				OutputPath:     outputPath,
				TargetSizeKB:   size,
				Format:         outputFormat,
				MinQuality:     minQuality,
				MaxQuality:     maxQuality,
				InitialQuality: quality,
			})
		}
		pool.Done()
	}()

	// Collect all results with a live progress bar
	progressBar := ui.NewProgressBar(len(files))
	var results []compressor.Result
	for result := range pool.Results() {
		results = append(results, result)
		progressBar.Increment()
	}
	progressBar.Finish()

	// Print all file results after progress completes (clean output)
	for _, result := range results {
		ui.PrintFileResult(result, size)
	}
	fmt.Println()

	// Print summary
	elapsed := time.Since(startTime)
	ui.PrintSummary(results, elapsed, size)

	return nil
}
