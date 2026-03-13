package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"os"
)

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
