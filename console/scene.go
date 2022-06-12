package main

import (
	"errors"
	"log"

	"github.com/myoan/snake/engine"
)

var (
	ErrIngameHitWall = errors.New("ingame hit wall")
	ErrIngameQuited  = errors.New("ingame quited")
)

type MenuScene struct {
	UI *UserInterface
}

func NewMenuScene(ui *UserInterface) *MenuScene {
	return &MenuScene{
		UI: ui,
	}
}

func (scene *MenuScene) Start() {
	log.Printf("ID: %s, Score: %d", scene.UI.UUID, scene.UI.Score)
}

func (scene *MenuScene) Update() (engine.SceneType, error) {
	return SceneTypeMatchmaking, nil
}

func (scene *MenuScene) Finish() {}

type MatchmakingScene struct {
	UI *UserInterface
}

func NewMatchmakingScene(ui *UserInterface) *MatchmakingScene {
	return &MatchmakingScene{
		UI: ui,
	}
}

func (scene *MatchmakingScene) Start() {
	go scene.UI.ConnectWebSocket()
	log.Println("waiting for matchmaking...")
}

func (scene *MatchmakingScene) Update() (engine.SceneType, error) {
	if scene.UI.Status == StatusStart {
		return SceneTypeIngame, nil
	}
	if scene.UI.Status == StatusDrop {
		return SceneTypeMenu, nil
	}
	return SceneTypeMatchmaking, nil
}

func (scene *MatchmakingScene) Finish() {
	if scene.UI.Status == StatusDrop {
		scene.UI.CloseWebSocket()
	}
	log.Printf("move to ingame")
}

type IngameScene struct {
	UI *UserInterface
}

func NewIngameScene(ui *UserInterface) *IngameScene {
	return &IngameScene{
		UI: ui,
	}
}

func (scene *IngameScene) Start() {}

func (scene *IngameScene) Update() (engine.SceneType, error) {
	if scene.UI.Status == StatusDrop {
		return SceneTypeMenu, nil
	}
	return SceneTypeIngame, nil
}

func (scene *IngameScene) Finish() {
	scene.UI.CloseWebSocket()
}
