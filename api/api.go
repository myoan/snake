package api

const (
	MoveLeft = iota
	MoveRight
	MoveUp
	MoveDown
)

type EventRequest struct {
	ID        int `json:"id"`
	Eventtype int `json:"eventtype"`
	Key       int `json:"key"`
}

type PlayerResponse struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Size      int `json:"size"`
	Direction int `json:"direction"`
}

type ResponseBody struct {
	Board   []int            `json:"board"`
	Width   int              `json:"width"`
	Height  int              `json:"height"`
	Players []PlayerResponse `json:"players"`
}

type EventResponse struct {
	Status int          `json:"status"`
	Body   ResponseBody `json:"body"`
}

type RestApiGetIDResponse struct {
	ID int `json:"id"`
}

const (
	GameStatusOK = iota
	GameStatusError
	GameStatusWaiting
)
