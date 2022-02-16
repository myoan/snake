package engine

import (
	"fmt"
	"time"
)

type SceneType int

func NewSceneManager(fps int) *SceneManager {
	scenes := make(map[SceneType]Scene)
	return &SceneManager{
		fps:    fps,
		scenes: scenes,
	}
}

// SceneManager manages every scenes
// You must set all scenes and transitions before start
type SceneManager struct {
	fps              int
	scenes           map[SceneType]Scene
	initSceneType    SceneType
	initScene        Scene
	currentSceneType SceneType
	currentScene     Scene
}

// Execute executes the state machine
func (mng *SceneManager) Execute() error {
	interval := 1000 / mng.fps
	t := time.NewTicker(time.Duration(interval) * time.Millisecond)
	mng.currentSceneType = mng.initSceneType
	mng.currentScene = mng.initScene

	mng.currentScene.Start()
	for range t.C {
		stype, err := mng.currentScene.Update()
		if err != nil {
			return err
		}
		if stype != mng.currentSceneType {
			mng.MoveTo(stype)
		}
	}
	mng.currentScene.Finish()
	return nil
}

func (mng *SceneManager) Stop() {}

// SetInitialScene sets the first scene
// If scene not found, it will return error
func (mng *SceneManager) SetInitialScene(ty SceneType) error {
	if mng.scenes[ty] == nil {
		return fmt.Errorf("scene %d not found", ty)
	}
	mng.initSceneType = ty
	mng.initScene = mng.scenes[ty]
	return nil
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
