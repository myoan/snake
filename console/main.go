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
	Width  = 40
	Height = 40

	SceneTypeNone engine.SceneType = iota
	SceneTypeMenu
	SceneTypeMatchmaking
	SceneTypeIngame

	StatusInit = iota
	StatusStart
	StatusDrop
)

type Board struct{}

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	f, err := os.OpenFile("log/client.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		f.Sync()
		f.Close()
	}()

	log.Printf("========== GAME START ==========")

	flag.Parse()
	log.SetFlags(0)

	ge := engine.NewGameEngine(10)
	mng := ge.SceneManager
	defer mng.Stop()

	event := ge.GetEventStream()
	webEvent := make(chan engine.ControlEvent)

	ui := NewUserInterface("noname", event, webEvent)

	ui.AddHandler(api.GameStatusInit, func(message []byte) error {
		log.Printf("get init response: %s", string(message))
		var resp api.InitResponse
		err = json.Unmarshal(message, &resp)
		if err != nil {
			log.Println("unmarshal:", err)
			return err
		}
		ui.UUID = resp.ID
		return nil
	})
	ui.AddHandler(api.GameStatusOK, func(message []byte) error {
		var resp api.EventResponse
		err = json.Unmarshal(message, &resp)
		if err != nil {
			log.Println("unmarshal:", err)
			return err
		}
		if ui.Status == StatusInit {
			ui.Status = StatusStart
		}
		ui.Update(resp.Body.Players)
		return nil
	})
	ui.AddHandler(api.GameStatusError, func(message []byte) error {
		var resp api.EventResponse
		err = json.Unmarshal(message, &resp)
		if err != nil {
			log.Println("unmarshal:", err)
			return err
		}

		log.Printf("return from ConnectWebsocket read handler: %d", api.GameStatusError)
		ui.Status = StatusDrop
		for _, p := range resp.Body.Players {
			if p.ID == ui.UUID {
				ui.Score = p.Size
				break
			}
		}
		return fmt.Errorf("error")
	})
	ui.AddHandler(api.GameStatusWaiting, func(message []byte) error {
		log.Printf("Receive waiting event")
		return nil
	})

	mng.AddScene(SceneTypeMenu, NewMenuScene(ui))
	mng.AddScene(SceneTypeMatchmaking, NewMatchmakingScene(ui))
	mng.AddScene(SceneTypeIngame, NewIngameScene(ui))
	mng.SetInitialScene(SceneTypeMenu)

	mng.Execute()
	ui.Finish()
	log.Printf("========== GAME FINISH ==========")
}
