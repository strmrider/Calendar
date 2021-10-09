package calendar

import (
	"encoding/json"
	"os"
	"sync"
)

const ( 
	minYear = 1900
	maxYear = 2099
)

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
	lock sync.Locker
}

type eventFinished func()

func CreateCalendar(startYear int, endYear int) *Calendar {
	if startYear < minYear || startYear > maxYear || 
		endYear < minYear || endYear > maxYear{
		panic("Invalid year: years must be between 1900 to 2099")
	}

	calendar := new(Calendar)
	for i := startYear; i <= endYear; i++ {
		calendar.years = append(calendar.years, createYear(i, getStartDay(i)))
	}
	return calendar
}

func (calendar *Calendar) IsYearExist(year int) bool{
	var start, end int = calendar.GetRange()
	return year >= start && year <= end
}

func (calendar *Calendar) GetRange() (int, int) {
	var size int = len(calendar.years)
	return calendar.years[0].year, calendar.years[size-1].year
}

func (calendar *Calendar) GetYear(year int) *Year {
	if !calendar.IsYearExist(year){
		panic("Year doesn't exist in calendar")
	}
	var start, _ int = calendar.GetRange()
	return calendar.years[year-start]
}

func (calendar *Calendar) GetDay(year int, month int, day int) *Day {
	if !calendar.IsYearExist(year){
		panic("Year doesn't exist in calendar")
	} else if month < 0 || month > 11{
		panic("Invalid month: month's range must be between 0 to 11")
	}
	return calendar.GetYear(year).GetDay(month, day)
}

func (calendar *Calendar) ExpandYears(years int){
	_ , end := calendar.GetRange()
	if end + years > maxYear{
		panic("Expansion exceeds maximal year")
	}
	for i :=0; i <= years; i++{
		end++
		calendar.years = append(calendar.years, createYear(end, getStartDay(end)))
	}
}

func (calendar *Calendar) SaveToDisk(path string) {
	var err error
	var bin []byte
	var file *os.File
	bin, err = calendar.Serialize()
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

func (calendar *Calendar) ShareOnNetwork(port int){
	if port <=0 || port > 65535{
		panic("Invlid port")
	}
	var network *Network = createNetwork(calendar.calendarToNetwork)
	syncwg.operate(wgAdd)
	go network.Listen(port)
}

func (calendar *Calendar) Serialize() ([]byte, error){
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
	var start, end int = calendar.GetRange()
	current, err =  json.Marshal(yearsBin)

	return json.Marshal(SerializedCalendar{start, end, current})
}

func Deserialize(bin []byte, reminder reminderAlert) *Calendar{
	if len(bin) == 0{
		panic("Cannot decode empty bytes array")
	}
	var years [][]byte
	var yearContainer SerializedYear
	var calendarContainer SerializedCalendar
	var err error
	var calendar *Calendar;

	err = json.Unmarshal(bin, &calendarContainer)
	if err == nil{
		calendar = CreateCalendar(calendarContainer.Start, calendarContainer.End)
		err = json.Unmarshal(calendarContainer.Events, &years)
		for _, year := range years{
			err = json.Unmarshal(year, &yearContainer)
			if err == nil{
				calendar.GetYear(yearContainer.Year).setEventsFromBin(yearContainer.Bin, reminder)
			}
		}
	}

	return calendar
}

func LoadFromFile(path string, reminder reminderAlert) *Calendar{
	var bin []byte
	var err error

	bin, err = os.ReadFile(path)
	if err == nil{
		return Deserialize(bin, reminder)
	}else if os.IsNotExist(err) {
		panic("Path does not exist")
	}
	return nil
}

func (calendar *Calendar) calendarToNetwork() []byte{
	bin, err := calendar.Serialize()
	if err == nil{
		return bin
	}
	return nil
}