package calendar

import (
	"encoding/json"
	"fmt"
)

const (
	January = iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func getMonthName(month int) string {
	switch month {
	case January:
		return "January"
	case February:
		return "February"
	case March:
		return "March"
	case April:
		return "April"
	case May:
		return "May"
	case June:
		return "June"
	case July:
		return "July"
	case August:
		return "August"
	case September:
		return "September"
	case October:
		return "October"
	case November:
		return "November"
	case December:
		return "December"
	default:
		return "Undefined"
	}
}

func getMonthDays(month int, leap bool) int {
	switch month {
	case January:
		return 31
	case February:
		if leap {
			return 29
		} else {
			return 28
		}
	case March:
		return 31
	case April:
		return 30
	case May:
		return 31
	case June:
		return 30
	case July:
		return 31
	case August:
		return 31
	case September:
		return 30
	case October:
		return 31
	case November:
		return 30
	case December:
		return 31
	default:
		return 0
	}
}

type SerializedDay struct{
	Day int
	Bin []byte
}

type Month struct {
	index int
	name  string
	days  []*Day
}

func createMonth(monthType int, firstDay int, leap bool) *Month {
	month := new(Month)
	month.index = monthType + 1
	month.name = getMonthName(monthType)
	var j int = firstDay
	for i := 0; i < getMonthDays(monthType, leap); i++ {
		//fmt.Println(getDayName(j))
		month.days = append(month.days, createDay(j))
		j++
		if j > Saturday {
			j = Sunday
		}
	}
	return month
}

func (month *Month) GetMonth() string {
	return month.name
}

func (month *Month) NumOfDays() int {
	return len(month.days)
}

func (month *Month) GetDay(day int) *Day {
	if day < 0 || day > month.NumOfDays(){
		panic("Day doesn't exist")
	}
	return month.days[day-1]
}

func (month *Month) getLastDayName() int {
	return month.days[month.NumOfDays()-1].index
}

func (month *Month) NumOfEvents() int {
	var sum int = 0
	for _, day := range month.days{
		sum += day.NumOfEvents()
	}
	return sum
}

func (month *Month) Summary() {
	fmt.Printf("%s\n", month.GetMonth())
	var i int = 0
	for _, day := range month.days{
		fmt.Printf("%d %s\n", i+1, day.Summary())
		i++
	}
}

func (month *Month) serialize() ([]byte, error){
	var daySerials [][]byte
	var err error
	var currentDay []byte
	var i int = 0
	for _, day := range month.days{
		if day.NumOfEvents() > 0{
			currentDay, err = day.serialize()
			if err == nil{
				currentDay, err = json.Marshal( SerializedDay{i, currentDay})
				daySerials = append(daySerials, currentDay)
			}
		}
		i++
	}
	return json.Marshal(daySerials)
}

func (month *Month) setEventsFromBin(bin []byte, handler reminderAlert){
	var days [][]byte
	var err error
	var container SerializedDay
	err = json.Unmarshal(bin, &days)
	if err == nil{
		for _, dayBin := range days{
			err = json.Unmarshal(dayBin, &container)
			if err == nil{
				month.GetDay(container.Day+1).setEventsFromBin(container.Bin, handler)
			}
		}
	}
}

func (month *Month) Print(){
	fmt.Printf("\n%s\n", month.GetMonth())
	fmt.Println("Sun\tMon\tTue\tWed\tThu\tFri\tSat")
	fmt.Println("====================================================")
	for i:=0; i<month.GetDay(1).index-1; i++ {
		fmt.Printf("\t")
	}

	for i:=1; i<=month.NumOfDays(); i++{
		if month.GetDay(i).NumOfEvents() > 0{
			fmt.Printf("*")
		}
		fmt.Printf("%d\t", i)
		if month.GetDay(i).index-1 == Saturday{
			fmt.Println()
			fmt.Println("====================================================")
		}
	}
	if month.GetDay(month.NumOfDays() -1 ).index-1 != Saturday{
		fmt.Println()
	}
}