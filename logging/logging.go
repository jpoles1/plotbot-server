package logging

import (
	"log"

	raven "github.com/getsentry/raven-go"
)

//LoadSentry initializes the connection to the Sentry error logging server
func LoadSentry(SentryDSN string) {
	raven.SetDSN(SentryDSN)
}

//Error is a helper function to check if an error has been produced, and respond appropriately
func Error(taskDescription string, err error) {
	if err != nil {
		log.Println("Error w/ " + taskDescription + ": " + err.Error())
		raven.CaptureError(err, map[string]string{"msg": taskDescription})
	}
}

//Fatal is a helper function to check if a fatal error has been produced, and respond appropriately
func Fatal(taskDescription string, err error) {
	if err != nil {
		log.Fatal("Fatal Error w/ " + taskDescription + ": " + err.Error())
		raven.CaptureError(err, map[string]string{"msg": taskDescription})
	}
}
