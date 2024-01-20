package styles

import (
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/constraints"
)

var (
	Indent = lipgloss.NewStyle().PaddingLeft(2)

	Green  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9ece6a"))
	Orange = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fd9864"))
	Red    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fd665f"))
)

type Floatable interface {
	constraints.Float | constraints.Integer
}

func ColorizeWithMax[T Floatable](x T, max T) lipgloss.Style {
	p := (float64(x) / float64(max)) * 100
	switch {
	// err
	case p >= 90:
		return Red.Copy()
	// warn
	case p >= 80:
		return Orange.Copy()
	// ok
	default:
		return Green.Copy()
	}
}

func ColorizeWithMin[T Floatable](x T, min T) lipgloss.Style {
	p := (float64(x) / float64(min)) * 100
	switch {
	// err
	case p <= 10:
		return Red.Copy()
	// warn
	case p <= 20:
		return Orange.Copy()
	// ok
	default:
		return Green.Copy()
	}
}
