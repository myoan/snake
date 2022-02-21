package main

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gorilla/websocket"
	"github.com/myoan/snake/engine"
)

/*
UserInterface represents the user interface, screen and controller.
This struct must be initialized at the beginning of the program and must live until the end.
You shouldn't re-create this struct.
*/

type UserInterface struct {
	Status   int
	screen   tcell.Screen
	webEvent chan engine.ControlEvent
	webDone  chan struct{}
	conn     *websocket.Conn
}

// NewUserInterface creates a new UserInterface.
// You must call this method before using the UserInterface.
// UserInterface is listening user controlle events and sends them to the event channel.
// webEvent is a channel for sending to web server.
func NewUserInterface(event chan<- engine.ControlEvent, webEvent chan engine.ControlEvent) *UserInterface {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorPurple)
	s, err := tcell.NewScreen()
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		logger.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	done := make(chan struct{})

	ui := &UserInterface{
		Status:   0,
		screen:   s,
		webEvent: webEvent,
		webDone:  done,
	}
	go ui.runController(event, webEvent)

	return ui
}

/*
Finish is called when the entire game is over.
*/
func (ui *UserInterface) Finish() {
	ui.screen.Fini()
}

// Draw shows the board on the screen.
// You should call this method periodically.
func (ui *UserInterface) Draw(b *Board) {
	snakeStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
	appleStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	screenWidth, screenHeight := ui.screen.Size()

	originX := (screenWidth - b.width) / 2
	originY := (screenHeight - b.height) / 2

	ui.screen.Clear()
	for y, row := range b.board {
		for x, cell := range row {
			if cell == 0 {
				// Empty
				ui.screen.SetContent(originX+x+2, originY+y+1, ' ', nil, defStyle)
			} else if cell > 0 {
				// Snake
				ui.screen.SetContent(originX+x+2, originY+y+1, ' ', nil, snakeStyle)
			} else {
				// Applg
				ui.screen.SetContent(originX+x+2, originY+y+1, ' ', nil, appleStyle)
			}
		}
	}

	// Draw borders
	for col := 0; col <= b.width+1; col++ {
		ui.screen.SetContent(originX+col, originY, tcell.RuneHLine, nil, defStyle)
		ui.screen.SetContent(originX+col, originY+b.height+1, tcell.RuneHLine, nil, defStyle)
	}
	for row := 0; row <= b.height+1; row++ {
		ui.screen.SetContent(originX, originY+row, tcell.RuneVLine, nil, defStyle)
		ui.screen.SetContent(originX+b.width+2, originY+row, tcell.RuneVLine, nil, defStyle)
	}

	// draw corners
	ui.screen.SetContent(originX, originY, tcell.RuneULCorner, nil, defStyle)
	ui.screen.SetContent(originX, originY+b.height+1, tcell.RuneLLCorner, nil, defStyle)
	ui.screen.SetContent(originX+b.width+2, originY, tcell.RuneURCorner, nil, defStyle)
	ui.screen.SetContent(originX+b.width+2, originY+b.height+1, tcell.RuneLRCorner, nil, defStyle)

	ui.screen.Show()
}

// DrawMenu displays input list.
// Each line is centered.
func (ui *UserInterface) DrawMenu(strs []string) {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	width, height := ui.screen.Size()
	row := len(strs)*2 - 1
	ui.screen.Clear()
	starth := (height - row) / 2

	for i, str := range strs {
		startw := (width - len(str)) / 2
		h := starth + i*2
		for j := 0; j < len(str); j++ {
			ui.screen.SetContent(startw+j, h, rune(str[j]), nil, style)
		}
	}

	ui.screen.Show()
}

func (ui *UserInterface) runController(event chan<- engine.ControlEvent, webEvent chan<- engine.ControlEvent) {
	logger.Printf("runController")
	for {
		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			ui.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				event <- engine.ControlEvent{Eventtype: 0, ID: 1}
				webEvent <- engine.ControlEvent{Eventtype: 0, ID: 1}
			} else if ev.Rune() == 'a' || ev.Key() == tcell.KeyLeft {
				event <- engine.ControlEvent{Eventtype: 0, ID: 2}
				webEvent <- engine.ControlEvent{Eventtype: 0, ID: MoveLeft}
			} else if ev.Rune() == 'd' || ev.Key() == tcell.KeyRight {
				event <- engine.ControlEvent{Eventtype: 0, ID: 3}
				webEvent <- engine.ControlEvent{Eventtype: 0, ID: MoveRight}
			} else if ev.Rune() == 'w' || ev.Key() == tcell.KeyUp {
				event <- engine.ControlEvent{Eventtype: 0, ID: 4}
				webEvent <- engine.ControlEvent{Eventtype: 0, ID: MoveUp}
			} else if ev.Rune() == 's' || ev.Key() == tcell.KeyDown {
				event <- engine.ControlEvent{Eventtype: 0, ID: 5}
				webEvent <- engine.ControlEvent{Eventtype: 0, ID: MoveDown}
			} else if ev.Rune() == ' ' || ev.Key() == tcell.KeyEnter {
				event <- engine.ControlEvent{Eventtype: 0, ID: 6}
				// webEvent <- engine.ControlEvent{Eventtype: 0, ID: 6}
			}
		}
	}
}

type EventResponse struct {
	Status  int              `json:"status"`
	Board   []int            `json:"board"`
	Width   int              `json:"width"`
	Height  int              `json:"height"`
	Players []PlayerResponse `json:"players"`
}

type PlayerResponse struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Size      int `json:"size"`
	Direction int `json:"direction"`
}

func (ui *UserInterface) ConnectWebSocket() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ingame"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Printf("dial: %v", err)
		return
	}
	ui.conn = c

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logger.Println("read:", err)
				return
			}
			var resp EventResponse
			err = json.Unmarshal(message, &resp)
			if err != nil {
				logger.Println("unmarshal:", err)
				return
			}

			if resp.Status != 0 {
				ui.Status = 1
				logger.Printf("return from ConnectWebsocket read handler: %d", resp.Status)
				return
			}

			rawBoard := make([][]int, resp.Height)
			for i := 0; i < resp.Height; i++ {
				rawBoard[i] = make([]int, resp.Width)
				for j := 0; j < resp.Width; j++ {
					rawBoard[i][j] = resp.Board[i*resp.Width+j]
				}
			}
			board := &Board{
				board:  rawBoard,
				width:  resp.Width,
				height: resp.Height,
			}
			ui.Draw(board)
		}
	}()

	for {
		select {
		case ctrl := <-ui.webEvent:
			event := &EventRequest{
				Eventtype: ctrl.Eventtype,
				ID:        ctrl.ID,
			}
			bytes, _ := json.Marshal(&event)
			err := c.WriteMessage(websocket.TextMessage, bytes)
			if err != nil {
				logger.Println("write:", err)
				return
			}
		case <-ui.webDone:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Println("write close:", err)
				return
			}
			select {
			case <-ui.webDone:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (ui *UserInterface) CloseWebSocket() {
	ui.webDone <- struct{}{}
	time.Sleep(time.Duration(300) * time.Millisecond)
	ui.Status = 0
	ui.conn.Close()
}
