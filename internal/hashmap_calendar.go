package internal

import "time"

type HashMapCalendar struct {}

func InitHashMapCalendar() HashMapCalendar {
	return HashMapCalendar{}
}

// REFACTOR
func (hc HashMapCalendar) InitModel(gitlog map[string][]string) Model {

	datesFromGitLog := parseDataFromGitlog(gitlog)
	parsedData := hc.ParseCalToView(datesFromGitLog)

	todayX, todayY := hc.GetDateIndex(time.Now())

	return Model{
		selectedX: todayX,
		selectedY: todayY,
		calData:   datesFromGitLog,
		viewData:  parsedData,
		hc:        hc,
	}
}

/**
 * Returns the index of the date in the viewData array
 * x is the number of weeks ago
 * y is the day of the week
 */
func (hc HashMapCalendar) GetDateIndex(date time.Time) (int, int) {
	// Max index - number of weeks ago
	x, y := dateToIndex(date)
	return x, y
}

// MOVE to utils
func dateToIndex(date time.Time) (int, int) {
	// Max index - number of weeks ago
	x := 51 - WeeksAgo(date)
	y := WeekDay(date)

	return x, y
}

// TEST
/**
 * Returns the date at the index in the viewData array
 * 		x is the number of weeks ago
 * 		y is the day of the week
 */
func (hc HashMapCalendar) GetIndexDate(x int, y int) time.Time {
	// compare the x,y to today and subtract
	today := time.Now()
	todayX, todayY := hc.GetDateIndex(today)

	diffX := todayX - x
	diffY := todayY - y

	diffDays := diffX*7 + diffY

	targetDate := today.AddDate(0, 0, -diffDays)
	return targetDate
}

/**
 * Parse the calData into a viewData array
 * calData is a list of CalDataPoint
 * viewData is a 2D array of ViewDataPoint
 */
func (hc HashMapCalendar) ParseCalToView(calData []CalDataPoint) [52][7]ViewDataPoint {

	viewData := [52][7]ViewDataPoint{}

	for _, calDataPoint := range calData {
		x, y := hc.GetDateIndex(calDataPoint.Date)
		// asign
		if x > -1 && y > -1 &&
			x < 52 && y < 7 {
			viewData[x][y].actual += calDataPoint.CommitCount
			viewData[x][y].commits = calDataPoint.CommitMessages
		}
	}

	viewData = normalizeViewData(viewData)
	return viewData
}

func normalizeViewData(data [52][7]ViewDataPoint) [52][7]ViewDataPoint {
	min, max := MinAndMax(data)
	// Normalize the data
	for i, row := range data {
		for j, val := range row {
			data[i][j].normalized = (val.actual - min) / (max - min)
		}
	}
	return data
}

func MinAndMax(data [52][7]ViewDataPoint) (float64, float64) {
	var min float64
	var max float64

	// Find min/max
	min = data[0][0].actual
	max = data[0][0].actual

	// find min and max
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

	return min, max

}
