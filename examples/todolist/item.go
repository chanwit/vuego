package main

import "syscall/js"

//go:generate go run github.com/chanwit/vuego/gen
type Item struct {
	js.Value
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
}
