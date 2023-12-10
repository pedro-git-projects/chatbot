package app

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func (app *Application) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := app.upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.logger.Println(err)
		return
	}

	app.clients[conn] = true

	app.logger.Println("Client connected")
}

func (app *Application) handleConnections() {
	disconnectCh := make(chan *websocket.Conn)

	go func() {
		for conn := range app.clients {
			_, _, err := conn.ReadMessage()
			if err != nil {
				disconnectCh <- conn
			}
		}
	}()

	for {
		select {
		case conn := <-disconnectCh:
			app.logger.Println("Client disconnected")
			conn.Close()
			delete(app.clients, conn)
		}
	}
}

func (app *Application) handleMessages() {
	for {
		msg := <-app.broadcast
		app.logger.Printf("Broadcasting message: %s\n", msg)
		for client := range app.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				app.logger.Println(err)
				client.Close()
				delete(app.clients, client)
			}
		}
	}
}
