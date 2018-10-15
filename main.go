package main

import (
	"fmt"
	"log"
	"net/http"
	"plotbot-server/envload"
	"plotbot-server/logging"

	"github.com/fatih/color"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

type stepperCommand struct {
	Degrees   int  `json:"deg"`
	Clockwise bool `json:"dir"`
}

type plotCommand struct {
	LeftStepper  stepperCommand `json:"left"`
	RightStepper stepperCommand `json:"right"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
	router.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test response"))
		})
	})
	router.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
			if err != nil {
				logging.Error("Upgrading websocket connection", err)
			}
			for {
				// Read message from browser
				msgType, msg, err := conn.ReadMessage()
				if err != nil {
					return
				}

				// Print the message to the console
				fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

				// Write message back to browser
				if err = conn.WriteMessage(msgType, msg); err != nil {
					return
				}
			}
		})
	})
	color.Green("Starting web server at: http://%s:%s", envData.BindIP, envData.BindPort)
	log.Fatal(http.ListenAndServe(envData.BindIP+":"+envData.BindPort, router))
}
