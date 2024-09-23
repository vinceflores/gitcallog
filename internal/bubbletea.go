package internal

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	view "github.com/vinceflores/gitcallog/view"
)

type Model struct {
	selectedX int
	selectedY int
	calData   []CalDataPoint
	viewData  [52][7]ViewDataPoint // Hardcoded to one year for now

	hc HashMapCalendar
}

// updates view based on user input
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.selectedY > 0 {
				m.selectedY--
			} else if m.selectedY == 0 && m.selectedX > 0 {
				// Scroll to the end of the previous week if on Sunday
				// and not at the beginning of the calendar.
				m.selectedY = 6
				m.selectedX--
			}

		case "down", "j":
			// Don't allow user to scroll beyond today
			if m.selectedY < 6 &&
				(m.selectedX != 51 ||
					m.selectedY < int(time.Now().Weekday())) {
				m.selectedY++
			} else if m.selectedY == 6 && m.selectedX != 51 {
				// Scroll to the beginning of next week if on Saturday
				// and not at the end of the calendar.
				m.selectedY = 0
				m.selectedX++
			}
		case "right", "l":
			// Don't allow users to scroll beyond today from the previous column
			if m.selectedX < 50 ||
				(m.selectedX == 50 &&
					m.selectedY <= int(time.Now().Weekday())) {
				m.selectedX++
			}
		case "left", "h":
			if m.selectedX > 0 {
				m.selectedX--
			}

		}
	}
	return m, nil
}

func (m Model) View() string {
	theTime := m.hc.GetIndexDate(m.selectedX, m.selectedY) // time.Now()
	s := "v2"
	s += view.Title(m, theTime)
	// selectedDetail := "  Commits: " +
	// 	fmt.Sprint(m.viewData[m.selectedX][m.selectedY].actual) +
	// 	" normalized: " +
	// 	fmt.Sprint(m.viewData[m.selectedX][m.selectedY].normalized) +
	// 	"\n\n"

	// s += selectedDetail
	
	// styles

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	boxStyle := lipgloss.NewStyle().
		PaddingRight(1).
		Foreground(lipgloss.Color(ScaleColors[2]))
	boxSelectedStyle := boxStyle.Copy().
		Background(lipgloss.Color("#9999ff")).
		Foreground(lipgloss.Color(ScaleColors[0]))

	// Month Labels
	var currMonth time.Month
	s += "  "

	for j := 0; j < 52; j++ {
		// Check the last day of the week for that column
		jMonth := m.hc.GetIndexDate(j, 6).Month()

		if currMonth != jMonth {
			currMonth = jMonth
			s += labelStyle.Render(m.hc.GetIndexDate(j, 6).Format("Jan") + " ")

			// Skip the length of the label we just added
			j += 1
		} else {
			s += "  "
		}
	}

	s += "\n"

	// Days of the week and calendar
	for j := 0; j < 7; j++ {
		// Add day of week labels
		switch j {
		case 0:
			s += labelStyle.Render("S ")
		case 1:
			s += labelStyle.Render("M ")
		case 2:
			s += labelStyle.Render("T ")
		case 3:
			s += labelStyle.Render("W ")
		case 4:
			s += labelStyle.Render("T ")
		case 5:
			s += labelStyle.Render("F ")
		case 6:
			s += labelStyle.Render("S ")
		}

		// Add calendar days
		for i := 0; i < 52; i++ {
			// Selected Item
			if m.selectedX == i && m.selectedY == j {
				s += boxSelectedStyle.Copy().Foreground(
					lipgloss.Color(
						getScaleColor(
							m.viewData[i][j].normalized))).
					Render("■")
			} else if i == 51 &&
				j > int(time.Now().Weekday()) {

				// In the future
				s += boxStyle.Render(" ")
			} else {

				// Not Selected Item and not in the future
				s += boxStyle.Copy().
					Foreground(
						lipgloss.Color(
							getScaleColor(
								m.viewData[i][j].normalized))).
					Render("■")
			}
		}
		s += "\n"
	}
	s += "\n\n"

	s += "INSTRUCTIONS\n"
	s += "Press 'q' to quit\n\n"
	s += "Use arrow keys to move \nor type any from;\n"
	s += "[{h:left}, {j:down}, {k:up}, {l:right}]\n"

	// // list commit messages
	commits := m.viewData[m.selectedX][m.selectedY].commits

	if len(commits) > 0 {
		s += "\nCOMMITS\n"
		s += "  - " + strings.Join(commits, "\n  - ") + "\n\n"
	}

	return s
}

func (m Model) Init() tea.Cmd {
	return nil
}

/****************************HELPER FUNCTIONS************************************/

func localDate(date string) time.Time {
	d := strings.Split(date, "-")
	year, month, day := d[0], d[1], d[2]
	return time.Date(ParseInt(year), time.Month(ParseInt(month)), ParseInt(day), 0, 0, 0, 0, time.UTC)
}

func parseDataFromGitlog(gitlog map[string][]string) []CalDataPoint {
	datesFromGitLog := []CalDataPoint{}

	for date, value := range gitlog {
		d := localDate(date)
		datesFromGitLog = append(datesFromGitLog, CalDataPoint{
			Date:           d,
			CommitCount:    float64(len(value)),
			CommitMessages: value,
		},
		)
	}
	return datesFromGitLog
}
