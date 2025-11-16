package formatter

import "fmt"

// GetFormatter returns a formatter by name
func GetFormatter(name string) (Formatter, error) {
	switch name {
	case "full":
		return NewFullFormatter(), nil
	case "json":
		return NewJSONFormatter(), nil
	case "summary":
		return NewSummaryFormatter(), nil
	case "markdown", "md":
		return NewMarkdownFormatter(), nil
	default:
		return nil, fmt.Errorf("unknown formatter: %s (available: full, json, summary, markdown)", name)
	}
}

// AvailableFormatters returns a list of available formatter names
func AvailableFormatters() []string {
	return []string{"full", "json", "summary", "markdown"}
}
