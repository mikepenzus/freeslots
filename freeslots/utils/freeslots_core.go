package utils

import (
	"fmt"
	"log"
	"time"
)

type FreeSlotsCoreAlgorithm struct {
	ShowAllEvents    bool
	NoDays           int
	MinDuration      int
	FromTime         string
	ToTime           string
	Format           string
	SkipWeekends     bool
	StartDate        time.Time
	ShowSlotDuration bool
}

func (freeSlotsCoreAlgorithm FreeSlotsCoreAlgorithm) FreeSlotsCore(dailyAgendas []DailyAgenda) {
	if freeSlotsCoreAlgorithm.ShowAllEvents {
		freeSlotsCoreAlgorithm.PrintAllEvents(dailyAgendas)
	} else {
		freeSlotsCoreAlgorithm.PrintFreeSlots(dailyAgendas)
	}
}

func (freeSlotsCoreAlgorithm FreeSlotsCoreAlgorithm) PrintAllEvents(dailyAgendas []DailyAgenda) {
	switch freeSlotsCoreAlgorithm.Format {
	case "plain":
		for _, dailyAgenda := range dailyAgendas {
			if freeSlotsCoreAlgorithm.SkipWeekends && dailyAgenda.IsWeekend() {
				continue
			}
			dailyAgenda.Print(true, false)
		}
	case "html":
		fmt.Print("<html><style>table, th, td {  border: 1px solid black;  border-collapse: collapse;} </style> <body><table><tr><td>Date</td><td>Event</td><td>Description</td></tr>")
		for _, dailyAgenda := range dailyAgendas {
			if freeSlotsCoreAlgorithm.SkipWeekends && dailyAgenda.IsWeekend() {
				continue
			}
			dailyAgenda.PrintHtml(true, false)
		}
		fmt.Println("</table></body></html>")
	case "markdown":
		fmt.Println("| Date | Event | Description |")
		fmt.Println("| -------- | -------- | -------- |")
		for _, dailyAgenda := range dailyAgendas {
			if freeSlotsCoreAlgorithm.SkipWeekends && dailyAgenda.IsWeekend() {
				continue
			}
			dailyAgenda.PrintMarkdown(true, false)
		}
	}
}

func (freeSlotsCoreAlgorithm FreeSlotsCoreAlgorithm) PrintFreeSlots(dailyAgendas []DailyAgenda) {
	newDailyAgendas, err := FillInWithEmptyDays(dailyAgendas, freeSlotsCoreAlgorithm.StartDate,
		freeSlotsCoreAlgorithm.NoDays, freeSlotsCoreAlgorithm.SkipWeekends)
	if err != nil {
		log.Fatalf("Unable to create event lists for empty days: %v", err)
	}
	fromHours, fromMinutes := ParseTime(freeSlotsCoreAlgorithm.FromTime)
	toHours, toMinutes := ParseTime(freeSlotsCoreAlgorithm.ToTime)
	switch freeSlotsCoreAlgorithm.Format {
	case "html":
		fmt.Print("<html><style>table, th, td {  border: 1px solid black;  border-collapse: collapse;} </style> <body><table><tr><td>Date</td><td>Slot</td></tr>")
	case "markdown":
		fmt.Println("| Date | Slot |")
		fmt.Println("| ----------- | ----------- |")
	}
	for _, dailyAgenda := range newDailyAgendas {
		freeSlotsAgenda, err := dailyAgenda.GetFreeSlots(freeSlotsCoreAlgorithm.MinDuration,
			fromHours, fromMinutes, toHours, toMinutes)
		if err != nil {
			log.Fatalf("Unable to get free slots: %v", err)
			return
		}
		if freeSlotsAgenda.IsEmpty() {
			continue
		}
		switch freeSlotsCoreAlgorithm.Format {
		case "html":
			freeSlotsAgenda.PrintHtml(false, freeSlotsCoreAlgorithm.ShowSlotDuration)
		case "plain":
			freeSlotsAgenda.Print(false, freeSlotsCoreAlgorithm.ShowSlotDuration)
		case "markdown":
			freeSlotsAgenda.PrintMarkdown(false, freeSlotsCoreAlgorithm.ShowSlotDuration)
		}
	}
	if freeSlotsCoreAlgorithm.Format == "html" {
		fmt.Println("</table></body></html>")
	}
}
