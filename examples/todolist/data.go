package main

import (
	"strings"
	"syscall/js"
)

//go:generate go run github.com/chanwit/vuego/gen
type Data struct {
	js.Value
	Items   []*Item `json:"items"`
	Title   string  `json:"title"`
	NewItem string  `json:"newItem"`
}

func (this *Data) AddItem() {
	text := strings.TrimSpace(this.GetNewItem())
	if len(text) > 0 {
		item := &Item{
			Text: text,
			Checked: false,
		}

		this.AddToItems(item)
		this.SetNewItem("")
	}
}
