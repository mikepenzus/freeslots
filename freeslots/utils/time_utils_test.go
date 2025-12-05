package utils

import (
	"testing"
	"time"
)

func TestGetEndTime(t *testing.T) {
	timeNow := time.Now()
	startEvent := CreateDefaultCalendarEventFromString(timeNow, "08:02", 60, "X")
	endTime := startEvent.GetEndTime()
	if endTime.Hour() != 9 || endTime.Minute() != 2 {
		t.Errorf("Error while parsing %v", endTime)
	}
	startEvent = CreateDefaultCalendarEventFromString(timeNow, "22:51", 42, "Descr")
	endTime = startEvent.GetEndTime()
	if endTime.Hour() != 23 || endTime.Minute() != 33 {
		t.Errorf("Error while parsing %v", endTime)
	}
	startEvent = CreateDefaultCalendarEventFromString(timeNow, "0:1", 121, "DescYr")
	endTime = startEvent.GetEndTime()
	if endTime.Hour() != 2 || endTime.Minute() != 2 {
		t.Errorf("Error while parsing %v", endTime)
	}
}

func TestCreateDefaultCalendarEvent(t *testing.T) {
	now := time.Now()
	event := CreateDefaultCalendarEvent(now, 8, 2, 60, "X")
	if event.Duration != 60 || event.Description != "X" || event.StartTime.Hour() != 8 ||
		event.StartTime.Minute() != 2 {
		t.Errorf("Error while parsing %v", event)
	}
	event = CreateDefaultCalendarEvent(now, 22, 51, 42, "Descr")
	if event.Duration != 42 || event.Description != "Descr" || event.StartTime.Hour() != 22 ||
		event.StartTime.Minute() != 51 {
		t.Errorf("Error while parsing %v", event)
	}
	event = CreateDefaultCalendarEvent(now, 0, 1, 121, "Y")
	if event.Duration != 121 || event.Description != "Y" || event.StartTime.Hour() != 0 ||
		event.StartTime.Minute() != 1 {
		t.Errorf("Error while parsing %v", event)
	}
}
func TestCreateDefaultCalendarEventFromString(t *testing.T) {
	timeNow := time.Now()
	event := CreateDefaultCalendarEventFromString(timeNow, "08:02", 60, "X")
	if event.Duration != 60 || event.Description != "X" || event.StartTime.Hour() != 8 ||
		event.StartTime.Minute() != 2 {
		t.Errorf("Error while parsing %v", event)
	}
	event = CreateDefaultCalendarEventFromString(timeNow, "22:51", 42, "Descr")
	if event.Duration != 42 || event.Description != "Descr" || event.StartTime.Hour() != 22 ||
		event.StartTime.Minute() != 51 {
		t.Errorf("Error while parsing %v", event)
	}
	event = CreateDefaultCalendarEventFromString(timeNow, "0:1", 121, "Y")
	if event.Duration != 121 || event.Description != "Y" || event.StartTime.Hour() != 0 ||
		event.StartTime.Minute() != 1 {
		t.Errorf("Error while parsing %v", event)
	}
}

func TestParseSingleDayAgenda(t *testing.T) {
	timeNow := time.Now()
	agendas := []string{
		"d2025-08-30,m30,aX",
		"d2023-09-11,m60,aX",
	}
	expectedDates := []time.Time{
		time.Date(2025, time.August, 30, 0, 0, 0, 0, time.Local),
		time.Date(2023, time.September, 11, 0, 0, 0, 0, time.Local),
	}
	for agendaIndex, agenda := range agendas {
		resultingEvents, err := ParseSingleDayAgenda(agenda)
		if err != nil {
			t.Errorf("Error while parsing agenda %v", agenda)
			return
		}
		resultingDate := resultingEvents[0].StartTime
		expectedDate := expectedDates[agendaIndex]
		if CompareTimesWithoutTimeZone(resultingDate, expectedDate) != 0 {
			t.Errorf("Error while parsing single day %v %v %v", agenda, resultingDate, expectedDate)
		}
	}

	agendas = []string{
		"m30,s16,aXX--XY--XYYX--ZZZZZZ",
		"m60,s4,aX-Y-Z-A-B-C-D---",
		"m30,s4,aX-Y-Z-A-B-C-D---",
		"m60,s8,aXXXXYYYYZZZZAAAA",
		"m60,s8,aXXXXXXXXZZZZZZZZ",
		"m30,s16,aXXXXYYYYZZZZAAAA",
		"m30,s16,aXXXXXXXXZZZZZZZZ",
		"m120,aXXXXXXXYYYYY",
		"m60,s0,a--------------------",
		"m15,s32,aX-Y-Z-A-B-C-D---",
		"m60,s4,a-------------AAAA",
		"m30,aAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"m30,s16,aXX--XY--XYYX--ZZZZZZ,aA",
		"m30,s16,aXX--XY--XYYX--ZZZZZZ,a-AA",
		"m30,s16,a-XX--XY--XYYX--ZZZZZZ,aAA",
	}
	expectedEventLists := [][]CalendarEvent{
		{
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 60, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "10:00", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "10:30", 30, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "12:00", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "12:30", 60, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "13:30", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "15:00", 180, "Z"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "04:00", 60, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "06:00", 60, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 60, "Z"),
			CreateDefaultCalendarEventFromString(timeNow, "10:00", 60, "A"),
			CreateDefaultCalendarEventFromString(timeNow, "12:00", 60, "B"),
			CreateDefaultCalendarEventFromString(timeNow, "14:00", 60, "C"),
			CreateDefaultCalendarEventFromString(timeNow, "16:00", 60, "D"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "02:00", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "03:00", 30, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "04:00", 30, "Z"),
			CreateDefaultCalendarEventFromString(timeNow, "05:00", 30, "A"),
			CreateDefaultCalendarEventFromString(timeNow, "06:00", 30, "B"),
			CreateDefaultCalendarEventFromString(timeNow, "07:00", 30, "C"),
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 30, "D"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 240, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "12:00", 240, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "16:00", 240, "Z"),
			CreateDefaultCalendarEventFromString(timeNow, "20:00", 240, "A"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 480, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "16:00", 480, "Z"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 120, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "10:00", 120, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "12:00", 120, "Z"),
			CreateDefaultCalendarEventFromString(timeNow, "14:00", 120, "A"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 240, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "12:00", 240, "Z"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "00:00", 840, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "14:00", 600, "Y"),
		},
		{},
		{
			CreateDefaultCalendarEventFromString(timeNow, "8:00", 15, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "8:30", 15, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "9:00", 15, "Z"),
			CreateDefaultCalendarEventFromString(timeNow, "9:30", 15, "A"),
			CreateDefaultCalendarEventFromString(timeNow, "10:00", 15, "B"),
			CreateDefaultCalendarEventFromString(timeNow, "10:30", 15, "C"),
			CreateDefaultCalendarEventFromString(timeNow, "11:00", 15, "D"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "17:00", 240, "A"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "0:0", 1440, "A"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 60, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 30, "A"),
			CreateDefaultCalendarEventFromString(timeNow, "10:00", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "10:30", 30, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "12:00", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "12:30", 60, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "13:30", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "15:00", 180, "Z"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 60, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "08:30", 60, "A"),
			CreateDefaultCalendarEventFromString(timeNow, "10:00", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "10:30", 30, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "12:00", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "12:30", 60, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "13:30", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "15:00", 180, "Z"),
		},
		{
			CreateDefaultCalendarEventFromString(timeNow, "08:00", 60, "A"),
			CreateDefaultCalendarEventFromString(timeNow, "08:30", 60, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "10:30", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "11:00", 30, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "12:30", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "13:00", 60, "Y"),
			CreateDefaultCalendarEventFromString(timeNow, "14:00", 30, "X"),
			CreateDefaultCalendarEventFromString(timeNow, "15:30", 180, "Z"),
		},
	}
	for agendaIndex, agenda := range agendas {
		resultingEvents, _ := ParseSingleDayAgenda(agenda)
		expectedEvents := expectedEventLists[agendaIndex]
		if len(resultingEvents) != len(expectedEvents) {
			t.Errorf("Error while parsing: agenda index %v, agenda %v, no. expected %v, no. results %v. Different lengths",
				agendaIndex, agenda, len(expectedEvents), len(resultingEvents))
			PrintEventList(resultingEvents)
			return
		}
		for eventIndex, expectedEvent := range expectedEvents {
			resultingEvent := resultingEvents[eventIndex]
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
