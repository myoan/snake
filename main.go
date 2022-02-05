package main

import (
	"log"
	"os"
	"time"
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
	client, _ := NewGameClient(1, 40, 30)
	stateMachine := NewGameStateMachine(40, 30)
	stateMachine.AddGameClient(client)

	for range t.C {
		logger.Printf("tick")
		switch stateMachine.gs.State() {
		case GameInit:
			stateMachine.InitUpdate()
		case GameStart:
			logger.Printf("Stop tick")
			t.Stop()
			err := stateMachine.StartUpdate()
			if err != nil {
				os.Exit(0)
			}
			logger.Printf("Reset tick")
			t.Reset(d)
		case GameFinish:
			logger.Printf("execute finishUpdate")
			stateMachine.FinishUpdate()
		}
	}

	client.Finish()
	logger.Printf("========== GAME FINISH ==========")
}
