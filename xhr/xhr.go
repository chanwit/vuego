package xhr

import "syscall/js"

type XMLHttpRequest struct {
	js.Value
}

func New() *XMLHttpRequest {
	return &XMLHttpRequest{
		Value: js.Global().Get("XMLHttpRequest").New(),
	}
}

func (x *XMLHttpRequest) Open(method string, path string) {
	x.Value.Call("open", method, path)
}

func (x *XMLHttpRequest) OnLoad(fn func([]js.Value)) {
	cb := js.NewCallback(fn)
	x.Value.Set("onload", cb)
}

func (x *XMLHttpRequest) Send() {
	x.Value.Call("send")
}
func (x *XMLHttpRequest) GetResponseText() string {
	return x.Get("responseText").String()
}
