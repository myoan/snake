package engine

import "testing"

type DummyScene struct{}

func (s *DummyScene) Start() {}
func (s *DummyScene) Update() (SceneType, error) {
	return 0, nil
}
func (s *DummyScene) Finish() {}

func NewDummyScene() *DummyScene {
	return &DummyScene{}
}

func TestSceneManager_SetInitialScene(t *testing.T) {
	mng := NewSceneManager()
	mng.AddScene(SceneType(1), NewDummyScene())
	err := mng.SetInitialScene(SceneType(2))

	if err == nil {
		t.Errorf("SetInitialScene should return error if scene not found")
	}
}
