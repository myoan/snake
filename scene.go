package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var (
	ErrIngameHitWall = errors.New("ingame hit wall")
	ErrIngameQuited  = errors.New("ingame quited")
)

type SceneType int

const (
	SceneTypeNone SceneType = iota
	SceneTypeMenu
	SceneTypeIngame
)

func NewSceneManager() *SceneManager {
	scenes := make(map[SceneType]Scene)
	return &SceneManager{scenes: scenes}
}

// SceneManager manages every scenes
// You must set all scenes and transitions before start
type SceneManager struct {
	scenes           map[SceneType]Scene
	currentSceneType SceneType
	currentScene     Scene
}

// Execute executes the state machine
func (mng *SceneManager) Execute() error {
	t := time.NewTicker(100 * time.Millisecond)
	mng.currentScene.Start()
	for range t.C {
		stype, err := mng.currentScene.Update()
		if err != nil {
			return err
		}
		if stype != mng.currentSceneType {
			logger.Printf("change scene %d -> %d", mng.currentSceneType, stype)
			mng.MoveTo(stype)
		}
	}
	mng.currentScene.Finish()
	return nil
}

func (mng *SceneManager) Stop() {
	fmt.Println("stop")
}

// SetFirstScene sets the first scene
func (mng *SceneManager) SetFirstScene(ty SceneType) {
	mng.currentSceneType = ty
	mng.currentScene = mng.scenes[ty]
}

// AddScene adds a scene to the manager
// If you set same type scene, it will be overwritten
func (mng *SceneManager) AddScene(ty SceneType, scene Scene) {
	mng.scenes[ty] = scene
}

// MoveTo changes current scene
func (mng *SceneManager) MoveTo(ty SceneType) error {
	scene := mng.scenes[ty]
	if scene == nil {
		return fmt.Errorf("scene %d not found", ty)
	}
	mng.currentScene.Finish()
	mng.currentScene = scene
	mng.currentSceneType = ty
	mng.currentScene.Start()
	return nil
}

type Scene interface {
	Start()
	Update() (SceneType, error)
	Finish()
}

type MenuScene struct {
	Input *Input
	UI    *UserInterface
}

func (scene *MenuScene) Start() {
	scene.UI.DrawMenu([]string{
		"Press Space / Enter to Start",
		"Press Esc to Quit",
	})
}

func (scene *MenuScene) Finish() {}

func (scene *MenuScene) Update() (SceneType, error) {

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

func NewMenuScene(ui *UserInterface, input *Input) *MenuScene {
	return &MenuScene{
		UI:    ui,
		Input: input,
	}
}

type IngameSceneStartArgs struct {
	width  int
	height int
}
type IngameScene struct {
	Input *Input
	UI    *UserInterface
}

func NewIngameScene(ui *UserInterface, input *Input) *IngameScene {
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

func (scene *IngameScene) Update() (SceneType, error) {
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
