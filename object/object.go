// Package object defines the object system for Monkey lang.
// All values encountered in evaluating Monkey source code
// are represented as an Object.
package object

import (
	"bytes"
	"fmt"
	"strings"
	
	"github.com/adamwoolhether/monkeyLang/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
)

// Object defines the contract for all values in Monkey.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer holds the value of
// ast.IntegerLiteral objects.
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// Boolean holds the value of
// ast.Boolean objects
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Null represents the absence of any value.
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue represents a return value in Monkey.
// It's essentially just a wrapper around Object.
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error represents an internal error in Monkey. Errors for
// wrong operators, unsupported operations, and other user
// or internal errors that can arise during execution.
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Function represents a Function internally, holding the function
// Body, Parameters. It also has an Env field, beacuse monkey
// functions carry their own environment, which allows for closures.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n}")
	out.WriteString(f.Body.String())
	out.WriteString("\n")
	
	return out.String()
}

// String allows representation of strings in Monkey. This is
// simplified due to go's native support for string.
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// BuiltinFunction allows implementation of native functions in Monkey.
// The only restriction is that they need to accept zero or more
// object.Object as args and return an object.Object.
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// Array allows arrays to be used in Monkey. It uses
// go's slice behind the scenes.
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	
	return out.String()
}
