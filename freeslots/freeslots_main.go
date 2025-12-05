package main

import (
	"log"
	"time"

	"freeslots/utils"

	"github.com/alexflint/go-arg"
)

type InputArgs struct {
	UserEmail               string `arg:"--useremail,required" help:"Full user email of the requestor. Mandatory field"`
	ShowAllEvents           bool   `arg:"--showallevents" help:"If present, show all events, otherwise show only free slots among events"`
	CredentialsFileName     string `arg:"--creds" default:"credentials.json" help:"credentials.json file from Google"`
	TokenFileName           string `arg:"--token" default:"token.json" help:"token.json file created by this app with the auth token from Google"`
	WebserverAddressAndPort string `arg:"--listen" default:"localhost:8080" help:"server address and port to open to get token from Google auth process"`
	NoDays                  int    `arg:"--nodays" default:"14" help:"Number of days after today"`
	MinDuration             int    `arg:"--minduration" default:"60" help:"Min duration of slots to search for"`
	FromTime                string `arg:"--from" default:"09:00" help:"From what time to start reporting free slots"`
	ToTime                  string `arg:"--to" default:"18:00" help:"To what time reporting free slots"`
	Format                  string `arg:"--format" default:"plain" help:"Output format. Can be: plain, html, markdown"`
	SkipWeekends            bool   `arg:"--skipweekends" help:"If present, skip weekends"`
	StartDate               string `arg:"--startdate" default:"" help:"From what date to start reporting free slots. Format accepted: yyyy-MM-dd"`
	ShowSlotDuration        bool   `arg:"--showslotduration" help:"If present, show the free slot duration"`
}

func main() {
	var inputArgs InputArgs
	arg.MustParse(&inputArgs)
	calendarExporterStatus := utils.CalendarExporterStatus{}
	calendarExporterStatus.CredentialsFileName = inputArgs.CredentialsFileName
	calendarExporterStatus.TokenFileName = inputArgs.TokenFileName
	calendarExporterStatus.WebserverAddressAndPort = inputArgs.WebserverAddressAndPort

	calendarService, err := utils.CreateCalendarService(calendarExporterStatus)
	if err != nil {
		log.Fatalf("Unable to create Google Calendar service: %v", err)
	}

	startDate := utils.GetPureDate(time.Now())
	if inputArgs.StartDate != "" {
		startDate, err = time.Parse(time.DateOnly, inputArgs.StartDate)
		if err != nil {
			log.Fatalf("Bad start date: %v", err)
		}
		// TODO: fix time zone
	}
	dailyAgendas, err := utils.GetEventsFromPrimaryCalendar(calendarService, startDate, inputArgs.NoDays, inputArgs.UserEmail)
	if err != nil {
		log.Fatalf("Unable to retrieve Google Calendar events: %v", err)
	}

	freeSlotsCoreAlgorithm := utils.FreeSlotsCoreAlgorithm{
		ShowAllEvents:    inputArgs.ShowAllEvents,
		NoDays:           inputArgs.NoDays,
		MinDuration:      inputArgs.MinDuration,
		FromTime:         inputArgs.FromTime,
		ToTime:           inputArgs.ToTime,
		Format:           inputArgs.Format,
		SkipWeekends:     inputArgs.SkipWeekends,
		StartDate:        startDate,
		ShowSlotDuration: inputArgs.ShowSlotDuration,
	}
	freeSlotsCoreAlgorithm.FreeSlotsCore(dailyAgendas)
}
