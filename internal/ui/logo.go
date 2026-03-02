package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ASCII art logo lines for "SHRINKR"
var logoLines = []string{
	`  ███████╗██╗  ██╗██████╗ ██╗███╗   ██╗██╗  ██╗██████╗ `,
	`  ██╔════╝██║  ██║██╔══██╗██║████╗  ██║██║ ██╔╝██╔══██╗`,
	`  ███████╗███████║██████╔╝██║██╔██╗ ██║█████╔╝ ██████╔╝`,
	`  ╚════██║██╔══██║██╔══██╗██║██║╚██╗██║██╔═██╗ ██╔══██╗`,
	`  ███████║██║  ██║██║  ██║██║██║ ╚████║██║  ██╗██║  ██║`,
	`  ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝ ╚═══╝╚═╝  ╚═╝╚═╝  ╚═╝`,
}

// Gradient colors for the logo reveal animation
var logoGradient = []lipgloss.Color{

	ColorLogo2, // orange

}

// PrintLogo prints the Shrinkr ASCII art logo with an animated reveal.
// Each line drops in with a gradient color, then the tagline fades in.
func PrintLogo() {
	fmt.Println()

	// Reveal each logo line with a staggered delay
	for i, line := range logoLines {
		color := logoGradient[i%len(logoGradient)]
		style := lipgloss.NewStyle().Foreground(color).Bold(true)
		fmt.Println(style.Render(line))
		time.Sleep(45 * time.Millisecond)
	}

	// Brief pause before tagline
	// time.Sleep(80 * time.Millisecond)

	// tagline := lipgloss.NewStyle().
	// 	Foreground(ColorGray).
	// 	Italic(true).
	// 	Render("          Image Compression Tool")

	// fmt.Println(tagline)

	// Animate the separator line drawing across
	sep := "─"
	width := 58
	for i := 1; i <= width; i++ {
		fmt.Printf("\r%s", DimStyle.Render(strings.Repeat(sep, i)))
		time.Sleep(4 * time.Millisecond)
	}
	fmt.Println()
	fmt.Println()
}
