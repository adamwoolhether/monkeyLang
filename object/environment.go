package object

// Environment keeps track of values by associating
// them to an Object with a name. The outer field
// contains a reference to another object.Environment,
// which mirrors how variable scopes are perceived.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnclosedEnvironment allows extending an environment by
// creating a new Environment with a pointer to the enclosing
// environment that it extends. This allows preserving
// previous bindings of a function call.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	
	return env
}

// NewEnvironment returns a new *Environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
