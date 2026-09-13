package main

import (
	"log"

	"lazytrees/internal/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if _, err := tea.NewProgram(ui.Bootstrap()).Run(); err != nil {
		log.Fatal(err)
	}
}
