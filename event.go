package calendar

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	title    string
	content  string
	creation time.Time
	alertTime time.Time
	reminder reminderAlert
	timer *time.Timer
}

type reminderAlert func(event *Event)

func createEvent(title, content, date string, reminder reminderAlert) *Event{
	event := new(Event)
	event.title = title
	event.content = content
	event.creation = time.Now()
	event.reminder = reminder
	t, err:= time.Parse(time.RFC3339, date)
	if (err == nil){
		event.alertTime = t
	}

	return event
}

func (event *Event) runTimer(){
	diff := event.alertTime.Sub(time.Now())
	fmt.Println(event.alertTime)
	fmt.Println(time.Now())
	diff -= 3
	fmt.Println(diff)
	if diff < 0{
		diff *= (-1)
	}
	event.timer = time.NewTimer(diff)
	go func(event *Event){
		<-event.timer.C
		if !event.timer.Stop(){
			event.reminder(event)
		}
	}(event)
}

func (event *Event) stopTimer(){
	event.timer.Stop()
}

func (event *Event) print(){
	fmt.Printf("%s\t%s\n", event.title, event.alertTime.Format(time.RFC1123))
	fmt.Println(event.content)
	fmt.Println()
}

func (event *Event) serialize() ([]byte, error){
	return json.Marshal([]string{event.title, event.content, event.alertTime.Format(time.RFC1123)})
}