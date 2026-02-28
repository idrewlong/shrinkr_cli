package ui

import (
	"fmt"
	"strings"

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

// Gradient colors for each line of the logo
var logoColors = []lipgloss.Color{
	ColorLogo1, // red
	ColorLogo2, // orange
	ColorLogo3, // yellow
	ColorLogo4, // green
	ColorLogo5, // cyan
	ColorLogo6, // purple
}

// PrintLogo prints the Shrinkr ASCII art logo with gradient coloring.
func PrintLogo() {
	fmt.Println()

	for i, line := range logoLines {
		color := logoColors[i%len(logoColors)]
		style := lipgloss.NewStyle().Foreground(color).Bold(true)
		fmt.Println(style.Render(line))
	}

	tagline := lipgloss.NewStyle().
		Foreground(ColorGray).
		Italic(true).
		Render("          Image Compression Tool by Andrew Long")

	fmt.Println(tagline)
	fmt.Println(strings.Repeat("─", 58))
	fmt.Println()
}
