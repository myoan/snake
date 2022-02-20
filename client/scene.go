package main

import (
	"errors"
	"math/rand"

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
	board := NewBoard(Width, Height)
	board.GenerateApple()

	event := make(chan Event)

	localGame = &LocalGame{
		board:     board,
		event:     event,
		x:         rand.Intn(Width),
		y:         rand.Intn(Height),
		size:      3,
		direction: rand.Intn(4),
	}

	localGame.GenerateSnake()
}

func (scene *IngameScene) Update() (engine.SceneType, error) {
	if scene.Input.KeyA {
		logger.Printf("turn <-")
		localGame.changeDirection(MoveLeft)
	}
	if scene.Input.KeyD {
		logger.Printf("turn ->")
		localGame.changeDirection(MoveRight)
	}
	if scene.Input.KeyW {
		logger.Printf("turn ^")
		localGame.changeDirection(MoveUp)
	}
	if scene.Input.KeyS {
		logger.Printf("turn v")
		localGame.changeDirection(MoveDown)
	}

	err := localGame.MovePlayer()
	if err != nil {
		localGame.board.Reset()
		return SceneTypeMenu, nil
	}
	localGame.board.Update()
	scene.UI.Draw(localGame.board)
	return SceneTypeIngame, nil
}

func (scene *IngameScene) Finish() {}
