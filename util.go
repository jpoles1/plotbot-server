package main

import (
	"encoding/json"
	"net/http"
	"plotbot-server/logging"
)

func sendErrorCode(w http.ResponseWriter, errCode int, msg string, err error) {
	http.Error(w, http.StatusText(errCode)+" - "+msg, errCode)
	if err != nil {
		logging.Error(msg, err)
	}
}

func sendResponseJSON(w http.ResponseWriter, data interface{}) (err error) {
	responseString, err := json.Marshal(data)
	if err != nil {
		sendErrorCode(w, 500, "Marshalling data to json", err)
		return
	}
	w.Write(responseString)
	return
}
