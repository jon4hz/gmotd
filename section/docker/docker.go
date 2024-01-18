package docker

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/styles"
)

type Section struct{}

func (Section) String() string { return "docker" }

func (Section) Gather(c *context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	containers, err := cli.ContainerList(c, types.ContainerListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to get docker container list: %w", err)
	}

	dcs := make([]context.DockerContainer, len(containers))
	for i, container := range containers {
		dcs[i] = context.DockerContainer{
			Name:    strings.TrimPrefix(container.Names[0], "/"),
			State:   container.State,
			Healthy: !strings.Contains(container.Status, "unhealthy"),
		}
	}

	sort.Slice(dcs, func(i, j int) bool {
		return dcs[i].Name < dcs[j].Name
	})

	c.Docker = &context.Docker{
		Containers: dcs,
	}

	return nil
}

var (
	keyStyle        = lipgloss.NewStyle()
	valueStyle      = lipgloss.NewStyle().Bold(true)
	leftColumnStyle = lipgloss.NewStyle().PaddingRight(4)
)

func (Section) Print(ctx *context.Context) string {
	c := ctx.Docker
	if c == nil {
		return ""
	}

	containers := make([]string, 0, len(c.Containers))
	oddAdd := 0
	if cap(containers)%2 != 0 {
		oddAdd = 1
	}

	containers = append(containers, renderColumn(c.Containers[:len(c.Containers)/2+oddAdd])...)
	containers = append(containers, renderColumn(c.Containers[len(c.Containers)/2+oddAdd:])...)

	leftColumn := lipgloss.JoinVertical(lipgloss.Top, containers[:len(containers)/2+oddAdd]...)
	bothColumns := lipgloss.JoinHorizontal(lipgloss.Left,
		leftColumnStyle.Render(leftColumn),
		lipgloss.JoinVertical(lipgloss.Top,
			containers[len(containers)/2+oddAdd:]...,
		),
	)

	return lipgloss.JoinVertical(lipgloss.Top,
		"docker status:",
		styles.Indent.Render(bothColumns),
	)
}

func renderColumn(c []context.DockerContainer) []string {
	lnw := longestNameWidth(c)
	containers := make([]string, len(c))
	for i, container := range c {
		containers[i] = renderContainer(container, lnw)
	}
	return containers
}

func renderContainer(c context.DockerContainer, lnw int) string {
	switch c.State {
	case "running":
		if c.Healthy {
			valueStyle.Foreground(styles.Green.GetForeground())
			c.State = "up"
		} else {
			valueStyle.Foreground(styles.Red.GetForeground())
			c.State = "unhealthy"
		}
	case "exited":
		valueStyle.Foreground(styles.Orange.GetForeground())
	}

	return keyStyle.Width(lnw).Render(c.Name+":") + "  " + valueStyle.Render(c.State)
}

func longestNameWidth(c []context.DockerContainer) int {
	l := 0
	for _, container := range c {
		if w := lipgloss.Width(container.Name) + 1; w > l {
			l = w
		}
	}
	return l
}
