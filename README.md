# freeslots
Shell application written in Go language to retrieve events from Google Calendar, display them and compute free slots in your agenda.

# freeslots

# How to use the application

## Install dependencies
* go get golang.org/x/oauth2
* go get golang.org/x/oauth2/google
* go get google.golang.org/api/calendar/v3
* go get github.com/alexflint/go-arg

## Get Google Calendar API credentials
* Go to the Google Cloud Console
* Create a new project (or select an existing one)
* Enable the Google Calendar API
* Create consent screen
* Create OAuth 2.0 credentials (Desktop app) 
* Download the credentials JSON file
* Save it as credentials.json in the same directory as your Go program

## Run the program as follows
```bash
Usage: mytestapps --useremail USEREMAIL [--showallevents] [--creds CREDS] [--token TOKEN] [--listen LISTEN] [--nodays NODAYS] [--minduration MINDURATION] [--from FROM] [--to TO] [--format FORMAT] [--skipweekends] [--startdate STARTDATE] [--showslotduration]

Options:
  --useremail USEREMAIL
                         Full user email of the requestor. Mandatory field
  --showallevents        If present, show all events, otherwise show only free slots among events
  --creds CREDS          credentials.json file from Google [default: credentials.json]
  --token TOKEN          token.json file created by this app with the auth token from Google [default: token.json]
  --listen LISTEN        server address and port to open to get token from Google auth process [default: localhost:8080]
  --nodays NODAYS        Number of days after today [default: 14]
  --minduration MINDURATION
                         Min duration of slots to search for [default: 60]
  --from FROM            From what time to start reporting free slots [default: 09:00]
  --to TO                To what time reporting free slots [default: 18:00]
  --format FORMAT        Output format. Can be: plain, html, markdown [default: plain]
  --skipweekends         If present, skip weekends
  --startdate STARTDATE
                         From what date to start reporting free slots. Format accepted: yyyy-MM-dd
  --showslotduration     If present, show the free slot duration
  --help, -h             display this help and exit


Examples:
go run . 

go run . --useremail sample@gmail.com --nodays 14 --from 09:00 --to 18:00

go run . --useremail sample@gmail.com --listen localhost:8080

go run . --useremail sample@gmail.com --creds credentials.json --token token.json 

go run . --useremail sample@gmail.com --startdate 2025-12-03

```

## First-time authentication

* The program will display a URL
* Open it in your browser
* Sign in with your Google account
* Grant all displayed permissions
* The application will call the local web server opened on localhost:80
* In case the port is wrong, make sure to change the code to point to a new port and reflect it in the "redirect_uris' key of the *-credentials.json file and rerun the application
