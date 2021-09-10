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

func (day *Day) getDay() string {
	return day.name
}

func (day *Day) addEvent(event *Event) {
	day.events = append(day.events, *event)
}

func (day *Day) removeEvent(index int) {
	day.events = append(day.events[:index], day.events[index+1:]...)
}

func (day *Day) getEvent(index int) *Event {
	return &day.events[index]
}

func (day *Day) clearEvents(index int) {
	for _, event := range day.events {
		event.stopTimer()
	}
	day.events = nil
}

func (day *Day) numOfEvents() int {
	return len(day.events)
}

func (day *Day) summary() string {
	return fmt.Sprintf("%s: %d events", day.name, len(day.events))
}

func (day *Day) print() {
	fmt.Printf("%s: %d events\n", day.name, len(day.events))
	for _, event := range day.events{
		event.print()
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
				day.events = append(day.events, *createEvent(fields[0], fields[1], fields[2], handler))
			}
		}
	}
}