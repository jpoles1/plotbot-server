package main

import (
	"log"
	"net/http"
	"plotbot-server/envload"
	"plotbot-server/logging"

	"github.com/fatih/color"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var plotterClients = make(map[*websocket.Conn]bson.ObjectId)
var monitorClients = make(map[*websocket.Conn]bool)
var broadcast = make(chan wsMessage) // broadcast channel

func wsClose(ws *websocket.Conn) {
	delete(plotterClients, ws)
	delete(monitorClients, ws)
	ws.Close()
}
func main() {
	envData, err := envload.LoadEnv(".env")
	if err != nil {
		log.Fatal("Failed to load environment config.")
	}
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Get("/", fileServer("static", true))
	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		if err != nil {
			logging.Error("Upgrading websocket connection", err)
		}
		defer wsClose(ws)
		monitorClients[ws] = true
		for {
			// Read message from browser
			var msg wsMessage
			err := ws.ReadJSON(&msg)
			if err != nil {
				logging.Error("Parsing incoming JSON", err)
				break
			}
			if msg.MessageType == 0 {
				if msg.Payload.(registrationMessage).ClientType == plotter {
					delete(monitorClients, ws)
					plotterClients[ws] = bson.NewObjectId()
					ws.WriteJSON(wsMessage{
						MessageType: status,
						Payload:     "New Plotter ID Assigned: " + plotterClients[ws].Hex(),
					})
				}
			}
			if msg.MessageType == 2 {
				if _, ok := plotterClients[ws]; ok {
					currentCommand := msg.Payload.(int)
					if currentCommand < len(commandList) {
						ws.WriteJSON(wsMessage{
							MessageType: plotCommand,
							Payload:     commandList[currentCommand],
						})
					}
				}
			}
		}
	})
	color.Green("Starting web server at: http://%s:%s", envData.BindIP, envData.BindPort)
	log.Fatal(http.ListenAndServe(envData.BindIP+":"+envData.BindPort, router))
}
