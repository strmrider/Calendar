package calendar

import (
	"encoding/json"
	"fmt"
)

const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func getDayName(day int) string {
	switch day {
	case Sunday:
		return "Sunday"
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	default:
		return "Undefined"
	}
}

type Day struct {
	index     int
	name      string
	events    []Event
	
}

func createDay(dayNumber int) *Day {
	day := new(Day)
	day.index = dayNumber + 1
	day.name = getDayName(dayNumber)
	return day
}

func (day *Day) GetDay() string {
	return day.name
}

func (day *Day) AddEvent(event *Event) {
	day.events = append(day.events, *event)
}

func (day *Day) RemoveEvent(index int) {
	if index < 0 || index > day.NumOfEvents()-1{
		panic("Inavlid index")
	} else {
		day.events = append(day.events[:index], day.events[index+1:]...)
	}
}

func (day *Day) GetEvent(index int) *Event {
	if index < 0 || index > day.NumOfEvents()-1{
		panic("Inavlid index")
	}
	return &day.events[index]
}

func (day *Day) ClearEvents() {
	for _, event := range day.events {
		event.StopTimer()
	}
	day.events = nil
}

func (day *Day) NumOfEvents() int {
	return len(day.events)
}

func (day *Day) Summary() string {
	return fmt.Sprintf("%s: %d events", day.name, len(day.events))
}

func (day *Day) Print() {
	fmt.Printf("%s: %d events\n", day.name, len(day.events))
	for _, event := range day.events{
		event.Print()
	}
}

func (day *Day) serialize() ([]byte, error){
	var eventSerials [][]byte
	var currentEvent []byte
	var err error
	for _, event := range day.events{
		currentEvent, err = event.serialize()
		if err == nil {
		eventSerials = append(eventSerials, currentEvent)
		}
	}
	return json.Marshal(eventSerials)
}

func (day *Day) setEventsFromBin(bin []byte, handler reminderAlert){
	var eventsBin [][]byte
	var err error
	var fields []string
	err = json.Unmarshal(bin, &eventsBin)
	if err == nil{
		for _, event := range eventsBin{
			err = json.Unmarshal(event, &fields)
			if err == nil{
				day.events = append(day.events, *CreateEvent(fields[0], fields[1], fields[2], handler))
			}
		}
	}
}