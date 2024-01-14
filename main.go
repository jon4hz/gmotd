package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/defaults"
	"github.com/jon4hz/gmotd/message"
	"golang.org/x/term"
)

func main() {
	ctx := context.New()
	defer ctx.Cancel()

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
	for _, pipe := range message.Message {
		if err := pipe.Gather(ctx); err != nil {
			log.Println(err)
			continue
		}
		if msg := pipe.Print(ctx); msg != "" {
			messages = append(messages, msg, "")
		}
	}

	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, messages...))
}
