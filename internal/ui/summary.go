package ui

import (
	"fmt"
	"time"

	"github.com/idrewlong/shrinkr_cli/internal/compressor"
	"github.com/idrewlong/shrinkr_cli/internal/util"
)

// PrintFileResult prints the result for a single file compression.
func PrintFileResult(result compressor.Result, targetSizeKB int) {
	if !result.Success {
		fmt.Printf("  %s %s\n",
			ErrorStyle.Render("✗"),
			ErrorStyle.Render(result.Filename+": "+result.Error))
		return
	}

	// Check if file exceeded target
	exceededTarget := result.CompressedSize > int64(targetSizeKB)*1024
	sizeStyle := SuccessStyle
	if exceededTarget {
		sizeStyle = WarningStyle
	}

	fmt.Printf("  %s %s\n",
		AccentStyle.Render(">>"),
		FileStyle.Render(result.Filename))

	fmt.Printf("     %s %s  %s  %s %s\n",
		LabelStyle.Render(util.FormatBytes(result.OriginalSize)),
		DimStyle.Render("→"),
		sizeStyle.Render(util.FormatBytes(result.CompressedSize)),
		DimStyle.Render("|"),
		StatStyle.Render(fmt.Sprintf("Q:%d  Saved: %.1f%%", result.Quality, result.CompressionPct)))

	if exceededTarget {
		fmt.Printf("     %s\n",
			WarningStyle.Render(fmt.Sprintf("⚠ Exceeds target of %d KB", targetSizeKB)))
	}
}

// PrintSummary prints the final summary in a styled box.
func PrintSummary(results []compressor.Result, elapsed time.Duration, targetSizeKB int) {
	var totalOriginal, totalCompressed int64
	var successCount, failCount int

	for _, r := range results {
		if r.Success {
			successCount++
			totalOriginal += r.OriginalSize
			totalCompressed += r.CompressedSize
		} else {
			failCount++
		}
	}

	savedBytes := totalOriginal - totalCompressed
	savedPct := 0.0
	if totalOriginal > 0 {
		savedPct = (1.0 - float64(totalCompressed)/float64(totalOriginal)) * 100
	}

	speed := 0.0
	if elapsed.Seconds() > 0 {
		speed = float64(len(results)) / elapsed.Seconds()
	}

	// Build summary content
	title := SuccessStyle.Render("  Compression Complete")
	lines := []string{
		title,
		"",
		fmt.Sprintf("  %s  %s",
			LabelStyle.Render("Success:"),
			SuccessStyle.Render(fmt.Sprintf("%d files", successCount))),
	}

	if failCount > 0 {
		lines = append(lines, fmt.Sprintf("  %s   %s",
			LabelStyle.Render("Failed:"),
			ErrorStyle.Render(fmt.Sprintf("%d files", failCount))))
	}

	lines = append(lines,
		fmt.Sprintf("  %s %s",
			LabelStyle.Render("Original:"),
			ValueStyle.Render(util.FormatBytes(totalOriginal))),
		fmt.Sprintf("  %s  %s",
			LabelStyle.Render("Output:"),
			ValueStyle.Render(util.FormatBytes(totalCompressed))),
		fmt.Sprintf("  %s    %s",
			LabelStyle.Render("Saved:"),
			SuccessStyle.Render(fmt.Sprintf("%s (%.1f%%)", util.FormatBytes(savedBytes), savedPct))),
		fmt.Sprintf("  %s     %s",
			LabelStyle.Render("Time:"),
			ValueStyle.Render(fmt.Sprintf("%.1fs", elapsed.Seconds()))),
		fmt.Sprintf("  %s    %s",
			LabelStyle.Render("Speed:"),
			ValueStyle.Render(fmt.Sprintf("%.1f files/sec", speed))),
	)

	content := ""
	for i, line := range lines {
		content += line
		if i < len(lines)-1 {
			content += "\n"
		}
	}

	fmt.Println(SummaryBox.Render(content))
}
