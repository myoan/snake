// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	Width  = 80
	Height = 40
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var logger *log.Logger

type EventRequest struct {
	Eventtype int `json:"eventtype"`
	ID        int `json:"id"`
}

type PlayerResponse struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Size      int `json:"size"`
	Direction int `json:"direction"`
}

type EventResponse struct {
	Status  int              `json:"status"`
	Board   []int            `json:"board"`
	Width   int              `json:"width"`
	Height  int              `json:"height"`
	Players []PlayerResponse `json:"players"`
}

type WebClient struct {
	conn *websocket.Conn
}

var client *WebClient
var Count int
var Delta int

func (c *WebClient) Send() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		Count += Delta
		event := &EventResponse{}
		bytes, _ := json.Marshal(&event)
		logger.Println("write:", string(bytes))
		err := c.conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			logger.Println("write:", err)
			return
		}
	}
}

func ingameHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Print("upgrade:", err)
		return
	}
	client = &WebClient{
		conn: c,
	}
	defer c.Close()

	// --- create game engine ---

	event := make(chan Event)
	game := NewGame(Width, Height, event)

	go game.Run()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			logger.Println("read:", err)
			break
		}
		logger.Printf("recv: %s", message)

		var req EventRequest
		json.Unmarshal(message, &req)

		err = game.changeDirection(req.ID)
		if err != nil {
			logger.Println("read:", err)
			return
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ingame")
}

func main() {
	f, err := os.OpenFile("log/server.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}
	logger = log.New(f, "", log.LstdFlags)
	defer func() {
		f.Sync()
		f.Close()
	}()

	Count = 0
	Delta = 1
	flag.Parse()
	logger.SetFlags(0)
	http.HandleFunc("/ingame", ingameHandler)
	http.HandleFunc("/", home)
	logger.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
