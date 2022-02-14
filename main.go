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

	StatusInit = iota
	StatusStart
	StatusFinish
	StatusQuit
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
	scene = NewIngameScene(event)
Loop:
	for range t.C {
		switch stateMachine.gs.State() {
		case GameInit:
			// ここはいわゆるゲームエンジンとしての設計で良い。つまり、inputは一つだけ

			logger.Printf("--- GameInit(%d)", runtime.NumGoroutine())
			// TODO: Add CPU Player
			args := IngameSceneStartArgs{
				width:  Width,
				height: Height,
			}
			scene.Start(args)
			stateMachine.sm.Run(stateMachine.gs.Start, GameArgument{scene: scene, status: StatusInit})
		case GameStart:
			// ここは複数のクライアントが集約するイメージ
			logger.Printf("--- GameStart(%d)", runtime.NumGoroutine())
			err, status := scene.Update(nil)
			if err == ErrIngameQuited {
				break Loop
			}

			stateMachine.sm.Run(stateMachine.gs.Finish, GameArgument{scene: scene, status: status})
		case GameFinish:
			logger.Printf("--- GameFinish(%d)", runtime.NumGoroutine())
			scene.Finish(nil)

			err := stateMachine.sm.Run(stateMachine.gs.Restart, GameArgument{scene: scene, status: StatusFinish})
			if err != nil {
				break Loop
			}
		}
	}
	client.Finish()
	scene.UI.Finish()

	logger.Printf("========== GAME FINISH ==========")
}
