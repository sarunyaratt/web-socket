package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Message struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ws", func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Println("Error while upgrading connection:", err)
			return err
		}
		defer ws.Close()

		clients[ws] = true

		go func() {
			for {
				var msg Message
				if err := ws.ReadJSON(&msg); err != nil {
					log.Println("Error reading JSON:", err)
					delete(clients, ws)
					break
				}
				broadcast <- msg
			}
		}()

		for {
			msg := <-broadcast
			for client := range clients {
				if err := client.WriteJSON(msg); err != nil {
					log.Printf("WebSocket error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	})

	e.Logger.Fatal(e.Start(":3000"))
}
