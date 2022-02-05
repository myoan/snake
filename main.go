package main

import (
	"log"
	"os"
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
	client, _ := NewGameClient(1, 40, 30)
	stateMachine := NewGameStateMachine(40, 30)
	stateMachine.AddGameClient(client)
	for i := 0; i < 5; i++ {
		stateMachine.InitUpdate()
		stateMachine.StartUpdate()
		stateMachine.FinishUpdate()
	}

	client.Finish()

	/*
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
	*/
	logger.Printf("========== GAME FINISH ==========")
}
