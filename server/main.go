package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/myoan/snake/api"
)

const (
	Width    = 80
	Height   = 40
	InitX    = 20
	InitY    = 20
	InitSize = 3
)

const (
	EventClientConnect = iota
	EventClientFinish
	EventClientRestart
)

type Observer interface {
	Update(data interface{}) error
}

type TriggerArgument struct {
	EventType int
	Client    Client
}

func ingameHandler(mng *SceneManager, w http.ResponseWriter, r *http.Request) {
	log.Printf("receive new ingame handler")
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	stream := make(chan []byte)
	NewWebClient(mng, c, stream)
	// TODO: ここでNewClientしているのは、mngをObserverとして登録してNotifyするためなのだが、分かりづらい

	// player := NewPlayer(client, stream)
	// event := make(chan Event)
	// game := NewGame(Width, Height, event, player)
	// go game.Run()
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ingame")
}

func main() {
	addr := flag.String("addr", "localhost:8080", "http service address")
	ge := NewGameEngine()

	ge.SceneMng.AddTrigger(EventClientConnect, func(args interface{}) {
		ta := args.(TriggerArgument)

		switch ge.SceneMng.CurrentSceneID {
		case SceneMatchmaking:
			log.Printf("Scene: MatchMaking (%d)\n", len(ge.Clients))
			ge.AddClient(ta.Client)
			if ge.ReachMaxClient() {
				ge.SceneMng.MoveScene(SceneIngame)

				ge.ExecuteIngame()
			} else {
				data := &api.EventResponse{
					Status: api.GameStatusWaiting,
				}

				bytes, _ := json.Marshal(&data)
				ta.Client.Send(bytes)
			}
		case SceneIngame:
			log.Printf("Scene: Ingame, ignore\n")
			// TODO: Should I disconnect client?
		}
	})
	ge.SceneMng.AddTrigger(EventClientFinish, func(args interface{}) {
		// ta := args.(TriggerArgument)
		log.Printf("Trigger: EventClientFinish\n")
		ge.DeleteClient(1)
	})
	ge.SceneMng.AddTrigger(EventClientRestart, func(args interface{}) {
		// ta := args.(TriggerArgument)
		log.Printf("Trigger: EventClientRestart\n")
	})

	flag.Parse()
	http.HandleFunc("/ingame", func(w http.ResponseWriter, r *http.Request) {
		ingameHandler(ge.SceneMng, w, r)
	})
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
