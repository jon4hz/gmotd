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
	ctx.Runtime.Width = w
	ctx.Runtime.Height = h

	for _, d := range defaults.Defaulters {
		d.Default(ctx)
	}

	var messages []string
	for _, pipe := range pipeline.Pipeline {
		if err := pipe.Gather(ctx); err != nil {
			log.Println(err)
			continue
		}
		messages = append(messages, pipe.Print(ctx))
	}

	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, messages...))
}
