package calendar

import (
	"encoding/json"
	"fmt"
	"os"
)

// from 1900 only
func getStartDay(year int) int{
	i:= year % 100
	var yearCode int = (i + (i/4)) % 7
	var centuryCode int
	if (year / 100) == 19{
		centuryCode = 0
	} else if (year / 100) == 20{
		centuryCode = 6
	}
	var sum int = yearCode + 0 + centuryCode + 1
	if year % 4 == 0{
		sum -= 1
	}
	return (sum % 7)
}

type SerializedYear struct{
	Year int
	Bin []byte
}

type SerializedCalendar struct{
	Start int
	End int
	Events []byte
}

type Calendar struct {
	years []*Year
}

func createCalendar(startYear int, endYear int) *Calendar {
	calendar := new(Calendar)
	for i := startYear; i <= endYear; i++ {
		calendar.years = append(calendar.years, createYear(i, getStartDay(i)))
	}

	return calendar
}

func (calendar *Calendar) getYearsRange() (int, int) {
	var size int = len(calendar.years)
	return calendar.years[0].year, calendar.years[size-1].year
}

func (calendar *Calendar) getYear(year int) *Year {
	start := calendar.years[0].year
	return calendar.years[year-start]
}

func (calendar *Calendar) getDay(year int, month int, day int) *Day {
	return calendar.getYear(year).getDay(month, day)
}

func (calendar *Calendar) expandYears(years int){
	_ , end := calendar.getYearsRange()
	for i :=0; i <= years; i++{
		end++
		calendar.years = append(calendar.years, createYear(end, getStartDay(end)))
	}
	
}

func (calendar *Calendar) saveToDisk(path string) {
	var err error
	var bin []byte
	var file *os.File
	bin, err = calendar.serialize()
	if err == nil{
		file, err = os.Open(path)
		defer file.Close()
		if err != nil{
			if os.IsNotExist(err){
				file, err = os.Create(path)
			}
		}
	}
	if err == nil{
		file.Write(bin)
	}
}

func (calendar *Calendar) serialize() ([]byte, error){
	var yearsBin [][]byte
	var err error
	var current []byte

	for _, year := range calendar.years{
		if year.numOfEvents() > 0{
			current, err = year.serialize()
			if err == nil{
				current, err = json.Marshal(SerializedYear{year.year, current})
				
				if err == nil{
					yearsBin = append(yearsBin, current)
				}
			}
		}
	}
	var start, end int = calendar.getYearsRange()
	current, err =  json.Marshal(yearsBin)

	return json.Marshal(SerializedCalendar{start, end, current})
}

func deserialize(bin []byte, reminder reminderAlert) *Calendar{
	var years [][]byte
	var yearContainer SerializedYear
	var calendarContainer SerializedCalendar
	var err error
	var calendar *Calendar;

	err = json.Unmarshal(bin, &calendarContainer)
	if err == nil{
		calendar = createCalendar(calendarContainer.Start, calendarContainer.End)
		err = json.Unmarshal(calendarContainer.Events, &years)
		for _, year := range years{
			err = json.Unmarshal(year, &yearContainer)
			fmt.Println(yearContainer.Year)
			if err == nil{
				calendar.getYear(yearContainer.Year).setEventsFromBin(yearContainer.Bin, reminder)
			}
		}
	}

	return calendar
}

func loadFromFile(path string, reminder reminderAlert) *Calendar{
	var bin []byte
	var err error

	bin, err = os.ReadFile(path)
	if err == nil{
		return deserialize(bin, reminder)
	}
	return nil
}
