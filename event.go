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

func twoDigits(number int) string{
	if number < 10{
		return fmt.Sprintf("0%d", number)
	}
	return fmt.Sprintf("%d", number)
}

func BuildTime(year, month, day, hour, minute, second int) string{
	return fmt.Sprintf("%s, %s %s %d %s:%s:%s MST", 
						getDayName(day)[:3], twoDigits(day), getMonthName(month)[:3], 
						year, twoDigits(hour), twoDigits(minute), twoDigits(second))
}

func CreateEvent(title, content, date string, reminder reminderAlert) *Event{
	event := new(Event)
	event.title = title
	event.content = content
	event.creation = time.Now()
	event.reminder = reminder
	t, err:= time.Parse(time.RFC1123, date)
	if (err == nil){
		event.alertTime = t
	}

	return event
}

/*Getters and setters*/

func (event *Event) Title() string{
	return event.title
}

func (event *Event) Content() string{
	return event.title
}

func (event *Event) SetTitle(title string){
	event.title = title
}

func (event *Event) SetContent(content string){
	event.content = content
}

func (event *Event) ResetAlertTime(date time.Time){
	event.StopTimer()
	event.alertTime = date
}

func (event *Event) RunTimer(){
	diff := event.alertTime.Sub(time.Now())
	fmt.Println(event.alertTime)
	fmt.Println(time.Now())
	diff -= 3
	fmt.Println(diff)
	if diff < 0{
		diff *= (-1)
	}
	event.timer = time.NewTimer(diff)
	syncwg.operate(wgAdd)
	go func(event *Event){
		<-event.timer.C
		if !event.timer.Stop(){
			syncwg.operate(wgDone)
			event.reminder(event)
		}
	}(event)
}

func (event *Event) StopTimer(){
	event.timer.Stop()
}

func (event *Event) Print(){
	fmt.Printf("%s\t%s\n", event.title, event.alertTime.Format(time.RFC1123))
	fmt.Println(event.content)
	fmt.Println()
}

func (event *Event) serialize() ([]byte, error){
	return json.Marshal([]string{event.title, event.content, event.alertTime.Format(time.RFC1123)})
}