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
	init := true
	var scene *IngameScene

	i := 0
	event := make(chan ControlEvent)
	for range t.C {
		switch stateMachine.gs.State() {
		case GameInit:
			// TODO: ここはいわゆるゲームエンジンとしての設計で良い。つまり、inputは一つだけ

			logger.Printf("--- GameInit")
			// TODO: Add CPU Player
			ingame := stateMachine.gc.NewIngameClient(stateMachine.gs.Game.FetchEvent())
			stateMachine.gs.ResetClient()
			stateMachine.gs.AddClient(ingame)
			stateMachine.gs.Game.ResetPlayers()
			stateMachine.gs.Game.AddPlayer(ingame)
			stateMachine.sm.Run(stateMachine.gs.Start, GameArgument{clients: stateMachine.gs.Clients})
		case GameStart:
			// TODO: ここは複数のクライアントが集約するイメージ
			logger.Printf("--- GameStart")
			if init {
				scene = NewIngameScene(event)
				ingame := stateMachine.gc.NewIngameClient(stateMachine.gs.Game.FetchEvent())
				scene.Start(Width, Height, ingame)
				init = false
			}
			scene.Update()

			/*
				err := stateMachine.gs.Game.Start(t)
				if err != nil {
					logger.Printf("[ERROR] %v", err)
					client.Finish()
					os.Exit(0)
				}
			*/
			stateMachine.sm.Run(stateMachine.gs.Finish, GameArgument{clients: stateMachine.gs.Clients})
			if i > 50 {
				client.Finish()
				os.Exit(0)
			}
			i++
		case GameFinish:
			logger.Printf("--- GameFinish")

			err := stateMachine.sm.Run(stateMachine.gs.Restart, GameArgument{clients: stateMachine.gs.Clients})
			if err != nil {
				client.Finish()
				os.Exit(0)
			}
		}
	}

	client.Finish()
	logger.Printf("========== GAME FINISH ==========")
}
