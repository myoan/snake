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

func (scene *MenuScene) Start() {
	scene.UI.DrawMenu([]string{
		"Press Space / Enter to Start",
		"Press Esc to Quit",
	})
}

func (scene *MenuScene) Finish() {}

func (scene *MenuScene) Update() (engine.SceneType, error) {
	if scene.Input.KeySpace {
		logger.Printf("push Space")
		return SceneTypeIngame, nil
	}
	if scene.Input.KeyEsc {
		logger.Printf("push Esc")
		return SceneTypeNone, ErrIngameQuited
	}

	scene.UI.DrawMenu([]string{
		"Press Space / Enter to Start",
		"Press Esc to Quit",
	})
	return SceneTypeMenu, nil
}

func NewMenuScene(input *engine.Input, ui *UserInterface) *MenuScene {
	return &MenuScene{
		Input: input,
		UI:    ui,
	}
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

func (scene *IngameScene) Start() {
	go scene.UI.ConnectWebSocket()
}

func (scene *IngameScene) Update() (engine.SceneType, error) {
	if scene.UI.Status != 0 {
		return SceneTypeMenu, nil
	}
	return SceneTypeIngame, nil
}

func (scene *IngameScene) Finish() {
	scene.UI.CloseWebSocket()
}
