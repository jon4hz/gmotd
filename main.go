package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/jon4hz/gmotd/config"
	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/message"
)

func main() {
	ctx := context.New()
	defer ctx.Cancel()

	for _, m := range message.Message {
		d, ok := m.(message.Defaulter)
		if !ok {
			continue
		}
		d.Default(ctx)
	}

	err := config.Load(ctx.Config)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, section := range message.Message {
		if !section.Enabled(ctx) {
			continue
		}
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
		if !section.Enabled(ctx) {
			continue
		}
		if msg := section.Print(ctx); msg != "" {
			messages = append(messages, msg, "")
		}
	}

	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, messages...))
}
