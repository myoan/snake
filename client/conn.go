package main

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/myoan/snake/api"
)

const (
	StatusInit = iota
	StatusWait
	StatusStart
	StatusDrop
)

type Conn struct {
	conn    *websocket.Conn
	event   chan int
	webDone chan struct{}
	funcMap map[int]func([]byte) error
	UUID    string
}

func NewConn() *Conn {
	event := make(chan int)
	done := make(chan struct{})
	fm := make(map[int]func([]byte) error)
	return &Conn{
		webDone: done,
		event:   event,
		funcMap: fm,
	}
}

func (conn *Conn) Connect(addr string) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	conn.conn = c

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			var resp map[string]interface{}
			err = json.Unmarshal(message, &resp)
			if err != nil {
				return
			}
			fstatus := resp["status"].(float64)
			istatus := int(fstatus)

			err = conn.funcMap[istatus](message)
			if err != nil {
				return
			}
		}
	}()

	for {
		select {
		case ctrl := <-conn.event:
			event := &api.EventRequest{
				UUID:      conn.UUID,
				Eventtype: 0,
				Key:       ctrl,
			}
			bytes, _ := json.Marshal(&event)
			err := c.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				return
			}
		case <-conn.webDone:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return
			}
			select {
			case <-conn.webDone:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// AddHandler adds a handler when server response is received.
func (conn *Conn) AddHandler(handler int, fn func([]byte) error) {
	conn.funcMap[handler] = fn
}

// CloseWebSocket closes disconnects to server.
// It is called when you exit ingame.
func (conn *Conn) Close() {
	conn.webDone <- struct{}{}

	// TODO: Does it exist other good way? (ex. wait for server response)
	time.Sleep(300 * time.Millisecond)

	conn.conn.Close()
}
