package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Start()
	Update() (SceneType, error)
	Finish()
	Draw(*ebiten.Image)
}

type SceneType string

type SceneManager struct {
	scenes           map[SceneType]Scene
	currentSceneType SceneType
	currentScene     Scene
}

func NewSceneManager() *SceneManager {
	scenes := make(map[SceneType]Scene)
	return &SceneManager{
		scenes: scenes,
	}
}

func (mng *SceneManager) AddScene(name SceneType, scene Scene) {
	mng.scenes[name] = scene
}

func (mng *SceneManager) SetInitialScene(name SceneType) {
	mng.currentScene = mng.scenes[name]
}

func (mng *SceneManager) Update() error {
	stype, err := mng.currentScene.Update()
	if err != nil {
		return err
	}

	if stype == mng.currentSceneType {
		return nil
	}
	return mng.moveTo(stype)
}

func (mng *SceneManager) Draw(screen *ebiten.Image) {
	mng.currentScene.Draw(screen)
}

// moveTo changes current scene
func (mng *SceneManager) moveTo(ty SceneType) error {
	fmt.Printf("%s -> %s\n", string(mng.currentSceneType), string(ty))
	scene := mng.scenes[ty]
	if scene == nil {
		return fmt.Errorf("scene %s not found", ty)
	}
	mng.currentScene.Finish()
	mng.currentScene = scene
	mng.currentSceneType = ty
	mng.currentScene.Start()
	return nil
}
