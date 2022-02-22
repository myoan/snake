package api

type EventRequest struct {
	Eventtype int `json:"eventtype"`
	ID        int `json:"id"`
}

type PlayerResponse struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Size      int `json:"size"`
	Direction int `json:"direction"`
}

type EventResponse struct {
	Status  int              `json:"status"`
	Board   []int            `json:"board"`
	Width   int              `json:"width"`
	Height  int              `json:"height"`
	Players []PlayerResponse `json:"players"`
}
