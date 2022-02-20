package main

type Client interface {
	ID() int
	Update(x, y, size, dir int, state string, board [][]int)
	Finish()
	Run(chan<- Event)
	Quit()
}

type RandomClient struct {
	id        int
	width     int
	height    int
	x         int
	y         int
	dir       int
	forceTurn chan int
	done      chan int
}

func (c *RandomClient) ID() int {
	return c.id
}

func (c *RandomClient) Update(x, y, size, dir int, state string, board [][]int) {
	c.x = x
	c.y = y
	c.dir = dir

	nextX, nextY := c.getNextCell()
	if nextX < 0 || nextX == c.width || nextY < 0 || nextY == c.height {
		c.forceTurn <- 1
	}
}

func (c *RandomClient) Finish() {
	c.done <- 1
}
func (c *RandomClient) Quit() {
	c.done <- 1
}
func (c *RandomClient) Run(event chan<- Event) {
	for {
		select {
		case <-c.forceTurn:
			var dir int
			switch c.dir {
			case MoveLeft:
				dir = MoveUp
			case MoveRight:
				dir = MoveDown
			case MoveUp:
				dir = MoveRight
			case MoveDown:
				dir = MoveLeft
			}
			event <- Event{
				ID:        c.ID(),
				Type:      "move",
				Direction: dir,
			}
		case <-c.done:
			return
		}
	}
}

func (c *RandomClient) getNextCell() (int, int) {
	switch c.dir {
	case MoveLeft:
		return c.x - 1, c.y
	case MoveRight:
		return c.x + 1, c.y
	case MoveUp:
		return c.x, c.y - 1
	case MoveDown:
		return c.x, c.y + 1
	}
	return 0, 0
}

func NewRandomClient(id, w, h int) (*RandomClient, error) {
	stream := make(chan int)
	doneStream := make(chan int)
	client := &RandomClient{
		id:        id,
		width:     w,
		height:    h,
		forceTurn: stream,
		done:      doneStream,
	}
	return client, nil
}

type GameClient struct {
	id    int
	event chan Event
}

func NewGameClient(id, w, h int) (*GameClient, error) {
	event := make(chan Event)
	client := &GameClient{
		id:    id,
		event: event,
	}
	return client, nil
}

func (c *GameClient) Finish() {
}

func (c *GameClient) Controller() chan<- Event {
	return c.event
}

type localClient struct{}

func (c *localClient) ID() int {
	return 0
}

func (c *localClient) Update(x, y, size, dir int, state string, board [][]int) {}
func (c *localClient) Finish()                                                 {}
func (c *localClient) Quit()                                                   {}
func (c *localClient) Run(event chan<- Event)                                  {}
