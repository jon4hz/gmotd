package hostname

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	_ "embed"

	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/gmotd/context"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/lukesampson/figlet/figletlib"
)

type Section struct{}

func (Section) String() string { return "hostname" }

func (Section) Enabled(c *context.Context) bool {
	return !c.Config.Hostname.Disabled
}

func (Section) Default(ctx *context.Context) {
	ctx.Config.Hostname.Figlet = true
	ctx.Config.Hostname.Color = "rainbow"
	ctx.Config.Hostname.FigletFont = "standard"
	ctx.Config.Hostname.FigletFontDir = "/usr/share/figlet/fonts"
}

func (Section) Gather(c *context.Context) error {
	h, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname: %w", err)
	}
	c.Hostname = &context.Hostname{
		Hostname: h,
	}
	return nil
}

//go:embed fonts/standard.flf
var defaultFont []byte

func (Section) Print(c *context.Context) string {
	if !c.Config.Hostname.Figlet {
		return c.Hostname.Hostname
	}

	f, err := figletlib.GetFontByName(c.Config.Hostname.FigletFontDir, c.Config.Hostname.FigletFont)
	if err != nil {
		f, err = figletlib.ReadFontFromBytes(defaultFont)
		if err != nil {
			return ""
		}
	}

	renderStr := figletlib.SprintMsg(c.Hostname.Hostname, f, 200, f.Settings(), "left")

	colors := func() string {
		colors := colorGrid(lipgloss.Width(renderStr), lipgloss.Height(renderStr))
		b := strings.Builder{}
		for i, x := range parseHostname(renderStr) {
			for j, y := range x {
				s := lipgloss.NewStyle().SetString(y).Foreground(lipgloss.Color(colors[i][j]))
				b.WriteString(s.String())
			}
			if i < len(colors)-1 {
				b.WriteByte('\n')
			}
		}
		return b.String()
	}()

	return colors
}

func parseLine(l string) []string {
	var line []string
	chars := strings.Split(l, "")
	line = append(line, chars...)
	return line
}

func parseHostname(h string) [][]string {
	var hostname [][]string
	lines := strings.Split(h, "\n")
	for _, line := range lines {
		pl := parseLine(line)
		hostname = append(hostname, pl)
	}
	return hostname
}

func colorGrid(xSteps, ySteps int) [][]string {
	c0, c1 := genGradientTable(xSteps)

	x0y0, _ := colorful.Hex(c0)
	x1y0, _ := colorful.Hex(c1)

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y0, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y0, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := 0; y < xSteps; y++ {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}

var colorChoices = []string{
	"#9e0142",
	"#d53e4f",
	"#f46d43",
	"#fdae61",
	"#fee08b",
	"#ffffbf",
	"#e6f598",
	"#abdda4",
	"#66c2a5",
	"#3288bd",
	"#5e4fa2",
}

func genGradientTable(xSteps int) (string, string) {
	req := int(float64(xSteps) / 9)
	if len(colorChoices) < req {
		req = len(colorChoices)
	}

	var offset int
	if req < len(colorChoices) {
		offset = rand.Intn(len(colorChoices) - req)
	}

	var c0, c1 string
	c0 = colorChoices[offset]

	if offset+req >= len(colorChoices) {
		c1 = colorChoices[len(colorChoices)-(offset+req)]
	} else {
		c1 = colorChoices[offset+req]
	}

	return c0, c1
}
