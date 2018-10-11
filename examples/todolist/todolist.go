package main

import (
	. "github.com/chanwit/vuego"
)

func main() {

	New(&Vue{
		El: "#app",
		Data: &Data{
			Items: []*Item {
				{Text: "Bananas", Checked: true},
				{Text: "Apples", Checked: false},
			},
			Title: "Vuego Framework 0.1",
			NewItem: "",
		},
		Methods: &Methods{
			"addItem": func(this interface{}) {
				this.(*Data).AddItem()
			},
		},
	})

	select {}
}
