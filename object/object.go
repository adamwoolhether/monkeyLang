// Package object defines the object system for Monkey lang.
// All values encountered in evaluating Monkey source code
// are represented as an Object.
package object

type ObjectType string

// Object defines the contract for all values in Monkey
type Object interface {
	Type() ObjectType
	Inspect() string
}
