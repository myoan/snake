package main

type ServerUserInterface struct {
}

func NewUserInterface() *ServerUserInterface {
	return &ServerUserInterface{}
}

func (ui *ServerUserInterface) Finish() {}

func NewIngameScene() {}
