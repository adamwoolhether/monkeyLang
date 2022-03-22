package object

// Environment keeps track of values by associating
// them to an Object with a name.
type Environment struct {
	store map[string]Object
}

// NewEnvironment returns a new *Environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
