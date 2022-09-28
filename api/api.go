package api

const (
	MoveLeft = iota
	MoveRight
	MoveUp
	MoveDown
)

type Message struct {
	UUID string `json:"uuid"`
	Path string `json:"path"`
	Body []byte `json:"body"`
}

type EventRequest struct {
	UUID      string `json:"uuid"`
	Eventtype int    `json:"eventtype"`
	Key       int    `json:"key"`
}

type PlayerResponse struct {
	ID        string `json:"id"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Size      int    `json:"size"`
	Direction int    `json:"direction"`
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

type InitResponse struct {
	Status int    `json:"status"`
	ID     string `json:"id"`
}

const (
	GameStatusInit = iota
	GameStatusOK
	GameStatusError
	GameStatusWaiting
)
