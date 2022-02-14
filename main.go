package main

import (
	"log"
	"os"
	"runtime"
	"time"
)

const (
	Width    = 80
	Height   = 30
	interval = 100
)

func main() {
	f, err := os.OpenFile("game.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}
	logger = log.New(f, "", log.LstdFlags)
	d := time.Duration(interval) * time.Millisecond
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
	var scene *IngameScene

	event := make(chan ControlEvent)
	var isQuit = false
	var isFinish = false
	scene = NewIngameScene(event)
Loop:
	for range t.C {
		switch stateMachine.gs.State() {
		case GameInit:
			// ここはいわゆるゲームエンジンとしての設計で良い。つまり、inputは一つだけ

			logger.Printf("--- GameInit(%d)", runtime.NumGoroutine())
			// TODO: Add CPU Player
			scene.Start(Width, Height)
			stateMachine.sm.Run(stateMachine.gs.Start, GameArgument{clients: stateMachine.gs.Clients, isFinish: false, isQuit: isQuit})
		case GameStart:
			// ここは複数のクライアントが集約するイメージ
			logger.Printf("--- GameStart(%d)", runtime.NumGoroutine())
			err := scene.Update()
			if err == ErrIngameHitWall {
				isFinish = true
			} else if err == ErrIngameQuited {
				isQuit = true
				break Loop
			}

			stateMachine.sm.Run(stateMachine.gs.Finish, GameArgument{clients: stateMachine.gs.Clients, isFinish: isFinish, isQuit: isQuit})
		case GameFinish:
			logger.Printf("--- GameFinish(%d)", runtime.NumGoroutine())
			scene.Finish()
			isFinish = false

			if isQuit {
				break Loop
			}

			err := stateMachine.sm.Run(stateMachine.gs.Restart, GameArgument{clients: stateMachine.gs.Clients, isFinish: isQuit})
			if err != nil {
				break Loop
			}
		}
	}
	scene.Finish()
	client.Finish()
	scene.UI.Finish()

	logger.Printf("========== GAME FINISH ==========")
}
