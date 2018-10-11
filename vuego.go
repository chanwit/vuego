package vuego

import (
	"fmt"
	"reflect"
	"syscall/js"
)

type Vue struct {
	js.Value
	El      string
	Data    interface{}
	Methods *Methods

	Created func()
	Mounted func()
}

var ctor = js.Global().Get("Object")
var arr = js.Global().Get("Array")

func ValueOf(rv reflect.Value) js.Value {
	fmt.Printf(">> ValueOf: DEBUG calling ValueOf type=%s\n", rv.Type().Name())
	switch rv.Kind() {
	case reflect.Slice:
		fmt.Println(">>>> case slice")
		a := arr.New(rv.Len())
		for i := 0; i < rv.Len(); i++ {
			a.SetIndex(i, ToJSON(rv.Index(i).Interface()))
		}
		return a
	default:
		fmt.Println(">>>> case default")
		return js.ValueOf(rv.Interface())
	}
}

func ToJSON(this interface{}) js.Value {
	// this.Value
	thisValue := ctor.New()

	v := reflect.Indirect(reflect.ValueOf(this))
	r := v.Type()

	fmt.Printf(">>> DEBUG ToJSON type=%s\n", r.Name())

	// starting with field 1
	for i := 1; i < r.NumField(); i++ {
		// field to init
		// must be tagged with json:"name"
		name := r.Field(i).Tag.Get("json")
		if name != "" {
			thisValue.Set(name, ValueOf(v.Field(i)))
		}
	}

	v.Field(0).Set(reflect.ValueOf(thisValue))

	return thisValue
}

/*
func (this *M) SetMessage(s string) {
	this.Value.Set("message", js.ValueOf(s))
}
*/

func New(v *Vue) *Vue {
	arg := ctor.New()
	arg.Set("el", v.El)
	arg.Set("data", ToJSON(v.Data))

	methods := ctor.New()
	for name, fn := range *v.Methods {
		cb := js.NewCallback(func(args []js.Value) {
			fmt.Println(">>>>>>>>>>>>> add item")
			fn(v.Data)
		});
		methods.Set(name, js.ValueOf(cb))
	}
	arg.Set("methods", methods)

	v.Value = js.Global().Get("Vue").New(arg)
	return v
}

func (v *Vue) Watch(prop string, callback func(newValue, oldValue js.Value)) {
	cb := js.NewCallback(func(args []js.Value) {
		callback(args[0], args[1])
	})
	v.Value.Call("$watch", js.ValueOf(prop), cb)
}
