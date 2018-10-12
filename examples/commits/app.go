package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	. "github.com/chanwit/vuego"
	"github.com/chanwit/vuego/xhr"
)

const (
	apiURL = "https://api.github.com/repos/vuejs/vue/commits?per_page=3&sha="
)

type Commits struct {
	js.Value
	Branches      []string `json:"branches"`
	CurrentBranch string   `json:"currentBranch"`
	Commits       []Commit `json:"commits"`
}

type Commit = map[string]interface{}

func (this *Commits) FetchData() {
	req := xhr.New()
	req.Open("GET", apiURL + this.CurrentBranch)
	req.OnLoad(func(args []js.Value) {
		this.Commits = []Commit{}
		json.Unmarshal([]byte(req.GetResponseText()), &this.Commits)
		fmt.Println(this.Commits[0]["html_url"])
		arrayConstructor := js.Global().Get("Array")
		a := arrayConstructor.New(len(this.Commits))
		for i, s := range this.Commits {
			a.SetIndex(i, js.ValueOf(s))
		}
		this.Value.Set("commits", a)
	})
	req.Send()
}

func main() {

	New(&Vue{
		El: "#demo",
		Data: &Commits{
			Branches:      []string{"master", "dev"},
			CurrentBranch: "master",
			Commits:       nil,
		},
		Created: func(this Any, args []js.Value) {
			this.(*Commits).FetchData()
		},
		Watch: &Watch{
			"currentBranch": "fetchData",
		},
		Filters: &Filters{
			"truncate": func(v Any) Any {
				s := v.(js.Value).String()
				return strings.TrimSuffix(s, "\n")
			},
			"formatDate": func(v Any) Any {
				return strings.Replace(v.(js.Value).String(), "Z", " ", -1)
			},
		},
		Methods: &Methods{
			"fetchData": func(this Any) {
				this.(*Commits).FetchData()
			},
		},
	})

	select { }

}
