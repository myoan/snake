package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	Width    = 80
	Height   = 40
	InitX    = 20
	InitY    = 20
	InitSize = 3
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

var client *WebClient

type WebClient struct {
	conn *websocket.Conn
}

func NewWebClient(conn *websocket.Conn, dataStream chan []byte) *WebClient {
	client = &WebClient{
		conn: conn,
	}
	go client.run(dataStream)
	return client
}

func (c *WebClient) ID() int {
	return 1
}

func (c *WebClient) Send(data []byte) error {
	err := client.conn.WriteMessage(websocket.TextMessage, data)
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

func ingameHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	stream := make(chan []byte)
	client = NewWebClient(c, stream)

	player := NewPlayer(stream)

	// --- create game engine ---

	event := make(chan Event)
	game := NewGame(Width, Height, event, player)

	go game.Run()
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ingame")
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	flag.Parse()
	http.HandleFunc("/ingame", ingameHandler)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
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
