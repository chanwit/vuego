package vuego

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"reflect"
	"syscall/js"
)

type Vue struct {
	js.Value
	El   string
	Data interface{}

	Methods *Methods
	Filters *Filters
	Watch   *Watch

	Created func(interface{}, []js.Value)
	Mounted func(interface{}, []js.Value)
}

var ctor = js.Global().Get("Object")
var arr = js.Global().Get("Array")

func ValueOf(rv reflect.Value) js.Value {
	fmt.Printf(">> ValueOf: DEBUG calling ValueOf type=%s\n", rv.Type().Name())
	switch rv.Kind() {
	case reflect.Slice:
		fmt.Printf(">>>> case slice (%v)\n", rv.Interface())
		a := arr.New(rv.Len())
		for i := 0; i < rv.Len(); i++ {
			fmt.Printf(">>>>>> case slice elem (%v)\n", rv.Index(i).Interface())
			a.SetIndex(i, ToJSON(rv.Index(i).Interface()))
		}
		return a
	default:
		fmt.Printf(">>>> case default (%v)\n", rv.Interface())
		return js.ValueOf(rv.Interface())
	}
}

func ToJSON(this interface{}) js.Value {
	// this.Value
	thisValue := ctor.New()

	v := reflect.Indirect(reflect.ValueOf(this))
	r := v.Type()

	if r.Name() == "string" {
		return js.ValueOf(v.Interface())
	}

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

func New(v *Vue) *Vue {
	arg := ctor.New()
	// arg.Set("_VUEGO_ID_", id)
	arg.Set("el", v.El)
	arg.Set("data", ToJSON(v.Data))
	arg.Set("watch", js.ValueOf(*v.Watch))
	arg.Set("created", js.NewCallback(func(args []js.Value) {
		v.Created(v.Data, args)
	}))

	filters := ctor.New()
	for name, fn := range *v.Filters {
		id := "VUEGO_FILTER_" + ksuid.New().String()
		output := id + "_OUTPUT"
		// VUEGO_FILTER_1BThKORJGFF1Ayjinc54GDdFlEb_OUTPUT
		cb := js.NewCallback(func(args []js.Value) {
			js.Global().Set(output, js.ValueOf(fn(args[0])))
		})
		js.Global().Set(id, cb)
		wrapper := js.Global().Call("eval", fmt.Sprintf("(function(v){ %s(v); return global.%s; })", id, output));
		filters.Set(name, wrapper)
	}
	arg.Set("filters", filters)

	methods := ctor.New()
	for name, fn := range *v.Methods {
		cb := js.NewCallback(func(args []js.Value) {
			fn(v.Data)
		});
		methods.Set(name, js.ValueOf(cb))
	}
	arg.Set("methods", methods)

	v.Value = js.Global().Get("Vue").New(arg)
	return v
}

func (v *Vue) Watch0(prop string, callback func(newValue, oldValue js.Value)) {
	cb := js.NewCallback(func(args []js.Value) {
		callback(args[0], args[1])
	})
	v.Value.Call("$watch", js.ValueOf(prop), cb)
}
