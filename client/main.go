package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/myoan/snake/api"
	"github.com/myoan/snake/engine"
)

const (
	Width  = 80
	Height = 30

	SceneTypeNone engine.SceneType = iota
	SceneTypeMenu
	SceneTypeMatchmaking
	SceneTypeIngame

	StatusInit = iota
	StatusStart
	StatusDrop
)

var (
	logger *log.Logger
)

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

	ui := NewUserInterface("noname", event, webEvent)

	ui.AddHandler(api.GameStatusInit, func(message []byte) error {
		logger.Printf("get init response: %s", string(message))
		var resp api.InitResponse
		err = json.Unmarshal(message, &resp)
		if err != nil {
			logger.Println("unmarshal:", err)
			return err
		}
		ui.UUID = resp.ID
		return nil
	})
	ui.AddHandler(api.GameStatusOK, func(message []byte) error {
		var resp api.EventResponse
		err = json.Unmarshal(message, &resp)
		if err != nil {
			logger.Println("unmarshal:", err)
			return err
		}
		board := generateBoard(resp.Body.Width, resp.Body.Height, resp.Body.Board)
		ui.Draw(board)
		return nil
	})
	ui.AddHandler(api.GameStatusError, func(message []byte) error {
		var resp api.EventResponse
		err = json.Unmarshal(message, &resp)
		if err != nil {
			logger.Println("unmarshal:", err)
			return err
		}

		logger.Printf("return from ConnectWebsocket read handler: %d", api.GameStatusError)
		ui.Status = StatusDrop
		ui.Score = resp.Body.Players[0].Size
		return fmt.Errorf("error")
	})
	ui.AddHandler(api.GameStatusWaiting, func(message []byte) error {
		logger.Printf("Receive waiting event")
		return nil
	})

	mng.AddScene(SceneTypeMenu, NewMenuScene(input, ui))
	mng.AddScene(SceneTypeMatchmaking, NewMatchmakingScene(input, ui))
	mng.AddScene(SceneTypeIngame, NewIngameScene(input, ui))
	mng.SetInitialScene(SceneTypeMenu)

	mng.Execute()
	ui.Finish()
	logger.Printf("========== GAME FINISH ==========")
}
