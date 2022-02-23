package main

import (
	"errors"

	"github.com/myoan/snake/engine"
)

var (
	ErrIngameHitWall = errors.New("ingame hit wall")
	ErrIngameQuited  = errors.New("ingame quited")
)

type MenuScene struct {
	Input *engine.Input
	UI    *UserInterface
}

func NewMenuScene(input *engine.Input, ui *UserInterface) *MenuScene {
	return &MenuScene{
		Input: input,
		UI:    ui,
	}
}

func (scene *MenuScene) Start() {
	scene.UI.DrawMenu([]string{
		"Press Space / Enter to Start",
		"Press Esc to Quit",
	})
}

func (scene *MenuScene) Update() (engine.SceneType, error) {
	if scene.Input.KeySpace {
		return SceneTypeMatchmaking, nil
	}
	if scene.Input.KeyEsc {
		return SceneTypeNone, ErrIngameQuited
	}

	scene.UI.DrawMenu([]string{
		"Press Space / Enter to Start",
		"Press Esc to Quit",
	})
	return SceneTypeMenu, nil
}

func (scene *MenuScene) Finish() {}

type MatchmakingScene struct {
	Input *engine.Input
	UI    *UserInterface
}

func NewMatchmakingScene(input *engine.Input, ui *UserInterface) *MatchmakingScene {
	return &MatchmakingScene{
		UI:    ui,
		Input: input,
	}
}

func (scene *MatchmakingScene) Start() {
	go scene.UI.ConnectWebSocket()
	scene.UI.DrawMenu([]string{
		"waiting for matchmaking...",
	})
}

func (scene *MatchmakingScene) Update() (engine.SceneType, error) {
	return SceneTypeIngame, nil
}

func (scene *MatchmakingScene) Finish() {
	logger.Printf("move to ingame")
}

type IngameScene struct {
	Input *engine.Input
	UI    *UserInterface
}

func NewIngameScene(input *engine.Input, ui *UserInterface) *IngameScene {
	return &IngameScene{
		UI:    ui,
		Input: input,
	}
}

func (scene *IngameScene) Start() {}

func (scene *IngameScene) Update() (engine.SceneType, error) {
	if scene.UI.Status != 0 {
		return SceneTypeMenu, nil
	}
	return SceneTypeIngame, nil
}

func (scene *IngameScene) Finish() {
	scene.UI.CloseWebSocket()
}
