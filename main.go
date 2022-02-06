package main

import (
	"log"
	"os"
	"time"
)

const (
	Width  = 40
	Height = 30
)

func main() {
	f, err := os.OpenFile("game.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}
	logger = log.New(f, "", log.LstdFlags)
	d := time.Duration(100) * time.Millisecond
	t := time.NewTicker(d)
	defer func() {
		f.Sync()
		f.Close()
		t.Stop()
	}()

	logger.Printf("========== GAME START ==========")
	client, _ := NewGameClient(1, Width, Height)
	stateMachine := NewGameStateMachine(Width, Height)
	stateMachine.AddGameClient(client)

	for range t.C {
		switch stateMachine.gs.State() {
		case GameInit:
			logger.Printf("--- GameInit")
			stateMachine.InitUpdate()
		case GameStart:
			// t.Stop()
			logger.Printf("--- GameStart")
			err := stateMachine.StartUpdate(t)
			if err != nil {
				client.Finish()
				os.Exit(0)
			}
			// t.Reset(d)
		case GameFinish:
			logger.Printf("--- GameFinish")
			stateMachine.FinishUpdate()
		}
	}

	client.Finish()
	logger.Printf("========== GAME FINISH ==========")
}
