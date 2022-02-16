package main

import (
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

func main() {
	f, err := os.OpenFile("game.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}
	logger = log.New(f, "", log.LstdFlags)
	defer func() {
		f.Sync()
		f.Close()
	}()

	logger.Printf("========== GAME START ==========")

	ge := engine.NewGameEngine(10)
	mng := ge.SceneManager
	defer mng.Stop()

	event := ge.GetEventStream()
	input := ge.Input
	ui := NewUserInterface(event)
	mng.AddScene(SceneTypeMenu, NewMenuScene(input, ui))
	mng.AddScene(SceneTypeIngame, NewIngameScene(input, ui))
	mng.SetInitialScene(SceneTypeMenu)

	mng.Execute()
	ui.Finish()
	logger.Printf("========== GAME FINISH ==========")
}
