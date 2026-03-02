package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
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

// wizardKeyMap returns a huh KeyMap with both ctrl+c and esc bound to quit,
// so pressing either key goes back one step (or exits at the first step).
func wizardKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc"))
	return km
}

// stepResult controls the state machine flow.
type stepResult int

const (
	stepNext  stepResult = iota
	stepBack
	stepAbort
)

// wizardState holds all collected values across wizard steps.
type wizardState struct {
	inputFolder         string
	outputFolder        string
	formatChoice        string
	presetChoice        string
	customSizeStr       string
	customQualityStr    string
	customMinQualityStr string
	customMaxQualityStr string
	customWorkersStr    string
}

func runWizard() error {
	ui.PrintLogo()

	state := &wizardState{
		formatChoice:        "webp",
		presetChoice:        "recommended",
		outputFolder:        "compressed",
		customSizeStr:       "500",
		customQualityStr:    "85",
		customMinQualityStr: "60",
		customMaxQualityStr: "90",
		customWorkersStr:    strconv.Itoa(runtime.NumCPU()),
	}

	steps := []func(*wizardState) stepResult{
		runStepFolder,
		runStepFormat,
		runStepPreset,
		runStepCustom,  // auto-skipped if preset != "custom"
		runStepOutput,
		runStepConfirm,
	}

	i := 0
	for i < len(steps) {
		result := steps[i](state)
		switch result {
		case stepNext:
			i++
		case stepBack:
			if i == 0 {
				fmt.Println(ui.DimStyle.Render("\n  Cancelled."))
				os.Exit(130)
			}
			i--
			// Skip custom step when going backwards if preset != custom
			if i == 3 && state.presetChoice != "custom" {
				i--
			}
		case stepAbort:
			fmt.Println(ui.DimStyle.Render("\n  Cancelled."))
			os.Exit(130)
		}
	}

	config := buildConfig(state)
	fmt.Println()
	return execute(config)
}

// ── Step 1: Input folder ────────────────────────────────────────────────────

func runStepFolder(state *wizardState) stepResult {
	detected := detectImageFolders()

	for {
		if len(detected) == 0 {
			// No nearby folders — show browse vs manual entry options
			var choice string
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Select image folder").
						Description("No image folders detected nearby.  Press Esc or Ctrl+C to cancel.").
						Options(
							huh.NewOption("Browse for folder...", "browse"),
							huh.NewOption("Enter folder path manually...", "manual"),
						).
						Value(&choice),
				),
			).WithTheme(huh.ThemeCharm()).WithKeyMap(wizardKeyMap())

			err := form.Run()
			if err == huh.ErrUserAborted {
				return stepBack
			}
			if err != nil {
				return stepAbort
			}

			if choice == "browse" {
				folder, cancelled := pickFolderFinder("Select your image folder")
				if cancelled {
					continue
				}
				if folder != "" {
					state.inputFolder = folder
					return stepNext
				}
			} else {
				folder, cancelled := pickFolderManual()
				if cancelled {
					continue
				}
				if folder != "" {
					state.inputFolder = folder
					return stepNext
				}
			}
			continue
		}

		// Build quick-select options from detected folders
		var options []huh.Option[string]
		for _, f := range detected {
			options = append(options, huh.NewOption(
				fmt.Sprintf("%s  (%d images)", f.path, f.count),
				f.path,
			))
		}
		options = append(options, huh.NewOption("Browse for another folder...", browsePlaceholder))

		var choice string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select image folder").
					Description("Detected nearby.  Press Esc or Ctrl+C to cancel.").
					Options(options...).
					Value(&choice),
			),
		).WithTheme(huh.ThemeCharm()).WithKeyMap(wizardKeyMap())

		err := form.Run()
		if err == huh.ErrUserAborted {
			return stepBack
		}
		if err != nil {
			return stepAbort
		}

		if choice == browsePlaceholder {
			folder, cancelled := pickFolder("Select your image folder")
			if cancelled {
				// User closed Finder — loop back to folder list, don't abort
				continue
			}
			if folder != "" {
				state.inputFolder = folder
				return stepNext
			}
			continue
		}

		state.inputFolder = choice
		return stepNext
	}
}

// ── Step 2: Output format ───────────────────────────────────────────────────

func runStepFormat(state *wizardState) stepResult {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Output format").
				Description("WebP is recommended for most use cases.  Press Esc or Ctrl+C to go back.").
				Options(
					huh.NewOption("WebP  (recommended)", "webp"),
					huh.NewOption("AVIF  (smaller, slower encode)", "avif"),
					huh.NewOption("JPEG  (universal compatibility)", "jpeg"),
					huh.NewOption("PNG   (lossless)", "png"),
				).
				Value(&state.formatChoice),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(wizardKeyMap())

	err := form.Run()
	if err == huh.ErrUserAborted {
		return stepBack
	}
	if err != nil {
		return stepAbort
	}
	return stepNext
}

// ── Step 3: Preset ──────────────────────────────────────────────────────────

func runStepPreset(state *wizardState) stepResult {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Compression settings").
				Description("Choose a preset or configure manually.  Press Esc or Ctrl+C to go back.").
				Options(
					huh.NewOption("Recommended — 500 KB, quality 85, balanced", "recommended"),
					huh.NewOption("Web Optimized — 200 KB, quality 75, aggressive", "web"),
					huh.NewOption("High Quality — 2 MB, quality 95, minimal compression", "high-quality"),
					huh.NewOption("Custom — choose your own settings", "custom"),
				).
				Value(&state.presetChoice),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(wizardKeyMap())

	err := form.Run()
	if err == huh.ErrUserAborted {
		return stepBack
	}
	if err != nil {
		return stepAbort
	}
	return stepNext
}

// ── Step 4: Custom settings (skipped if preset != "custom") ────────────────

func runStepCustom(state *wizardState) stepResult {
	if state.presetChoice != "custom" {
		return stepNext
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Target file size (KB)").
				Description("Images will be compressed to fit under this size.").
				Placeholder("500").
				Value(&state.customSizeStr).
				Validate(validatePositiveInt),

			huh.NewInput().
				Title("Initial quality (1-100)").
				Placeholder("85").
				Value(&state.customQualityStr).
				Validate(validateQuality),

			huh.NewInput().
				Title("Min quality (1-100)").
				Description("Quality floor — compression won't go below this.").
				Placeholder("60").
				Value(&state.customMinQualityStr).
				Validate(validateQuality),

			huh.NewInput().
				Title("Max quality (1-100)").
				Description("Quality ceiling — compression won't exceed this.").
				Placeholder("90").
				Value(&state.customMaxQualityStr).
				Validate(validateQuality),

			huh.NewInput().
				Title("Worker count").
				Description(fmt.Sprintf("Your machine has %d CPU cores.  Press Esc or Ctrl+C to go back.", runtime.NumCPU())).
				Placeholder(strconv.Itoa(runtime.NumCPU())).
				Value(&state.customWorkersStr).
				Validate(validatePositiveInt),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(wizardKeyMap())

	err := form.Run()
	if err == huh.ErrUserAborted {
		return stepBack
	}
	if err != nil {
		return stepAbort
	}
	return stepNext
}

// ── Step 5: Output folder ───────────────────────────────────────────────────

func runStepOutput(state *wizardState) stepResult {
	for {
		var choice string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Output folder").
					Description("Where should compressed images be saved?  Press Esc or Ctrl+C to go back.").
					Options(
						huh.NewOption("compressed/  (default, created in current folder)", "compressed"),
						huh.NewOption("Browse for a custom output location...", browsePlaceholder),
					).
					Value(&choice),
			),
		).WithTheme(huh.ThemeCharm()).WithKeyMap(wizardKeyMap())

		err := form.Run()
		if err == huh.ErrUserAborted {
			return stepBack
		}
		if err != nil {
			return stepAbort
		}

		if choice != browsePlaceholder {
			state.outputFolder = choice
			return stepNext
		}

		// Browse: pick parent directory, then name the output folder
		parentDir, cancelled := pickFolder("Select where to save compressed images")
		if cancelled {
			// Finder closed — re-show the output selection
			continue
		}

		var folderName string
		nameForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Output folder name").
					Description(fmt.Sprintf("Will be created inside: %s", parentDir)).
					Placeholder("compressed").
					Value(&folderName).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("name is required")
						}
						return nil
					}),
			),
		).WithTheme(huh.ThemeCharm()).WithKeyMap(wizardKeyMap())

		if err := nameForm.Run(); err == huh.ErrUserAborted {
			continue // back to output selection
		}

		state.outputFolder = filepath.Join(parentDir, folderName)
		return stepNext
	}
}

// ── Step 6: Confirm ─────────────────────────────────────────────────────────

func runStepConfirm(state *wizardState) stepResult {
	// Scan images for count display
	files, _ := scanner.FindImages(state.inputFolder, false)
	printWizardSummary(state, len(files))

	confirmed := true
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Start compression?").
				Affirmative("Let's go!").
				Negative("Cancel").
				Value(&confirmed),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(wizardKeyMap())

	err := form.Run()
	if err == huh.ErrUserAborted {
		return stepBack
	}
	if err != nil {
		return stepAbort
	}
	if !confirmed {
		return stepBack
	}
	return stepNext
}

// ── Folder picker ───────────────────────────────────────────────────────────

// pickFolder opens a native Finder dialog on macOS, falling back to text input.
// Returns (folder, cancelled).
func pickFolder(prompt string) (string, bool) {
	if runtime.GOOS == "darwin" {
		if _, err := exec.LookPath("osascript"); err == nil {
			return pickFolderFinder(prompt)
		}
	}
	return pickFolderManual()
}

// pickFolderFinder opens a native macOS Finder folder selection dialog.
// Returns (folder, cancelled).
func pickFolderFinder(prompt string) (string, bool) {
	folder, err := macFolderDialog(prompt)
	if err != nil {
		return "", true
	}
	return folder, false
}

// pickFolderManual shows a text input for typing a folder path.
// Returns (folder, cancelled).
func pickFolderManual() (string, bool) {
	var folder string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Folder path").
				Description("Tip: paste a full path with Cmd+V, or drag a folder from Finder into this window.").
				Placeholder("/Users/you/Pictures").
				Value(&folder).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("path is required")
					}
					info, err := os.Stat(s)
					if err != nil {
						return fmt.Errorf("cannot access: %s", s)
					}
					if !info.IsDir() {
						return fmt.Errorf("not a directory: %s", s)
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(wizardKeyMap())

	if err := form.Run(); err != nil {
		return "", true
	}
	return folder, false
}

// macFolderDialog opens a native macOS Finder folder selection dialog.
func macFolderDialog(prompt string) (string, error) {
	script := fmt.Sprintf(`tell application "Finder" to activate
set chosenFolder to choose folder with prompt "%s"
return POSIX path of chosenFolder`, prompt)

	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", err
	}

	folder := strings.TrimSpace(string(out))
	folder = strings.TrimRight(folder, "/")
	return folder, nil
}

// ── Folder detection ────────────────────────────────────────────────────────

type imageFolder struct {
	path  string
	count int
}

func detectImageFolders() []imageFolder {
	var folders []imageFolder

	if count := countImages("."); count > 0 {
		folders = append(folders, imageFolder{path: ".", count: count})
	}

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

// ── Config builder ──────────────────────────────────────────────────────────

func buildConfig(state *wizardState) ShrinkConfig {
	config := ShrinkConfig{
		InputFolder: state.inputFolder,
		Output:      state.outputFolder,
		Format:      state.formatChoice,
		Recursive:   false,
	}

	if p, ok := presets[state.presetChoice]; ok {
		config.Size = p.Size
		config.Quality = p.Quality
		config.MinQuality = p.MinQuality
		config.MaxQuality = p.MaxQuality
		config.Workers = p.Workers
	} else {
		config.Size, _ = strconv.Atoi(state.customSizeStr)
		config.Quality, _ = strconv.Atoi(state.customQualityStr)
		config.MinQuality, _ = strconv.Atoi(state.customMinQualityStr)
		config.MaxQuality, _ = strconv.Atoi(state.customMaxQualityStr)
		config.Workers, _ = strconv.Atoi(state.customWorkersStr)
	}

	return config
}

// ── Summary ─────────────────────────────────────────────────────────────────

func printWizardSummary(state *wizardState, imageCount int) {
	cfg := buildConfig(state)

	content := fmt.Sprintf(
		"  %s  %s\n"+
			"  %s  %s\n"+
			"  %s  %s\n"+
			"  %s  %s\n"+
			"  %s  %s\n"+
			"  %s  %s\n"+
			"  %s  %s",
		ui.LabelStyle.Render("Input:"),
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

// ── Validators ───────────────────────────────────────────────────────────────

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
