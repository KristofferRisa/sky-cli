package ui

import "github.com/fatih/color"

var (
	// Bold makes text bold
	Bold = color.New(color.Bold).SprintFunc()

	// Blue for headers and borders
	Blue     = color.New(color.FgBlue).SprintFunc()
	BlueBold = color.New(color.FgBlue, color.Bold).SprintFunc()

	// Green for positive/good conditions
	Green     = color.New(color.FgGreen).SprintFunc()
	GreenBold = color.New(color.FgGreen, color.Bold).SprintFunc()

	// Yellow for warnings
	Yellow     = color.New(color.FgYellow).SprintFunc()
	YellowBold = color.New(color.FgYellow, color.Bold).SprintFunc()

	// Red for alerts/dangerous conditions
	Red     = color.New(color.FgRed).SprintFunc()
	RedBold = color.New(color.FgRed, color.Bold).SprintFunc()

	// Cyan for informational text
	Cyan     = color.New(color.FgCyan).SprintFunc()
	CyanBold = color.New(color.FgCyan, color.Bold).SprintFunc()
)

// DisableColors turns off color output
func DisableColors() {
	color.NoColor = true
}

// EnableColors turns on color output
func EnableColors() {
	color.NoColor = false
}

// Header prints a formatted section header
func Header(text string) string {
	line := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	return "\n" + BlueBold(line) + "\n" + Bold(text) + "\n" + BlueBold(line)
}
