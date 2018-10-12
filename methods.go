package vuego

type Any = interface{}

type Methods map[string]func(Any)

type Filters map[string]func(Any)Any

type Watch = map[string]Any