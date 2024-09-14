package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	selectedX int
	selectedY int
	calData   []CalDataPoint
	viewData  [52][7]viewDataPoint // Hardcoded to one year for now
}

type CalDataPoint struct {
	Date                 time.Time
	CommitMessages       []string
	Value                float64
}

type viewDataPoint struct {
	actual     float64
	normalized float64
	commits 	[]string
}

var scaleColors = []string{
	"#161b22", // Less
	"#0e4429",
	"#006d32",
	"#26a641",
	"#39d353", // - More
}

/******************************************/


// TEST
func weeksAgo(date time.Time) int {
	today := truncateToDate(time.Now())
	thisWeek := today.AddDate(0, 0, -int(today.Weekday())) // Most recent Sunday

	compareDate := date								  // truncate to date
	compareWeek := compareDate.AddDate(0, 0, -int(compareDate.Weekday())) // get teh previews week

	result := thisWeek.Sub(compareWeek).Hours() / 24 / 7 			  // get the number of weeks between the two dates	
	return int(result)
}

// TEST
// sets the datet to 00:00:00
func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Local().Year(), t.Local().Month(), t.Local().Day(), 0, 0, 0, 0, t.Local().Location())
}

// TEST
/**
	 * Returns the index of the date in the viewData array
	 * x is the number of weeks ago
	 * y is the day of the week
*/ 
func getDateIndex(date time.Time) (int, int) {
	// Max index - number of weeks ago
	x := 51 - weeksAgo(date) 

	y := int(date.Weekday()) 

	return x, y
}

// TEST
/**
	 * Returns the date at the index in the viewData array
	 * 		x is the number of weeks ago
	 * 		y is the day of the week
*/ 
func getIndexDate(x int, y int) time.Time {
	// compare the x,y to today and subtract
	today := time.Now()
	todayX, todayY := getDateIndex(today)

	diffX := todayX - x 
	diffY := todayY - y

	diffDays := diffX*7 + diffY

	targetDate := today.AddDate(0, 0, -diffDays)
	return targetDate
}

func parseCalToView(calData []CalDataPoint) [52][7]viewDataPoint {

	viewData := [52][7]viewDataPoint{}

	for _, calDataPoint := range calData {
		x, y := getDateIndex(calDataPoint.Date)
		// asign 
		if x > -1 && y > -1 &&
			x < 52 && y < 7 {
				viewData[x][y].actual += calDataPoint.Value
				viewData[x][y].commits = calDataPoint.CommitMessages
		}
	}

	viewData = normalizeViewData(viewData)
	return viewData
}

func normalizeViewData(data [52][7]viewDataPoint) [52][7]viewDataPoint {
	var min float64
	var max float64

	// Find min/max
	min = data[0][0].actual
	max = data[0][0].actual

	for _, row := range data {
		for _, val := range row {

			if val.actual < min {
				min = val.actual
			}
			if val.actual > max {
				max = val.actual
			}
		}
	}

	// Normalize the data
	for i, row := range data {
		for j, val := range row {
			data[i][j].normalized = (val.actual - min) / (max - min)
		}
	}
	return data
}

func parseInt(s string) int {
	i , _ :=  strconv.Atoi(s)
	return int(i)
}

func getScaleColor(value float64) string {
	const numColors = 5
	// Assume it's normalized between 0.0-1.0
	const max = 1.0
	// const min = 0.0
	norm := (value/max)*(numColors-1) 
	if value > 0 && value < 0.5 {
		return scaleColors[0.5 * (numColors -1) ]
	}
	return scaleColors[int(norm)]
}
  
func getLogMap() (map[string][]string , error)  {
  // Define the git command with arguments
  cmd := exec.Command("git", "log", "--pretty=format:%ad/%s", "--date=short")
  // Run the command and capture the output
  output, err:= cmd.Output()
  if err != nil {
      fmt.Printf("Error: %s\n", err)
      return nil , err
  }
  // fmt.Printf("Output: \n%s\n", output)  

  output_str_arr  := strings.Split(string(output), "\n")

  // make a hashmap
  m := make(map[string][]string)  
  for _, s := range output_str_arr {
    temp := strings.Split(s, "/")
    // iterate through temp and append to hashmap
    for i , v := range temp {
      if i == 0 {
        continue
      }else{
        m[string(temp[0])] = append(m[string(temp[0])], string(v))
      }
    }
  }

  // DELETE
  // fmt.Println("Printing the hashmap")
  // for key , value:= range m {
  //   fmt.Println(key)
  //   for _, v := range value {
  //     fmt.Println(v)
  //   }
  // }

  return m , nil
}

// func (m *Model) addCalData(date time.Time, commits []string ,  val float64) {
// 	// Create new cal data point and add to cal data
// 	newPoint := CalDataPoint{date, commits,  val}
// 	m.calData = append(m.calData, newPoint)
// }
/***************MODEL***************************/
func (m Model) Init() tea.Cmd{
  return nil
}

func InitModel( gitlog map[string][]string, ) Model {
  // make calData
  todayX, todayY := getDateIndex(time.Now())
  datesFromGitLog := []CalDataPoint{}
  for key, value := range gitlog {
	  date := strings.Split(key, "-")
	  year, month, day := date[0], date[1], date[2]
	  d := time.Date(parseInt(year), time.Month(parseInt(month)), parseInt(day),0,0,0,0, time.UTC)
	  datesFromGitLog = append(datesFromGitLog, CalDataPoint{
		Date: d,
		Value: float64(len(value)),
		CommitMessages: value,
	  },
	)     
  }

	parsedData := parseCalToView(datesFromGitLog)

  return Model{
    selectedX: todayX,
    selectedY: todayY,
    calData:   datesFromGitLog,
    viewData: parsedData,
  }
}

// updates view based on user input
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: ignore if not focused
	// if !m.focus { return m, nil }

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
		// case "enter", " ":
		// 	// Hard coded to add a new entry with value `1.0`
		// 	m.addCalData(
		// 		getIndexDate(m.selectedX, m.selectedY),
		// 		m.calData[m.selectedX].CommitMessages,
		// 		1.0)
		// 	m.viewData = parseCalToView(m.calData)

		}
	}
	return m, nil
}

func (m Model) View() string{
	theTime := getIndexDate(m.selectedX, m.selectedY) // time.Now()
	// theme := "dark"
	// title, _ := glamour.Render(theTime.Format("Monday, Januoary 02, 2006"), "#0f0KJJKJ")
	title, _ := glamour.Render(theTime.Format("# Monday, January 02, 2006"), "dark")
	s := title
	// s += "  theme: " +  theme + "\n"
	selectedDetail := "  Commits: " +
		fmt.Sprint(m.viewData[m.selectedX][m.selectedY].actual) +
		" normalized: " +
		fmt.Sprint(m.viewData[m.selectedX][m.selectedY].normalized) +
		"\n\n"

	s += selectedDetail
	s += "  Press 'q' to quit\n\n"	
	// styles 

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	boxStyle := lipgloss.NewStyle().
		PaddingRight(1).
		Foreground(lipgloss.Color(scaleColors[2]))

	boxSelectedStyle := boxStyle.Copy().
		Background(lipgloss.Color("#9999ff")).
		Foreground(lipgloss.Color(scaleColors[0]))

	// Month Labels
	var currMonth time.Month
	s += "  "
	for j := 0; j < 52; j++ {
		// Check the last day of the week for that column
		jMonth := getIndexDate(j, 6).Month()

		if currMonth != jMonth {
			currMonth = jMonth
			s += labelStyle.Render(getIndexDate(j, 6).Format("Jan") + " ")

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
  s += "Use arrow keys to move \nor type;\n"

  s += "h - left \n"
  s += "j - down \n"
  s += "k - up \n"
  s += "l - right \n"


	// // list commit messages
	commits := m.viewData[m.selectedX][m.selectedY].commits
	if(len(commits) > 0 ) {
		s+= "\nCOMMITS\n"
		s += "  - " + strings.Join(commits, "\n  - ") + "\n\n"
	}


	return s
}

func main(){
  m , err := getLogMap()
  if err != nil {
	fmt.Println("Error getting git log \n Make sure you are in a git repo")
    return 
  }
  p := tea.NewProgram(InitModel(m))
  
  if _, err := p.Run(); err != nil {
    fmt.Printf("Alas, there's been an error: %v", err)
    os.Exit(1)
  }
}
