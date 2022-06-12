package main

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/myoan/snake/api"
	"github.com/myoan/snake/engine"
)

type Direction int

const (
	DirectionLeft Direction = iota
	DirectionRight
	DirectionUp
	DirectionDown
)

/*
UserInterface represents the user interface, screen and controller.
This struct must be initialized at the beginning of the program and must live until the end.
You shouldn't re-create this struct.
*/

type UserInterface struct {
	Status   int
	webEvent chan engine.ControlEvent
	webDone  chan struct{}
	conn     *websocket.Conn
	funcMap  map[int]func([]byte) error
	Score    int
	UUID     string
}

// NewUserInterface creates a new UserInterface.
// You must call this method before using the UserInterface.
// UserInterface is listening user controlle events and sends them to the event channel.
// webEvent is a channel for sending to web server.
func NewUserInterface(uuid string, event chan<- engine.ControlEvent, webEvent chan engine.ControlEvent) *UserInterface {
	done := make(chan struct{})
	fm := make(map[int]func([]byte) error)

	ui := &UserInterface{
		Status:   StatusInit,
		Score:    -1,
		UUID:     uuid,
		webEvent: webEvent,
		webDone:  done,
		funcMap:  fm,
	}
	return ui
}

// Finish is called when the entire game is over.
func (ui *UserInterface) Finish() {}

// AddHandler adds a handler when server response is received.
func (ui *UserInterface) AddHandler(handler int, fn func([]byte) error) {
	ui.funcMap[handler] = fn
}

// ConnectWebSocket connects to server.
// It connects when ingame is started.
// So, it is recreate connections if you play ingame multiple times.
func (ui *UserInterface) ConnectWebSocket() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ingame"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("dial: %v", err)
		return
	}
	ui.conn = c

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var resp map[string]interface{}
			err = json.Unmarshal(message, &resp)
			if err != nil {
				log.Println("unmarshal:", err)
				return
			}
			fstatus := resp["status"].(float64)
			istatus := int(fstatus)

			err = ui.funcMap[istatus](message)
			if err != nil {
				log.Printf("return from ConnectWebsocket read handler: %d", istatus)
				return
			}
		}
	}()

	for {
		select {
		case ctrl := <-ui.webEvent:
			event := &api.EventRequest{
				UUID:      ui.UUID,
				Eventtype: ctrl.Eventtype,
				Key:       ctrl.Key,
			}
			bytes, _ := json.Marshal(&event)
			err := c.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-ui.webDone:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-ui.webDone:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// CloseWebSocket closes disconnects to server.
// It is called when you exit ingame.
func (ui *UserInterface) CloseWebSocket() {
	ui.webDone <- struct{}{}

	// TODO: Does it exist other good way? (ex. wait for server response)
	time.Sleep(300 * time.Millisecond)

	ui.Status = StatusInit
	ui.conn.Close()
}

func (ui *UserInterface) Update(players []api.PlayerResponse) {
	var x, y int
	var dir Direction
	for _, player := range players {
		if ui.UUID == player.ID {
			x = player.X
			y = player.Y
			dir = Direction(player.Direction)
		}
	}

	if x == 3 && dir == DirectionLeft {
		ui.webEvent <- engine.ControlEvent{Eventtype: 0, Key: api.MoveDown}
	}
	if x == Width-4 && dir == DirectionRight {
		ui.webEvent <- engine.ControlEvent{Eventtype: 0, Key: api.MoveUp}
	}
	if y == 3 && dir == DirectionUp {
		ui.webEvent <- engine.ControlEvent{Eventtype: 0, Key: api.MoveLeft}
	}
	if y == Height-4 && dir == DirectionDown {
		ui.webEvent <- engine.ControlEvent{Eventtype: 0, Key: api.MoveRight}
	}
}
