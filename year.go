package calendar

import (
	"encoding/json"
	"fmt"
)

type SerializedMonth struct{
	Month int
	Bin []byte
}

type Year struct {
	year   int
	months [12]*Month
}

func createYear(year int, firstDay int) *Year {
	newYear := new(Year)
	newYear.year = year
	var day int = firstDay
	for i := 0; i < 12; i++ {
		//fmt.Println(day)
		newYear.months[i] = createMonth(i, day, newYear.isLeap())
		day = newYear.months[i].getLastDayName()
		if day > Saturday {
			day = Sunday
		}
	}
	return newYear
}

func (year *Year) getMonth(month int) *Month {
	return year.months[month]
}

func (year *Year) getDay(month int, day int) *Day {
	return year.months[month].getDay(day)
}

func (year *Year) numOfEvents() int {
	var sum int = 0
	for _, month := range year.months{
		sum += month.numOfEvents()
	}
	return sum
}

func (year *Year) isLeap() bool {
	return year.year%4 == 0
}

func (year *Year) print() {
	fmt.Println(year.year)
	for _, month:= range year.months{
		month.print()
	}
}

func (year *Year) serialize() ([]byte, error){
	var monthsBins [][]byte
	var err error
	var current []byte
	var i int = 0

	for _, month := range year.months{
		if month.numOfEvents() > 0 {
			current, err = month.serialize()
			if err == nil{
				current, err = json.Marshal(SerializedMonth{i, current})
				if err == nil{

					monthsBins = append(monthsBins, current)
				}
			}
		}
		i++
	}
	return json.Marshal(monthsBins)
}

func (year *Year) setEventsFromBin(bin []byte, reminder reminderAlert) {
	var err error
	var months [][]byte
	var container SerializedMonth
	err = json.Unmarshal(bin, &months)
	if err == nil{
		for _, monthBin := range months{
			err = json.Unmarshal(monthBin, &container)
			if err == nil{
				year.months[container.Month].setEventsFromBin(container.Bin, reminder)
			}
		}
	}
}