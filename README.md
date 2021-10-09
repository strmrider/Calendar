Calenar managment and tasks scheduler

## How to use
Import package
```Go
package main

import (
  "fmt"
  calendar "github.com/strmrider/Calendar"
)

func main(){
}
```
### Calendar
Calendar years range is 1900 to 2099. Each year is initialized automatically.
```Go
 // creates a new calendar within the range of 1965 to 2012
 cal := calendar.CreateCalendar(1965, 2012)
 
 // returns whether a year exists in the calendar
 fmt.Println(cal.IsYearExist(1982))
 // returns calendar's range
 fmt.Println(cal.GetRange())
 // add years to the calendar
 cal.ExpandYears(4)
 fmt.Println(cal.GetRange())
 // returns Day type (see Day section)
 cal.GetDay(1967, calendar.June, 4)
 // returns Year type (see Year section)
 cal.GetYear(1978)
```
Output
```
true
1965 2012
1965 2017
```
#### Year
```Go
year := GetYear(1978)

// returns Day type (see Day section)
year.GetDay(calendar.April, 14)
// returns Month type (see Month section)
year.GetMonth(calendar.December)

// returns whether the year is a leap year
fmt.Println(year.IsLeap())
// prints all the mounths
year.Print()
```
#### Month
```Go
month := year.GetMonth(calendar.December)

// returns Day type (see Day section)
month.GetDay(12)
// retuns month's name
fmt.Println(month.GetMonth())
// returns number of days
fmt.Println(month.NumOfDays())
// returns number of events in month
fmt.Println(month.NumOfEvents())
// prints month summary
month.Summary()
// print month (a day with an event will be marked with *)
month.Print()
```
```Output
December
31
0
December
1 Friday: 0 events
2 Saturday: 0 events
3 Sunday: 0 events
4 Monday: 0 events
5 Tuesday: 0 events
6 Wednesday: 0 events
7 Thursday: 0 events
8 Friday: 0 events
9 Saturday: 0 events
10 Sunday: 0 events
11 Monday: 0 events
12 Tuesday: 0 events
13 Wednesday: 0 events
14 Thursday: 0 events
15 Friday: 0 events
16 Saturday: 0 events
17 Sunday: 0 events
18 Monday: 0 events
19 Tuesday: 0 events
20 Wednesday: 0 events
21 Thursday: 0 events
22 Friday: 0 events
23 Saturday: 0 events
24 Sunday: 0 events
25 Monday: 0 events
26 Tuesday: 0 events
27 Wednesday: 0 events
28 Thursday: 0 events
29 Friday: 0 events
30 Saturday: 0 events
31 Sunday: 0 events

December
Sun     Mon     Tue     Wed     Thu     Fri     Sat
====================================================
                                        1       2
====================================================
3       4       5       6       7       8       9
====================================================
10      11      12      13      14      15      16
====================================================
17      18      19      20      21      22      23
====================================================
24      25      26      27      28      29      30
====================================================
31
```
#### Day
```Go
day := year.GetDay(calendar.April, 14)

// adds a new event. receives Event type (see Events section)
day.AddEvent(NEW_EVENT)
// retuns an event by given index in day's events list
day.GetEvent(EVENT_INDEX)
// retuns the number of events in a day
day.NumOfEvents()
// removes an event by given index in day's events list
day.RemoveEvent(EVENT_INDEX)
// clears day's events list
day.ClearEvents()
```
### Events and tasks
Events include several properties:
* Title
* Content
* Date
* Reminder function; runs when event's time is up
```Go
// alert function, exceuted when time is up
func alert(Event* event){
  event.Print()
}

title := "sample event"
content = "Event is done"
// Y/M/D H/M/S format
alertDate = calendar.BuildTime(2011, calendar.January, 8, 23, 30, 0)
event := calendar.CreatEvent(title, content, alertDate, alert)
// starts timer
event.RunTimer()
// stops timer
event.StopTimer()
// prints event's details
event.Print()
```
### Serialization
Serialize the calendar into Json format in bytes
```Go
// serialize
bin := cal.Serialize()

// deserialize with a reminder function for existed events (using alert funcion from previous sections)
cal := calendar.Deserialize(bin, alert)
```
### Save and load from disk
```Go
path = "mycalendar.cl"
// saves the calendar in a file
cal.SaveToDisk(path)

// load calendar from file with a reminder function for existed events (using alert funcion from previous sections)
cal := calendar.LoadFromFile(path, alert)
```
### Network
Share and consume calendars via network
```Go
var port int = 4231
// the calendar listens for requests and send its data to the clients
cl.ShareOnNetwork(port)

calendar_host = "141.25.41.4"
// read calendar from via network with a reminder function for existed events (using alert funcion from previous sections)
newCal := calendar.ReadCalendarFromNetwork(calendar_host, port, alert)
```
