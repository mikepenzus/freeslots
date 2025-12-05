package utils

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

type CalendarEvent struct {
	StartTime   time.Time
	Duration    int
	Description string
	Timezone    string
}

func (calendarEvent CalendarEvent) GetEndTime() time.Time {
	timeDuration := time.Duration(calendarEvent.Duration) * time.Minute
	endDate := calendarEvent.StartTime.Add(timeDuration)
	return endDate
}

func (event CalendarEvent) Print() {
	fmt.Printf("%v - %v - %v - %v\n", event.StartTime,
		event.Duration, event.Description, event.Timezone)
}

type DailyAgenda struct {
	Date   time.Time
	Events []CalendarEvent
}

func (dailyAgenda DailyAgenda) Print(showDescription, showSlotDuration bool) {
	fmt.Printf("%s: ", dailyAgenda.Date.Format("2 Jan 2006"))
	for index, event := range dailyAgenda.Events {
		if index > 0 {
			fmt.Print(", ")
		}
		if showSlotDuration {
			fmt.Printf("%s-%s (%v')", event.StartTime.Format("15:04"), event.GetEndTime().Format("15:04 MST"),
				event.Duration)

		} else {
			fmt.Printf("%s-%s", event.StartTime.Format("15:04"), event.GetEndTime().Format("15:04 MST"))
		}
		if showDescription {
			fmt.Printf(" (%s)", event.Description)
		}
	}
	fmt.Println()
}

func (dailyAgenda DailyAgenda) PrintHtml(showDescription, showSlotDuration bool) {
	for _, event := range dailyAgenda.Events {
		if showSlotDuration {
			fmt.Printf("<tr><td>%s</td><td>%s-%s (%v')</td>", dailyAgenda.Date.Format("2 Jan 2006"),
				event.StartTime.Format("15:04"), event.GetEndTime().Format("15:04 MST"), event.Duration)
		} else {
			fmt.Printf("<tr><td>%s</td><td>%s-%s</td>", dailyAgenda.Date.Format("2 Jan 2006"),
				event.StartTime.Format("15:04"), event.GetEndTime().Format("15:04 MST"))
		}
		if showDescription {
			fmt.Printf("<td>%s</td>", event.Description)
		}
		fmt.Println("</tr>")
	}
}

func (dailyAgenda DailyAgenda) PrintMarkdown(showDescription, showSlotDuration bool) {
	for _, event := range dailyAgenda.Events {
		if showSlotDuration {
			fmt.Printf("| %s | %s-%s (%v') |", dailyAgenda.Date.Format("2 Jan 2006"),
				event.StartTime.Format("15:04"), event.GetEndTime().Format("15:04 MST"), event.Duration)
		} else {
			fmt.Printf("| %s | %s-%s |", dailyAgenda.Date.Format("2 Jan 2006"),
				event.StartTime.Format("15:04"), event.GetEndTime().Format("15:04 MST"))
		}
		if showDescription {
			fmt.Printf("| %s |", event.Description)
		}
		fmt.Println()
	}
}

func (dailyAgenda DailyAgenda) IsEmpty() bool {
	return len(dailyAgenda.Events) == 0
}

func (dailyAgenda DailyAgenda) IsWeekend() bool {
	return dailyAgenda.Date.Weekday() == time.Sunday || dailyAgenda.Date.Weekday() == time.Saturday
}

func GetTimeWithSpecificHoursMinutes(origTime time.Time, hours, minutes int) time.Time {
	return time.Date(origTime.Year(), origTime.Month(), origTime.Day(), hours, minutes, 0, 0, origTime.Location())
}

// print event list
func PrintEventList(eventList []CalendarEvent) {
	for _, event := range eventList {
		event.Print()
	}
}

func GetPureDate(origTime time.Time) time.Time {
	return GetTimeWithSpecificHoursMinutes(origTime, 0, 0)
}

func CompareTimesWithoutTimeZone(firstTime, secondTime time.Time) int {
	newFirstTime := time.Date(firstTime.Year(), firstTime.Month(), firstTime.Day(),
		firstTime.Hour(), firstTime.Minute(), firstTime.Second(), 0, secondTime.Location())
	return newFirstTime.Compare(secondTime)
}

func CreateDefaultCalendarEvent(currentDay time.Time, startTimeHour int, startTimeMin int, duration int, description string) CalendarEvent {
	now := currentDay
	calendarEvent := CalendarEvent{
		Duration:    duration,
		Description: description,
		Timezone:    now.Local().String(),
	}
	StartTime := time.Date(now.Year(), now.Month(), now.Day(), startTimeHour, startTimeMin,
		0, 0, now.Local().Location())
	calendarEvent.StartTime = StartTime
	return calendarEvent
}

func CreateDefaultCalendarEventFromString(currentDay time.Time, startTime string, duration int, description string) CalendarEvent {
	now := currentDay
	calendarEvent := CalendarEvent{
		Duration:    duration,
		Description: description,
		Timezone:    now.Local().String(),
	}
	startTimeParts := strings.Split(startTime, ":")
	hours, _ := strconv.Atoi(startTimeParts[0])
	mins, _ := strconv.Atoi(startTimeParts[1])
	StartTime := time.Date(now.Year(), now.Month(), now.Day(), hours, mins,
		0, 0, now.Local().Location())
	calendarEvent.StartTime = StartTime
	return calendarEvent
}

// parse time in the form "HH:MM"
func ParseTime(timeAsString string) (int, int) {
	parts := strings.Split(timeAsString, ":")
	hours, _ := strconv.Atoi(parts[0])
	mins, _ := strconv.Atoi(parts[1])
	return hours, mins
}

func SortEventListByStartTime(eventList *[]CalendarEvent) {
	slices.SortFunc(*eventList, func(a, b CalendarEvent) int {
		return a.StartTime.Compare(b.StartTime)
	})
}

// Get events from Google Calendar
func GetEventsFromPrimaryCalendar(srv *calendar.Service, tMin time.Time, noDays int, userMail string) ([]DailyAgenda, error) {
	// TODO: add option to manage multiple calendars
	tMinAsString := tMin.Format(time.RFC3339)
	tMaxAsString := tMin.AddDate(0, 0, noDays).Format(time.RFC3339)
	events, err := srv.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(tMinAsString).
		TimeMax(tMaxAsString).
		MaxResults(2500).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, err
	}

	eventList := []CalendarEvent{}
	for _, item := range events.Items {
		// scanning attendees to see if I declined the event
		// TODO: fix
		canAttend := true
		for _, attendee := range item.Attendees {
			if attendee.Email == userMail && attendee.ResponseStatus == "declined" {
				canAttend = false
				break
			}
		}
		if !canAttend {
			continue
		}

		newEvent := CalendarEvent{}
		newEvent.Description = item.Summary
		newEvent.Timezone = item.Start.TimeZone
		from := item.Start.DateTime
		if from == "" {
			newEvent.StartTime, _ = time.Parse(time.DateOnly, item.Start.Date)
		} else {
			newEvent.StartTime, _ = time.Parse(time.RFC3339, from)
		}
		end := item.End.DateTime
		if end == "" {
			endDate, _ := time.Parse(time.DateOnly, item.End.Date)
			newEvent.Duration = int(endDate.Sub(newEvent.StartTime).Minutes())
		} else {
			endDate, _ := time.Parse(time.RFC3339, end)
			newEvent.Duration = int(endDate.Sub(newEvent.StartTime).Minutes())
		}
		// skip all day-events
		if newEvent.Duration != 1440 {
			eventList = append(eventList, newEvent)
		}
	}

	SortEventListByStartTime(&eventList)
	var dailyAgendas []DailyAgenda = SplitCalendarEventsByDay(eventList)
	return dailyAgendas, nil
}

func ParseDailyAgenda(singleDayAgenda string) (DailyAgenda, error) {
	events, err := ParseSingleDayAgenda(singleDayAgenda)
	if err != nil {
		return DailyAgenda{}, err
	}
	singleDayAgendaParts := strings.Split(singleDayAgenda, ",")
	currentDay := GetPureDate(time.Now())
	for _, part := range singleDayAgendaParts {
		switch part[0] {
		case 'd':
			dateStr := part[1:]
			parsedTime, err := time.Parse(time.DateOnly, dateStr)
			if err != nil {
				return DailyAgenda{}, err
			}
			currentDay = parsedTime
		}
	}
	dailyAgenda := DailyAgenda{
		Events: events,
		Date:   currentDay,
	}
	return dailyAgenda, nil
}

func ParseSingleDayAgenda(singleDayAgenda string) ([]CalendarEvent, error) {
	// Format:
	// d<yyyy-mm-dd>,m<number>,s<number>,a<agenda>
	// d<yyyy-mm-dd> = date of the single day agenda (pay attention to the format!)
	// m<number> = duration in minutes of each slot
	// s<number> = skip the first <number> slots from the midnight
	// a<agenda> = list of characters representing a day agenda
	// The agenda starts at midnight and finishes at 23.59.
	// A character of the agenda could mean: dash=available empty slot,
	// alphabetic character=event lasting half an hour,
	// consecutive alphabetic characters are merged into longer events if they share the same character.
	// E.g.: "m30,s16,aXX--XY--XYYX--ZZZZZZ" means
	// (8:00, event X, 60 minutes), (10:00, event X, 30 minutes), (10:30, event Y, 30 minutes),
	// (12:00, event X, 30 minutes), (12:30, event Y, 60 minutes), (14:00, event Z, 180 minutes)
	// Multiple agenda items will be merged into a single list, e.g.:
	// "m30,s16,aXX--XY--XYYX--ZZZZZZ,aAAAAAA-BB"
	singleDayAgendaParts := strings.Split(singleDayAgenda, ",")
	slotDuration := 60
	slotsToSkip := 0
	dailyAgendasToScan := make([]string, 0)
	currentDay := GetPureDate(time.Now())
	for _, part := range singleDayAgendaParts {
		switch part[0] {
		case 'd':
			dateStr := part[1:]
			parsedTime, err := time.Parse(time.DateOnly, dateStr)
			if err != nil {
				return nil, err
			}
			currentDay = parsedTime
			// TODO: set time zone
		case 'm':
			minStr := part[1:]
			slotDuration, _ = strconv.Atoi(minStr)
		case 's':
			skipStr := part[1:]
			slotsToSkip, _ = strconv.Atoi(skipStr)
		case 'a':
			currentAgenda := part[1:]
			dailyAgendasToScan = append(dailyAgendasToScan, currentAgenda)
		}
	}
	totalCalendarEvents := make([]CalendarEvent, 0)

	for _, dailyAgendaToScan := range dailyAgendasToScan {
		calendarEvents := make([]CalendarEvent, 0)
		previousChar := '-'
		currentEventDescription := ""
		currentEventDuration := 0
		currentEventStartHour := 0
		currentEventStartMin := 0
		for currentEventSlotIndex, currentChar := range dailyAgendaToScan {
			if currentChar == '-' {
				if previousChar == '-' {
					continue
				}
				// finalize calendar event and store it
				newCalendarEvent := CreateDefaultCalendarEvent(currentDay, currentEventStartHour, currentEventStartMin,
					currentEventDuration, currentEventDescription)
				calendarEvents = append(calendarEvents, newCalendarEvent)
				previousChar = currentChar
				currentEventDescription = ""
				currentEventDuration = 0
				currentEventStartHour = 0
				currentEventStartMin = 0
			} else {
				if previousChar == '-' {
					// put a placeholder for a new calendar event to be finalized later
					currentEventDescription = string(currentChar)
					currentEventDuration = slotDuration
					tempEventStartMin := (currentEventSlotIndex + slotsToSkip) * slotDuration
					currentEventStartHour = tempEventStartMin / 60
					currentEventStartMin = tempEventStartMin % 60
					previousChar = currentChar
				} else {
					// check if the current event is different from the previous one, finalize the current event if so
					if previousChar == currentChar {
						currentEventDuration += slotDuration
					} else {
						newCalendarEvent := CreateDefaultCalendarEvent(currentDay, currentEventStartHour, currentEventStartMin,
							currentEventDuration, currentEventDescription)
						calendarEvents = append(calendarEvents, newCalendarEvent)
						currentEventDescription = string(currentChar)
						currentEventDuration = slotDuration
						tempEventStartMin := (currentEventSlotIndex + slotsToSkip) * slotDuration
						currentEventStartHour = tempEventStartMin / 60
						currentEventStartMin = tempEventStartMin % 60
						previousChar = currentChar
					}
				}
			}
		}
		if previousChar != '-' {
			// finalize calendar event and store it
			newCalendarEvent := CreateDefaultCalendarEvent(currentDay, currentEventStartHour, currentEventStartMin,
				currentEventDuration, currentEventDescription)
			calendarEvents = append(calendarEvents, newCalendarEvent)
		}

		totalCalendarEvents = MergeCalendarEventLists(totalCalendarEvents, calendarEvents)
	}

	return totalCalendarEvents, nil
}
