package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/gmotd/internal/context"
	"github.com/jon4hz/gmotd/internal/defaults"
	"github.com/jon4hz/gmotd/internal/pipeline"
	"golang.org/x/term"
)

func main() {
	ctx := context.New()

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	ctx.Width = w
	ctx.Height = h

	for _, d := range defaults.Defaulters {
		d.Default(ctx)
	}

	var messages []string
	for _, pipe := range pipeline.Pipeline {
		messages = append(messages, pipe.Message(ctx))
	}

	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, messages...))
}
