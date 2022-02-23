package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WebClient struct {
	stream    chan []byte
	conn      *websocket.Conn
	observers []Observer
	mu        sync.Mutex
}

func (c *WebClient) AddObserver(o Observer) {
	c.observers = append(c.observers, o)
}

func NewWebClient(mng *SceneManager, conn *websocket.Conn, dataStream chan []byte) *WebClient {
	obs := make([]Observer, 0)
	client := &WebClient{
		stream:    dataStream,
		conn:      conn,
		observers: obs,
	}
	go client.run(dataStream)
	client.AddObserver(mng)
	client.Notify(EventClientConnect)
	return client
}

func (c *WebClient) Notify(tp int) {
	for _, o := range c.observers {
		data := TriggerArgument{
			EventType: tp,
			Client:    c,
		}
		o.Update(data)
	}
}

func (c *WebClient) ID() int {
	return 1
}

func (c *WebClient) Send(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("write: ", err)
		return err
	}
	return nil
}

func (c *WebClient) Stream() chan []byte {
	return c.stream
}

func (c *WebClient) run(stream chan []byte) {
	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			close(stream)
			c.Close()
			return
		}
		if mt == websocket.CloseMessage {
			log.Println("close:", string(message))
			close(stream)
			c.Notify(EventClientFinish)
			c.Close()
			return
		}
		log.Printf("recv: %s", message)

		stream <- message
	}
}

func (c *WebClient) Close() {
	log.Printf("Close client %d", c.ID())
	c.conn.Close()
}
