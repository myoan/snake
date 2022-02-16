package main

import (
	"log"
	"os"
)

const (
	Width    = 80
	Height   = 30
	interval = 100
)

func main() {
	mng := NewSceneManager()
	defer mng.Stop()

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
	// client, _ := NewGameClient(1, Width, Height)

	event := make(chan ControlEvent)
	input := NewInput(event)
	ui := NewUserInterface(event)
	mng.AddScene(SceneTypeMenu, NewMenuScene(ui, input))
	mng.AddScene(SceneTypeIngame, NewIngameScene(ui, input))
	mng.SetFirstScene(SceneTypeMenu)

	mng.Execute()
	ui.Finish()
	logger.Printf("========== GAME FINISH ==========")
}
