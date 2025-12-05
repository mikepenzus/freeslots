package utils

import (
	"time"
)

// return agenda with time constraints
// maxHours can be 24 in order to make sure that we can scan the whole day
func (dailyAgenda DailyAgenda) Constrain(minHours, minMinutes, maxHours, maxMinutes int) DailyAgenda {
	resultDailyAgenda := DailyAgenda{
		Date:   dailyAgenda.Date,
		Events: []CalendarEvent{},
	}
	minTime := GetTimeWithSpecificHoursMinutes(dailyAgenda.Date, minHours, minMinutes)
	var maxTime time.Time
	if maxHours == 24 {
		maxTime = GetTimeWithSpecificHoursMinutes(dailyAgenda.Date.AddDate(0, 0, 1), 0, 0)
	} else {
		maxTime = GetTimeWithSpecificHoursMinutes(dailyAgenda.Date, maxHours, maxMinutes)
	}
	for _, event := range dailyAgenda.Events {
		currentEvent := CalendarEvent{
			StartTime:   event.StartTime,
			Duration:    event.Duration,
			Description: event.Description,
			Timezone:    event.Timezone,
		}
		if currentEvent.StartTime.Compare(minTime) < 0 {
			// trim current event
			durationToTrim := int(minTime.Sub(currentEvent.StartTime).Minutes())
			currentEvent.StartTime = minTime
			currentEvent.Duration -= durationToTrim
		}
		endTime := currentEvent.GetEndTime()
		if endTime.Compare(maxTime) > 0 {
			// trim current event
			durationToTrim := int(endTime.Sub(maxTime).Minutes())
			currentEvent.Duration -= durationToTrim
		}
		if currentEvent.Duration > 0 {
			resultDailyAgenda.Events = append(resultDailyAgenda.Events, currentEvent)
		}
	}
	return resultDailyAgenda
}

func (dailyAgenda DailyAgenda) GetFreeSlots(minDuration, fromHours, fromMinutes,
	toHours, toMinutes int) (DailyAgenda, error) {
	agendaWithOnlyFreeSlots := DailyAgenda{
		Date:   dailyAgenda.Date,
		Events: []CalendarEvent{},
	}
	newDailyAgenda := dailyAgenda.Constrain(fromHours, fromMinutes, toHours, toMinutes)
	/*
		fmt.Println("*** Agenda after constraining...") // TODO: debug...
		newDailyAgenda.Print(true)
	*/
	gluedEvents := GlueCalendarEvents(newDailyAgenda.Events)

	currentSlotStartTime := time.Date(newDailyAgenda.Date.Year(), newDailyAgenda.Date.Month(), newDailyAgenda.Date.Day(),
		fromHours, fromMinutes, 0, 0, newDailyAgenda.Date.Location())

	for _, busyEvent := range gluedEvents {
		switch busyEvent.StartTime.Compare(currentSlotStartTime) {
		case 1:
			// new free slot to append to the list of events to return
			currentSlotDuration := int(busyEvent.StartTime.Sub(currentSlotStartTime).Minutes())
			newEvent := CalendarEvent{
				StartTime:   currentSlotStartTime,
				Duration:    currentSlotDuration,
				Description: "*",
				Timezone:    currentSlotStartTime.Location().String(), // TODO: fix timezone
			}
			if newEvent.Duration >= minDuration {
				agendaWithOnlyFreeSlots.Events = append(agendaWithOnlyFreeSlots.Events, newEvent)
			}
			currentSlotStartTime = busyEvent.GetEndTime()
		case 0:
			// no new free slots, just pass to the next event in the list
			currentSlotStartTime = busyEvent.GetEndTime()
		case -1:
			// TODO: double check this point, it should never occur
			currentSlotStartTime = busyEvent.GetEndTime()
		}
	}
	endOfTodayAllowedTimeRange := time.Date(newDailyAgenda.Date.Year(), newDailyAgenda.Date.Month(), newDailyAgenda.Date.Day(),
		toHours, toMinutes, 0, 0, newDailyAgenda.Date.Location())
	if endOfTodayAllowedTimeRange.Compare(currentSlotStartTime) > 0 {
		// new free slot to append to the list of events to return
		currentSlotDuration := int(endOfTodayAllowedTimeRange.Sub(currentSlotStartTime).Minutes())
		newEvent := CalendarEvent{
			StartTime:   currentSlotStartTime,
			Duration:    currentSlotDuration,
			Description: "*",
			Timezone:    currentSlotStartTime.Location().String(), // TODO: fix timezone
		}
		if newEvent.Duration >= minDuration {
			agendaWithOnlyFreeSlots.Events = append(agendaWithOnlyFreeSlots.Events, newEvent)
		}
	}
	return agendaWithOnlyFreeSlots, nil
}

// splits calendar events into days
// assumption: they are sorted by StartTime
func SplitCalendarEventsByDay(inputEvents []CalendarEvent) []DailyAgenda {
	outputDailyAgendas := []DailyAgenda{}
	var currentDailyAgendaDate time.Time
	var currentDailyAgendaEvents []CalendarEvent
	initMode := true
	for _, currentEvent := range inputEvents {
		if initMode {
			currentDailyAgendaDate = GetPureDate(currentEvent.StartTime)
			currentDailyAgendaEvents = []CalendarEvent{currentEvent}
			initMode = false
			continue
		}
		dayOfCalendarEvent := GetPureDate(currentEvent.StartTime)
		if CompareTimesWithoutTimeZone(currentDailyAgendaDate, dayOfCalendarEvent) == 0 {
			// events are on the same day
			currentDailyAgendaEvents = append(currentDailyAgendaEvents, currentEvent)
		} else {
			// events are on different days
			outputDailyAgendas = append(outputDailyAgendas, DailyAgenda{
				Date:   currentDailyAgendaDate,
				Events: currentDailyAgendaEvents,
			})
			currentDailyAgendaDate = GetPureDate(currentEvent.StartTime)
			currentDailyAgendaEvents = []CalendarEvent{currentEvent}
		}
	}
	outputDailyAgendas = append(outputDailyAgendas, DailyAgenda{
		Date:   currentDailyAgendaDate,
		Events: currentDailyAgendaEvents,
	})

	return outputDailyAgendas
}

// merge events by creating a single list of events sorted by start date
func MergeCalendarEventLists(firstList, secondList []CalendarEvent) []CalendarEvent {
	mergedList := make([]CalendarEvent, 0, len(firstList)+len(secondList))
	mergedList = append(mergedList, firstList...)
	mergedList = append(mergedList, secondList...)
	SortEventListByStartTime(&mergedList)
	return mergedList
}

// glue events sorted by start date by creating a single list with the minimum set of slots,
// regardless of their descriptions
// Assumption: all events belong to the same day
func GlueCalendarEvents(eventList []CalendarEvent) []CalendarEvent {
	mergedList := make([]CalendarEvent, 0)
	if len(eventList) == 0 {
		return mergedList
	}
	if len(eventList) == 1 {
		mergedList = append(mergedList, eventList[0])
		return mergedList
	}
	var candidateEventToGlue CalendarEvent
	for index, currentEvent := range eventList {
		if index == 0 {
			candidateEventToGlue = currentEvent
			continue
		}
		endTimeOfCandidateEvent := candidateEventToGlue.GetEndTime()
		endTimeOfCurrentEvent := currentEvent.GetEndTime()
		switch endTimeOfCandidateEvent.Compare(currentEvent.StartTime) {
		case -1:
			// two events that can't be glued
			mergedList = append(mergedList, candidateEventToGlue)
			candidateEventToGlue = currentEvent
		case 0:
			// the 2 events are exactly one after the other - let's glue them
			candidateEventToGlue.Duration += currentEvent.Duration
		case 1:
			// the first event ends after the start of the 2nd event - let's glue them
			// if first event ends after the 2nd, don't change the first event
			if CompareTimesWithoutTimeZone(endTimeOfCandidateEvent, endTimeOfCurrentEvent) == -1 {
				gapBetweenEndOfTwoEvents := endTimeOfCurrentEvent.Sub(endTimeOfCandidateEvent)
				overlap := int(gapBetweenEndOfTwoEvents.Minutes())
				candidateEventToGlue.Duration += currentEvent.Duration - overlap
			}
		}
	}
	// append pending event
	mergedList = append(mergedList, candidateEventToGlue)
	return mergedList
}

func FillInWithEmptyDays(dailyAgendas []DailyAgenda, tMin time.Time, noDays int, skipWeekends bool) ([]DailyAgenda, error) {
	// map agendas to dates
	now := time.Now()
	newDailyAgendas := []DailyAgenda{}
	mapOfDailyAgendas := make(map[time.Time]DailyAgenda)
	for _, dailyAgenda := range dailyAgendas {
		// TODO: I neutralized the time zone, fix this issue
		agendaTime := dailyAgenda.Date
		timeWithNeutralLocation := time.Date(agendaTime.Year(), agendaTime.Month(), agendaTime.Day(),
			0, 0, 0, 0, now.Location())
		mapOfDailyAgendas[timeWithNeutralLocation] = dailyAgenda
	}
	for dayIndex := 0; dayIndex < noDays; dayIndex++ {
		// TODO: I neutralized the time zone, fix this issue
		agendaTime := tMin.AddDate(0, 0, dayIndex)
		currentDate := time.Date(agendaTime.Year(), agendaTime.Month(), agendaTime.Day(),
			0, 0, 0, 0, now.Location())
		if skipWeekends && (currentDate.Weekday() == time.Sunday || currentDate.Weekday() == time.Saturday) {
			continue
		}
		targetAgenda, agendaIndex := mapOfDailyAgendas[currentDate]
		if agendaIndex {
			// agenda found
			newDailyAgendas = append(newDailyAgendas, targetAgenda)
		} else {
			// need to create a new agenda
			newTargetAgenda := DailyAgenda{
				Date:   currentDate,
				Events: []CalendarEvent{},
			}
			newDailyAgendas = append(newDailyAgendas, newTargetAgenda)
		}
	}
	return newDailyAgendas, nil
}
