package utils

import (
	"testing"
	"time"
)

func TestGlueCalendarEvents(t *testing.T) {
	agendas := []string{
		"m30,s16,a--,a-",
		"m30,s16,a-X,a--",
		"m30,s16,aXX,a--YY",
		"m30,s16,aXX,a-YY",
		"m30,s16,aXX,a---YY",
		"m60,s10,aXX,a----YY",
		"m30,s16,aXXXX,aYY",
		"m30,s16,aXXXX,a-YY",
		"m30,s16,aXXXX,a--YY",
		"m30,s16,aXXXX,a---YZ",
		"m60,s4,aX-Y-Z-A-B-C-D---,a-X-Y-Z-A-B-C-D---",
	}
	expectedAgendas := []string{
		"m30,s16,a",
		"m30,s16,a-X",
		"m30,s16,aXXXX",
		"m30,s16,aXXX",
		"m30,s16,aXX-YY",
		"m60,s10,aXX--YY",
		"m30,s16,aXXXX",
		"m30,s16,aXXXX",
		"m30,s16,aXXXX",
		"m30,s16,aXXXXX",
		"m60,s4,aXXXXXXXXXXXXXX",
	}
	for agendaIndex, agenda := range agendas {
		inputEvents, _ := ParseSingleDayAgenda(agenda)
		expectedAgenda := expectedAgendas[agendaIndex]
		expectedOutputEvents, _ := ParseSingleDayAgenda(expectedAgenda)
		outputEvents := GlueCalendarEvents(inputEvents)
		if len(outputEvents) != len(expectedOutputEvents) {
			t.Errorf("Error while gluing: agenda index %v, agenda %v, expected agenda %v, no. expected %v, no. results %v. Different lengths",
				agendaIndex, agenda, expectedAgenda, len(outputEvents), len(expectedOutputEvents))
			PrintEventList(outputEvents)
			return
		}
		for eventIndex, expectedEvent := range expectedOutputEvents {
			outputEvent := outputEvents[eventIndex]
			if outputEvent.Duration != expectedEvent.Duration ||
				outputEvent.Description != expectedEvent.Description ||
				outputEvent.StartTime.Hour() != expectedEvent.StartTime.Hour() ||
				outputEvent.StartTime.Minute() != expectedEvent.StartTime.Minute() {
				t.Errorf("Error while parsing: agenda index %v, agenda %v, mismatching event index %v", agendaIndex, agenda, eventIndex)
				return
			}
		}
	}
}

func TestSplitCalendarEventsByDay(t *testing.T) {
	timeNow := time.Now()
	timeTomorrow := timeNow.AddDate(0, 0, 1)
	inputEvents := []CalendarEvent{
		CreateDefaultCalendarEventFromString(timeNow, "08:00", 60, "X"),
		CreateDefaultCalendarEventFromString(timeNow, "10:00", 30, "X"),
		CreateDefaultCalendarEventFromString(timeNow, "10:30", 30, "Y"),
		CreateDefaultCalendarEventFromString(timeNow, "12:00", 30, "X"),
	}
	dailyAgendas := SplitCalendarEventsByDay(inputEvents)
	if len(dailyAgendas) != 1 || len(dailyAgendas[0].Events) != len(inputEvents) {
		t.Errorf("Error parsing set of of events # 1")
		for _, inputEvent := range inputEvents {
			inputEvent.Print()
		}
		for _, dailyAgenda := range dailyAgendas {
			dailyAgenda.Print(true, true)
		}
		return
	}

	inputEvents = []CalendarEvent{
		CreateDefaultCalendarEventFromString(timeNow, "08:00", 60, "X"),
		CreateDefaultCalendarEventFromString(timeNow, "10:00", 30, "X"),
		CreateDefaultCalendarEventFromString(timeTomorrow, "10:30", 30, "Y"),
		CreateDefaultCalendarEventFromString(timeTomorrow, "12:00", 30, "X"),
		CreateDefaultCalendarEventFromString(timeTomorrow, "14:00", 30, "X"),
	}
	dailyAgendas = SplitCalendarEventsByDay(inputEvents)
	if len(dailyAgendas) != 2 || len(dailyAgendas[0].Events) != 2 || len(dailyAgendas[1].Events) != 3 {
		t.Errorf("Error parsing set of of events # 2")
		for _, inputEvent := range inputEvents {
			inputEvent.Print()
		}
		for _, dailyAgenda := range dailyAgendas {
			dailyAgenda.Print(true, true)
		}
		return
	}

	timeWithTZ1, _ := time.Parse(time.RFC3339, "2025-12-10T00:00:00+01:00")
	timeWithTZ2, _ := time.Parse(time.RFC3339, "2025-12-10T00:00:00+00:00")
	inputEvents = []CalendarEvent{
		CreateDefaultCalendarEventFromString(timeWithTZ1, "08:00", 60, "X"),
		CreateDefaultCalendarEventFromString(timeWithTZ1, "10:00", 30, "X"),
		CreateDefaultCalendarEventFromString(timeWithTZ2, "10:30", 30, "Y"),
		CreateDefaultCalendarEventFromString(timeWithTZ2, "12:00", 30, "X"),
		CreateDefaultCalendarEventFromString(timeWithTZ1, "14:00", 30, "X"),
	}
	dailyAgendas = SplitCalendarEventsByDay(inputEvents)
	if len(dailyAgendas) != 1 || len(dailyAgendas[0].Events) != 5 {
		t.Errorf("Error parsing set of of events # 3")
		for _, inputEvent := range inputEvents {
			inputEvent.Print()
		}
		for _, dailyAgenda := range dailyAgendas {
			dailyAgenda.Print(true, true)
		}
		return
	}

	timeWithTZ1, _ = time.Parse(time.RFC3339, "2025-12-10T00:00:00+01:00")
	timeWithTZ2, _ = time.Parse(time.DateOnly, "2025-12-10")
	inputEvents = []CalendarEvent{
		CreateDefaultCalendarEventFromString(timeWithTZ1, "08:00", 60, "X"),
		CreateDefaultCalendarEventFromString(timeWithTZ1, "10:00", 30, "X"),
		CreateDefaultCalendarEventFromString(timeWithTZ2, "10:30", 30, "Y"),
		CreateDefaultCalendarEventFromString(timeWithTZ2, "12:00", 30, "X"),
		CreateDefaultCalendarEventFromString(timeWithTZ1, "14:00", 30, "X"),
	}
	dailyAgendas = SplitCalendarEventsByDay(inputEvents)
	if len(dailyAgendas) != 1 || len(dailyAgendas[0].Events) != 5 {
		t.Errorf("Error parsing set of of events # 4")
		for _, inputEvent := range inputEvents {
			inputEvent.Print()
		}
		for _, dailyAgenda := range dailyAgendas {
			dailyAgenda.Print(true, true)
		}
		return
	}

	timeWithTZ1, _ = time.Parse(time.RFC3339, "2025-12-10T00:00:00+01:00")
	timeWithTZ2, _ = time.Parse(time.DateOnly, "2025-12-10")
	inputEvents = []CalendarEvent{
		{
			StartTime:   timeWithTZ1,
			Duration:    60,
			Description: "Test",
		},
		{
			StartTime:   timeWithTZ1,
			Duration:    60,
			Description: "Test",
		},
		{
			StartTime:   timeWithTZ2,
			Duration:    60,
			Description: "Test",
		},
		{
			StartTime:   timeWithTZ2,
			Duration:    60,
			Description: "Test",
		},
		{
			StartTime:   timeWithTZ1,
			Duration:    60,
			Description: "Test",
		},
	}
	dailyAgendas = SplitCalendarEventsByDay(inputEvents)
	if len(dailyAgendas) != 1 || len(dailyAgendas[0].Events) != 5 {
		t.Errorf("Error parsing set of of events # 5")
		for _, inputEvent := range inputEvents {
			inputEvent.Print()
		}
		for _, dailyAgenda := range dailyAgendas {
			dailyAgenda.Print(true, true)
		}
		return
	}

}

func TestConstrain(t *testing.T) {
	agendas := []string{
		"m30,s0,aXX--XY--XYYX--ZZZZZZ",
		"m60,s4,aX-Y-Z-A-B-C-D---",
		"m60,s8,aXXXXYYYYZZZZAAAA",
		"m120,aXXXXXXXYYYYY",
		"m60,s0,a--------------------",
		"m30,aAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
	}
	expectedAgendasAfterConstraining := []string{
		"m30,s16,aZZZZ",
		"m60,s4,a----Z-A-B-C-D---",
		"m60,s8,aXXXXYYYYZZ",
		"m120,a----XXXYY",
		"m60,s0,a--------------------",
		"m30,s16,aAAAAAAAAAAAAAAAAAAAA",
	}
	for agendaIndex, agenda := range agendas {
		inputAgenda, _ := ParseDailyAgenda(agenda)
		expectedAgendaAfterConstraining := expectedAgendasAfterConstraining[agendaIndex]
		expectedAgenda, _ := ParseDailyAgenda(expectedAgendaAfterConstraining)
		outputAgenda := inputAgenda.Constrain(8, 0, 18, 0)
		if len(expectedAgenda.Events) != len(outputAgenda.Events) {
			t.Errorf("Error while parsing: agenda index %v, agenda %v, no. expected %v, no. results %v. Different lengths",
				agendaIndex, agenda, len(expectedAgenda.Events), len(outputAgenda.Events))
			PrintEventList(outputAgenda.Events)
			return
		}
		for eventIndex, expectedEvent := range expectedAgenda.Events {
			resultingEvent := outputAgenda.Events[eventIndex]
			if resultingEvent.Duration != expectedEvent.Duration ||
				resultingEvent.Description != expectedEvent.Description ||
				resultingEvent.StartTime.Hour() != expectedEvent.StartTime.Hour() ||
				resultingEvent.StartTime.Minute() != expectedEvent.StartTime.Minute() {
				t.Errorf("Error while parsing: agenda index %v, agenda %v, mismatching event index %v", agendaIndex, agenda, eventIndex)
				return
			}
		}
	}
}

func TestGetFreeSlots(t *testing.T) {
	agendas := []string{
		"m30,s16,aXX--XY--XYYX--ZZZZZZ",
		"m30,s14,aXXXX--XY--XYYX--ZZZZZZ",
		"m30,s16,aXX--XY--XYYX--ZZZZZZ----ZZ",
		"m60,s0",
		"m60,s7,aXX",
		"m60,s8,a---------XX",
	}
	expectedAgendas := []string{
		"m30,s16,a--**--**----**",
		"m30,s16,a--**--**----**",
		"m30,s16,a--**--**----**",
		"m60,s8,a**********",
		"m60,s7,a--*********",
		"m60,s8,a*********",
	}
	for agendaIndex, agenda := range agendas {
		inputAgenda, _ := ParseDailyAgenda(agenda)
		expectedAgendaAsString := expectedAgendas[agendaIndex]
		expectedAgenda, _ := ParseDailyAgenda(expectedAgendaAsString)
		outputAgenda, err := inputAgenda.GetFreeSlots(30, 8, 0, 18, 0)
		if err != nil {
			t.Errorf("Error while parsing: agenda index %v, agenda %v, expected agenda %v, error %v",
				agendaIndex, agenda, expectedAgenda, err)
			return
		}
		if len(outputAgenda.Events) != len(expectedAgenda.Events) {
			t.Errorf("Error while parsing: agenda index %v, agenda %v, expected agenda %v, no. expected %v, no. results %v. Different lengths",
				agendaIndex, agenda, expectedAgenda, len(outputAgenda.Events), len(expectedAgenda.Events))
			return
		}
		for eventIndex, expectedEvent := range expectedAgenda.Events {
			outputEvent := outputAgenda.Events[eventIndex]
			if outputEvent.Duration != expectedEvent.Duration ||
				outputEvent.Description != expectedEvent.Description ||
				outputEvent.StartTime.Hour() != expectedEvent.StartTime.Hour() ||
				outputEvent.StartTime.Minute() != expectedEvent.StartTime.Minute() {
				t.Errorf("Error while parsing: agenda index %v, agenda %v, mismatching event index %v", agendaIndex,
					agenda, eventIndex)
				return
			}
		}
	}
}

func TestFillInWithEmptyDays(t *testing.T) {
	today := GetPureDate(time.Now())
	tomorrow := today.AddDate(0, 0, 1)
	calendarEvents0 := []CalendarEvent{
		CreateDefaultCalendarEvent(today, 0, 0, 0, "Home"),
	}
	calendarEvents1 := []CalendarEvent{
		CreateDefaultCalendarEvent(tomorrow, 0, 0, 0, "Home"),
		CreateDefaultCalendarEvent(tomorrow, 8, 30, 30, "A"),
		CreateDefaultCalendarEvent(tomorrow, 9, 0, 30, "B"),
		CreateDefaultCalendarEvent(tomorrow, 9, 30, 30, "C"),
		CreateDefaultCalendarEvent(tomorrow, 10, 0, 30, "D"),
		CreateDefaultCalendarEvent(tomorrow, 13, 0, 60, "E"),
	}
	dailyAgenda0 := DailyAgenda{
		Date:   today,
		Events: calendarEvents0,
	}
	dailyAgenda1 := DailyAgenda{
		Date:   tomorrow,
		Events: calendarEvents1,
	}
	dailyAgendas := []DailyAgenda{
		dailyAgenda0,
		dailyAgenda1,
	}
	filledInAgendas, err := FillInWithEmptyDays(dailyAgendas, today, 2, false)
	if err != nil {
		t.Errorf("Error while filling in %v", err)
		return
	}
	if len(filledInAgendas) != 2 {
		t.Errorf("Length mismatch about no. agendas: %v", len(filledInAgendas))
		return
	}
	todaysEvents := filledInAgendas[0].Events
	if len(todaysEvents) != 1 {
		t.Errorf("Length mismatch about no. events: %v", len(todaysEvents))
		return
	}
	todaysEvents = filledInAgendas[1].Events
	if len(todaysEvents) != 6 {
		t.Errorf("Length mismatch about no. events: %v", len(todaysEvents))
		return
	}
}
