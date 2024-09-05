package main

import (
	"fmt"
	"testing"
	"strings"
	"time"
)


func TestGitlog(t *testing.T) {

	m , error :=  getLogMap()
	if error != nil {
		t.Error("Error in getLogMap")
	}
	// fmt.Println(m)
	for key, val := range m {
		fmt.Println(key)
		fmt.Println(val)
		// fmt.Println(len(val))
		// fmt.Println(float64(len(val)))
	}
	datesFromGitLog := []CalDataPoint{}
	for key, value := range m {
		date := strings.Split(key, "-")
		year, month, day := date[0], date[1], date[2]
		
		// fmt.Printf("Year: %s, Month: %s, Day: %s\n", year, month, day)
		
		d := time.Date(parseInt(year), time.Month(parseInt(month)), parseInt(day),0,0,0,0, time.UTC)
		// fmt.Printf("Date: %s\n", d)
		datesFromGitLog = append(datesFromGitLog, CalDataPoint{
		  Date: d,
		  Value: float64(len(value)),
		  CommitMessages: value,
		},
	  )     
	}

	fmt.Println(datesFromGitLog)
	t.Log("Test Passed")
}


func testGetDateIndex(t *testing.T) {
	m , error :=  getLogMap()	
	if error != nil {
		t.Error("Error in getLogMap")
	}

	datesFromGitLog := []CalDataPoint{}
	for key, value := range m {
		date := strings.Split(key, "-")
		year, month, day := date[0], date[1], date[2]
		d := time.Date(parseInt(year), time.Month(parseInt(month)), parseInt(day),0,0,0,0, time.UTC)
		datesFromGitLog = append(datesFromGitLog, CalDataPoint{
		  Date: d.Local()  ,
		  // Value: float64(len(value)/10),
		  Value: float64(1),
		  CommitMessages: value,
		},
	  )     
	}

	for _ , val := range datesFromGitLog {
		x, y := getDateIndex(val.Date)
		fmt.Println(x)
		fmt.Println("-")
		fmt.Println(y)
	}


	t.Log("Test Passed testGetDateIndex")

}

func testDateNow(t *testing.T){
	today := time.Now()
	fmt.Println(today)
	t.Log("Test Passed testDateNow")
}