package main

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/jon4hz/gmotd/config"
	"github.com/jon4hz/gmotd/context"
	"github.com/jon4hz/gmotd/message"
	"github.com/muesli/termenv"
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

	// maybe make this a config option.
	lipgloss.SetColorProfile(termenv.TrueColor)

	var (
		sectionResults = make(map[string]string)
		mu             sync.Mutex
		wg             sync.WaitGroup
		doneC          = make(chan struct{})
	)

	go func() {
		for _, section := range message.Message {
			if !section.Enabled(ctx) {
				continue
			}
			wg.Add(1)
			go func(section message.Section) {
				defer wg.Done()
				if err := section.Gather(ctx); err != nil {
					log.Error("Section gather returned error!", "section", section.String(), "err", err)
					return
				}

				if msg := section.Print(ctx); msg != "" {
					mu.Lock()
					sectionResults[section.String()] = msg + "\n"
					mu.Unlock()
				}

			}(section)
		}
		wg.Wait()
		close(doneC)
	}()

	select {
	case <-ctx.Done():
		log.Fatal("Error: context expired")
		return
	case <-doneC:
	}

	var messages []string
	printed := make(map[string]struct{})

	// check if specific order is set
	if len(ctx.Config.Order) > 0 {
		for _, section := range ctx.Config.Order {
			if result, ok := sectionResults[section]; ok {
				messages = append(messages, result)
				printed[section] = struct{}{}
			}
		}
	}

	// print the rest
	for _, section := range message.Message {
		if !section.Enabled(ctx) {
			continue
		}
		// skip if already printed
		if _, ok := printed[section.String()]; ok {
			continue
		}
		if msg, ok := sectionResults[section.String()]; ok {
			messages = append(messages, msg)
		}
	}

	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, messages...))
}
