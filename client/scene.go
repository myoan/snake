package main

import (
	"errors"
	"fmt"

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
	if scene.UI.Score < 0 {
		scene.UI.DrawMenu([]string{
			fmt.Sprintf("ID: %s", scene.UI.UUID),
			"Score: -",
			"Press Space / Enter to Start",
			"Press Esc to Quit",
		})
	} else {
		scene.UI.DrawMenu([]string{
			fmt.Sprintf("ID: %s", scene.UI.UUID),
			fmt.Sprintf("Score: %d", scene.UI.Score),
			"Press Space / Enter to Start",
			"Press Esc to Quit",
		})
	}
}

func (scene *MenuScene) Update() (engine.SceneType, error) {
	if scene.Input.KeySpace {
		return SceneTypeMatchmaking, nil
	}
	if scene.Input.KeyEsc {
		return SceneTypeNone, ErrIngameQuited
	}

	if scene.UI.Score < 0 {
		scene.UI.DrawMenu([]string{
			fmt.Sprintf("ID: %s", scene.UI.UUID),
			"Score: -",
			"Press Space / Enter to Start",
			"Press Esc to Quit",
		})
	} else {
		scene.UI.DrawMenu([]string{
			fmt.Sprintf("ID: %s", scene.UI.UUID),
			fmt.Sprintf("Score: %d", scene.UI.Score),
			"Press Space / Enter to Start",
			"Press Esc to Quit",
		})
	}
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
	if scene.Input.KeyEsc {
		return SceneTypeNone, ErrIngameQuited
	}
	if scene.UI.Status == StatusStart {
		return SceneTypeIngame, nil
	}
	if scene.UI.Status == StatusDrop {
		return SceneTypeMenu, nil
	}
	scene.UI.DrawMenu([]string{
		"waiting for matchmaking...",
	})
	return SceneTypeMatchmaking, nil
}

func (scene *MatchmakingScene) Finish() {
	if scene.UI.Status == StatusDrop {
		scene.UI.CloseWebSocket()
	}
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
	if scene.UI.Status == StatusDrop {
		return SceneTypeMenu, nil
	}
	return SceneTypeIngame, nil
}

func (scene *IngameScene) Finish() {
	scene.UI.CloseWebSocket()
}
