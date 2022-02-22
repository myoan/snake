package main

import (
	"flag"
	"log"
	"os"

	"github.com/myoan/snake/engine"
)

const (
	Width  = 80
	Height = 30

	SceneTypeNone engine.SceneType = iota
	SceneTypeMenu
	SceneTypeIngame
)

var (
	logger *log.Logger
)

type EventRequest struct {
	Eventtype int `json:"eventtype"`
	ID        int `json:"id"`
}

type EventResponse struct {
	Status  int              `json:"status"`
	Board   []int            `json:"board"`
	Width   int              `json:"width"`
	Height  int              `json:"height"`
	Players []PlayerResponse `json:"players"`
}

type PlayerResponse struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Size      int `json:"size"`
	Direction int `json:"direction"`
}

type Board struct {
	board  [][]int
	width  int
	height int
}

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	f, err := os.OpenFile("log/client.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}
	logger = log.New(f, "", log.LstdFlags)
	defer func() {
		f.Sync()
		f.Close()
	}()

	logger.Printf("========== GAME START ==========")

	flag.Parse()
	log.SetFlags(0)

	ge := engine.NewGameEngine(10)
	mng := ge.SceneManager
	defer mng.Stop()

	event := ge.GetEventStream()
	input := ge.Input
	webEvent := make(chan engine.ControlEvent)

	ui := NewUserInterface(event, webEvent)

	mng.AddScene(SceneTypeMenu, NewMenuScene(input, ui))
	mng.AddScene(SceneTypeIngame, NewIngameScene(input, ui))
	mng.SetInitialScene(SceneTypeMenu)

	mng.Execute()
	ui.Finish()
	logger.Printf("========== GAME FINISH ==========")
}
