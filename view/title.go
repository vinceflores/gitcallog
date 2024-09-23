package view

import (
	"time"

	b "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func Title(m b.Model, theTime time.Time)string{
	title, _ := glamour.Render(theTime.Format("# Monday, January 02, 2006"), "dark")
	return title
}

