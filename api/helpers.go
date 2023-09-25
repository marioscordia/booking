package api

import (
	"fmt"
	"log"
	"runtime/debug"
)

type Error struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

func NewError(code int, msg string) Error {
	return Error{Code: code, Msg: msg}
}

func ErrorLog(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Output(2, trace)
}