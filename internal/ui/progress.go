package ui

import (
	"fmt"
	"strings"
	"sync"
)

// ProgressBar renders a styled progress bar that updates in-place.
type ProgressBar struct {
	total   int
	current int
	width   int
	mu      sync.Mutex
}

// NewProgressBar creates a new progress bar for the given total.
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		total: total,
		width: 40,
	}
}

// Start renders the initial empty progress bar before any results arrive.
func (p *ProgressBar) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.render()
}

// Increment advances the progress bar by one and re-renders.
func (p *ProgressBar) Increment() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current++
	p.render()
}

// render draws the progress bar on the current line using carriage return.
func (p *ProgressBar) render() {
	pct := float64(p.current) / float64(p.total)
	filled := int(pct * float64(p.width))
	if filled > p.width {
		filled = p.width
	}

	bar := ProgressFilled.Render(strings.Repeat("█", filled)) +
		ProgressEmpty.Render(strings.Repeat("░", p.width-filled))

	status := ProgressText.Render(fmt.Sprintf(" %3d%% | %d/%d files",
		int(pct*100), p.current, p.total))

	fmt.Printf("\r  [%s]%s", bar, status)
}

// Finish completes the progress bar and moves to the next line.
func (p *ProgressBar) Finish() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = p.total
	p.render()
	fmt.Println()
	fmt.Println()
}
