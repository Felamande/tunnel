package main

import "fmt"

const (
	SwitchIP byte = iota + 1
	SwitchData
	SwitchIPOK
	ErrUnknowFlag
)

type Error struct {
	Func   string
	Action string
	Msg    interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s when %s: %v", e.Func, e.Action, e.Msg)
}
