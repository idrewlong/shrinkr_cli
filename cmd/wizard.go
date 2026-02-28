package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/idrewlong/shrinkr_cli/internal/scanner"
	"github.com/idrewlong/shrinkr_cli/internal/ui"
)

type preset struct {
	Size       int
	Quality    int
	MinQuality int
	MaxQuality int
	Workers    int
}

var presets = map[string]preset{
	"recommended": {
		Size:       500,
		Quality:    85,
		MinQuality: 60,
		MaxQuality: 90,
		Workers:    runtime.NumCPU(),
	},
	"web": {
		Size:       200,
		Quality:    75,
		MinQuality: 50,
		MaxQuality: 80,
		Workers:    runtime.NumCPU(),
	},
	"high-quality": {
		Size:       2048,
		Quality:    95,
		MinQuality: 85,
		MaxQuality: 98,
		Workers:    runtime.NumCPU(),
	},
}

const browsePlaceholder = "__browse__"

func runWizard() error {
	ui.PrintLogo()

	// Step 1: Detect folders with images
	inputFolder, err := selectFolder()
	if err != nil {
		return err
	}

	var (
		formatChoice string
		presetChoice string
	)

	// Custom settings as strings (huh.NewInput binds to *string)
	customSizeStr := "500"
	customQualityStr := "85"
	customMinQualityStr := "60"
	customMaxQualityStr := "90"
	customWorkersStr := strconv.Itoa(runtime.NumCPU())

	formatChoice = "webp"

	// Step 2: Collect format, preset, and custom settings
	form := huh.NewForm(
		// Output format
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Output format").
				Description("WebP offers the best size-to-quality ratio for most use cases.").
				Options(
					huh.NewOption("WebP  (recommended)", "webp"),
					huh.NewOption("AVIF  (smaller, slower encode)", "avif"),
					huh.NewOption("JPEG  (universal compatibility)", "jpeg"),
					huh.NewOption("PNG   (lossless)", "png"),
				).
				Value(&formatChoice),
		),

		// Preset selection
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Compression settings").
				Description("Choose a preset or configure manually.").
				Options(
					huh.NewOption("Recommended — 500 KB, quality 85, balanced", "recommended"),
					huh.NewOption("Web Optimized — 200 KB, quality 75, aggressive", "web"),
					huh.NewOption("High Quality — 2 MB, quality 95, minimal compression", "high-quality"),
					huh.NewOption("Custom — choose your own settings", "custom"),
				).
				Value(&presetChoice),
		),

		// Custom settings (only shown when "custom" is selected)
		huh.NewGroup(
			huh.NewInput().
				Title("Target file size (KB)").
				Description("Images will be compressed to fit under this size.").
				Placeholder("500").
				Value(&customSizeStr).
				Validate(validatePositiveInt),

			huh.NewInput().
				Title("Initial quality (1-100)").
				Placeholder("85").
				Value(&customQualityStr).
				Validate(validateQuality),

			huh.NewInput().
				Title("Min quality (1-100)").
				Description("Quality floor — compression won't go below this.").
				Placeholder("60").
				Value(&customMinQualityStr).
				Validate(validateQuality),

			huh.NewInput().
				Title("Max quality (1-100)").
				Description("Quality ceiling — compression won't exceed this.").
				Placeholder("90").
				Value(&customMaxQualityStr).
				Validate(validateQuality),

			huh.NewInput().
				Title("Worker count").
				Description(fmt.Sprintf("Your machine has %d CPU cores.", runtime.NumCPU())).
				Placeholder(strconv.Itoa(runtime.NumCPU())).
				Value(&customWorkersStr).
				Validate(validatePositiveInt),
		).WithHideFunc(func() bool {
			return presetChoice != "custom"
		}),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return handleAbort(err)
	}

	// Build config from wizard answers
	config := buildConfigFromWizard(inputFolder, formatChoice, presetChoice,
		customSizeStr, customQualityStr, customMinQualityStr, customMaxQualityStr, customWorkersStr)

	// Scan for images to show count in summary
	files, err := scanner.FindImages(config.InputFolder, config.Recursive)
	if err != nil {
		return fmt.Errorf("error scanning folder: %w", err)
	}

	// Print summary before running
	printWizardSummary(config, len(files))

	// Step 3: Confirm
	confirmed := true
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Start compression?").
				Affirmative("Let's go!").
				Negative("Cancel").
				Value(&confirmed),
		),
	).WithTheme(huh.ThemeCharm())

	if err := confirmForm.Run(); err != nil {
		return handleAbort(err)
	}

	if !confirmed {
		fmt.Println(ui.DimStyle.Render("\n  Cancelled."))
		return nil
	}

	fmt.Println()
	return execute(config)
}

// selectFolder auto-detects folders with images and lets the user pick one,
// or browse the filesystem if none are found or the user wants a different folder.
func selectFolder() (string, error) {
	detected := detectImageFolders()

	if len(detected) == 0 {
		// No folders found — go straight to file picker
		return browseForFolder()
	}

	// Build options from detected folders
	var options []huh.Option[string]
	for _, f := range detected {
		label := fmt.Sprintf("%s  (%d images)", f.path, f.count)
		options = append(options, huh.NewOption(label, f.path))
	}
	options = append(options, huh.NewOption("Browse for another folder...", browsePlaceholder))

	var choice string
	selectForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select image folder").
				Description("Folders with images detected nearby.").
				Options(options...).
				Value(&choice),
		),
	).WithTheme(huh.ThemeCharm())

	if err := selectForm.Run(); err != nil {
		return "", handleAbort(err)
	}

	if choice == browsePlaceholder {
		return browseForFolder()
	}

	return choice, nil
}

// browseForFolder opens the huh FilePicker for filesystem navigation.
func browseForFolder() (string, error) {
	var folder string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Title("Browse for image folder").
				Description("Navigate to the folder containing your images.").
				CurrentDirectory(".").
				FileAllowed(false).
				DirAllowed(true).
				ShowHidden(false).
				ShowSize(true).
				Value(&folder),
		),
	).WithTheme(huh.ThemeCharm())

	if err := form.Run(); err != nil {
		return "", handleAbort(err)
	}

	return folder, nil
}

type imageFolder struct {
	path  string
	count int
}

// detectImageFolders scans the current directory and one level of subdirectories
// for folders that contain supported image files.
func detectImageFolders() []imageFolder {
	var folders []imageFolder

	// Check current directory
	if count := countImages("."); count > 0 {
		cwd, _ := os.Getwd()
		folders = append(folders, imageFolder{path: ".", count: count})
		_ = cwd
	}

	// Check immediate subdirectories
	entries, err := os.ReadDir(".")
	if err != nil {
		return folders
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name()[0] == '.' {
			continue
		}
		path := entry.Name()
		if count := countImages(path); count > 0 {
			folders = append(folders, imageFolder{path: path, count: count})
		}
	}

	// Also check parent directory's immediate children (sibling folders)
	parentEntries, err := os.ReadDir("..")
	if err != nil {
		return folders
	}

	cwd, _ := os.Getwd()
	currentBase := filepath.Base(cwd)

	for _, entry := range parentEntries {
		if !entry.IsDir() || entry.Name()[0] == '.' || entry.Name() == currentBase {
			continue
		}
		path := filepath.Join("..", entry.Name())
		if count := countImages(path); count > 0 {
			folders = append(folders, imageFolder{path: path, count: count})
		}
	}

	return folders
}

// countImages returns the number of supported image files in a directory (non-recursive).
func countImages(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() && scanner.IsSupportedImage(entry.Name()) {
			count++
		}
	}
	return count
}

func handleAbort(err error) error {
	if err == huh.ErrUserAborted {
		fmt.Println(ui.DimStyle.Render("\n  Cancelled."))
		os.Exit(130)
	}
	return fmt.Errorf("wizard error: %w", err)
}

func buildConfigFromWizard(
	inputFolder, formatChoice, presetChoice,
	sizeStr, qualityStr, minQStr, maxQStr, workersStr string,
) ShrinkConfig {
	config := ShrinkConfig{
		InputFolder: inputFolder,
		Output:      "compressed",
		Format:      formatChoice,
		Recursive:   false,
	}

	if p, ok := presets[presetChoice]; ok {
		config.Size = p.Size
		config.Quality = p.Quality
		config.MinQuality = p.MinQuality
		config.MaxQuality = p.MaxQuality
		config.Workers = p.Workers
	} else {
		config.Size, _ = strconv.Atoi(sizeStr)
		config.Quality, _ = strconv.Atoi(qualityStr)
		config.MinQuality, _ = strconv.Atoi(minQStr)
		config.MaxQuality, _ = strconv.Atoi(maxQStr)
		config.Workers, _ = strconv.Atoi(workersStr)
	}

	return config
}

func printWizardSummary(cfg ShrinkConfig, imageCount int) {
	content := fmt.Sprintf(
		"  %s  %s\n"+
			"  %s  %s\n"+
			"  %s  %s\n"+
			"  %s  %s\n"+
			"  %s  %s\n"+
			"  %s  %s\n"+
			"  %s  %s",
		ui.LabelStyle.Render("Folder:"),
		ui.ValueStyle.Render(cfg.InputFolder),
		ui.LabelStyle.Render("Images:"),
		ui.TitleStyle.Render(fmt.Sprintf("%d found", imageCount)),
		ui.LabelStyle.Render("Format:"),
		ui.AccentStyle.Render(cfg.Format),
		ui.LabelStyle.Render("Target:"),
		ui.ValueStyle.Render(fmt.Sprintf("%d KB", cfg.Size)),
		ui.LabelStyle.Render("Quality:"),
		ui.ValueStyle.Render(fmt.Sprintf("%d (range %d–%d)", cfg.Quality, cfg.MinQuality, cfg.MaxQuality)),
		ui.LabelStyle.Render("Workers:"),
		ui.ValueStyle.Render(strconv.Itoa(cfg.Workers)),
		ui.LabelStyle.Render("Output:"),
		ui.ValueStyle.Render(cfg.Output),
	)

	fmt.Println(ui.HeaderBox.Render(content))
	fmt.Println()
}

func validatePositiveInt(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("must be a number")
	}
	if n <= 0 {
		return fmt.Errorf("must be greater than 0")
	}
	return nil
}

func validateQuality(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("must be a number")
	}
	if n < 1 || n > 100 {
		return fmt.Errorf("must be between 1 and 100")
	}
	return nil
}
