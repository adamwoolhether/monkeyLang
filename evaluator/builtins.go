package evaluator

import (
	"github.com/adamwoolhether/monkeyLang/object"
)

/*
len([1, 2, 3]); // => 3
first([1, 2, 3]); // => 1
last([1, 2, 3]); // => 3
rest([1, 2, 3]); // => [2, 3]
push([1, 2, 3], 4); // => [1, 2, 3, 4]
puts("Hello World!"); // prints "Hello World!"
*/

// builtins enables built-in funcs to be accessed.
var builtins = map[string]*object.Builtin{
	"len":   object.GetBuiltinByName("len"),
	"puts":  object.GetBuiltinByName("puts"),
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	"rest":  object.GetBuiltinByName("rest"),
	"push":  object.GetBuiltinByName("push"),
}
