package main

import (
	"fmt"
	"log"
	"os"
	"sync"

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

	var wg sync.WaitGroup
	for _, section := range message.Message {
		wg.Add(1)
		go func(section message.Section) {
			defer wg.Done()
			if err := section.Gather(ctx); err != nil {
				log.Println(err)
				return
			}
		}(section)
	}
	wg.Wait()

	var messages []string
	for _, section := range message.Message {
		if msg := section.Print(ctx); msg != "" {
			messages = append(messages, msg, "")
		}
	}

	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, messages...))
}
