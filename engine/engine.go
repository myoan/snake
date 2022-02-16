package engine

type Engine struct {
	SceneManager *SceneManager
	Input        *Input
	event        chan ControlEvent
}

func NewGameEngine(fps int) *Engine {
	mng := NewSceneManager(fps)
	event := make(chan ControlEvent)
	interval := 1000 / fps
	input := NewInput(event, interval)

	return &Engine{
		SceneManager: mng,
		event:        event,
		Input:        input,
	}
}

func (engine *Engine) GetEventStream() chan<- ControlEvent {
	return engine.event
}
