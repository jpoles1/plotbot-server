package main

import (
	"fmt"
	"log"
	"net/http"
	"plotbot-server/envload"
	"plotbot-server/logging"
	"strconv"

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

var plotterClients = make(map[*websocket.Conn]plotterStatus)
var plotterBroadcasters = make(map[bson.ObjectId]*websocket.Conn)
var monitorClients = make(map[*websocket.Conn]bool)
var broadcast = make(chan wsMessage) // broadcast channel

var defaultPlotterConfig = plotterConfig{
	AnchorDistance: 39 * 2.54,
	SpoolDiameter:  30,
	StartCoord:     plotterCoordinate{20 * 2.54, 29 * 2.54},
}

func wsClose(ws *websocket.Conn) {
	fmt.Println("Closing socket!")
	if _, ok := plotterClients[ws]; ok {
		delete(plotterBroadcasters, plotterClients[ws].PlotterID)
	}
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
	router.Route("/api", func(r chi.Router) {
		r.Get("/deviceList", func(w http.ResponseWriter, r *http.Request) {
			var clientData []plotterStatus
			for _, value := range plotterClients {
				clientData = append(clientData, value)
			}
			sendResponseJSON(w, clientData)
		})
		r.Get("/deviceCommand", func(w http.ResponseWriter, r *http.Request) {
			queryValues := r.URL.Query()
			if bson.IsObjectIdHex(queryValues.Get("id")) {
				if x, err := strconv.ParseFloat(queryValues.Get("x"), 64); err == nil {
					if y, err := strconv.ParseFloat(queryValues.Get("y"), 64); err == nil {
						ws := plotterBroadcasters[bson.ObjectIdHex(queryValues.Get("id"))]
						updatedClient := plotterClients[ws]
						ws.WriteJSON(wsMessage{
							MessageType: "plotCommand",
							Payload: updatedClient.generatePlotMessage(
								plotterCoordinate{x, y},
							),
						})
						plotterClients[ws] = updatedClient
					}
				}
			}
		})
	})
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
			fmt.Println(msg)
			if msg.MessageType == "registration" {
				if msg.Payload.(string) == "plotter" {
					delete(monitorClients, ws)
					plotterClients[ws] = plotterStatus{
						bson.NewObjectId(),
						defaultPlotterConfig,
						"Online",
						plotterCoordinate{0, 0},
					}
					plotterBroadcasters[plotterClients[ws].PlotterID] = ws
					ws.WriteJSON(wsMessage{
						MessageType: "status",
						Payload:     "New Plotter ID Assigned: " + plotterClients[ws].PlotterID.Hex(),
					})
				}
			}
			if msg.MessageType == "commandRequest" {
				if _, ok := plotterClients[ws]; ok {
					currentCommand := int(msg.Payload.(float64))
					if currentCommand == 0 {
						updatedStatus := plotterClients[ws]
						updatedStatus.CurrentCoord = plotterCoordinate{0, 0}
						plotterClients[ws] = updatedStatus
					}
					if currentCommand < len(coordList) {
						updatedClient := plotterClients[ws]
						ws.WriteJSON(wsMessage{
							MessageType: "plotCommand",
							Payload:     updatedClient.generatePlotMessage(coordList[currentCommand]),
						})
						plotterClients[ws] = updatedClient
					}
				}
			}
		}
	})
	color.Green("Starting web server at: http://%s:%s", envData.BindIP, envData.BindPort)
	log.Fatal(http.ListenAndServe(envData.BindIP+":"+envData.BindPort, router))
}
