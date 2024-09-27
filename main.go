package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	internal "github.com/vinceflores/gitcallog/internal"
)

func main() {
	data, err := internal.GetLogMap()
	if err != nil {
		fmt.Println("Error getting git log \n Make sure you are in a git repo")
		return
	}

	hashMap := internal.InitHashMapCalendar(data)

	model := internal.InitModel(hashMap)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
