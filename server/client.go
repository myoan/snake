package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type WebClient struct {
	conn *websocket.Conn
}

func NewWebClient(conn *websocket.Conn, dataStream chan []byte) *WebClient {
	client := &WebClient{
		conn: conn,
	}
	go client.run(dataStream)
	return client
}

func (c *WebClient) ID() int {
	return 1
}

func (c *WebClient) Send(data []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("write: ", err)
		return err
	}
	return nil
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
			c.Close()
			return
		}
		log.Printf("recv: %s", message)

		stream <- message
	}
}

func (c *WebClient) Close() {
	c.conn.Close()
}
