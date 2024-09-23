package internal

import "time"

func TruncateToDate(t time.Time) time.Time {
	return time.Date(t.Local().Year(), t.Local().Month(), t.Local().Day(), 0, 0, 0, 0, t.Local().Location())
}

// MOVE to utils
func WeeksAgo(date time.Time) int {
	today := TruncateToDate(time.Now())
	thisWeek := today.AddDate(0, 0, -int(today.Weekday())) // Most recent Sunday

	compareDate := date                                                   // truncate to date
	compareWeek := compareDate.AddDate(0, 0, -int(compareDate.Weekday())) // get teh previews week

	result := thisWeek.Sub(compareWeek).Hours() / 24 / 7 // get the number of weeks between the two dates
	return int(result)
}

// MOVE to utils
func WeekDay(date time.Time) int {
	return int(date.Weekday())
}
