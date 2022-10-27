// Package compiler evaluates given instructions and compiles Bytecode.
// It should walk the AST recursively, find *ast.IntegerLiterals, evaluate
// them and turn them into *object.Integers before adding them to the
// constants field and adding OpConstant instructions to the internal
// instructions slice.
package compiler

import (
	"fmt"

	"github.com/adamwoolhether/monkeyLang/ast"
	"github.com/adamwoolhether/monkeyLang/code"
	"github.com/adamwoolhether/monkeyLang/object"
)

// Compiler holds generated bytecode('instruction') a a pool of constants.
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// Compile determines how to handle given base on the node type.
func (c *Compiler) Compile(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		for _, s := range n.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(n.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop) // clean the stack.
	case *ast.InfixExpression:
		err := c.Compile(n.Left)
		if err != nil {
			return err
		}

		err = c.Compile(n.Right)
		if err != nil {
			return err
		}

		switch n.Operator {
		case "+":
			c.emit(code.OpAdd)
		default:
			return fmt.Errorf("unknown operator %s", n.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: n.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	}

	return nil
}

// Bytecode returns Bytecode from the compiler-generations instructions.
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// addConstant appends an object.Object to the end of a compiler's constants
// slice, returning its index as an identifier
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)

	return len(c.constants) - 1
}

// emit will generate an instruction and adding them to a collection in memeory.
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	return pos
}
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstructions := len(c.instructions)
	c.instructions = append(c.instructions, ins...)

	return posNewInstructions
}

// Bytecode contains compiler-generated instructions and
// compiler-evaluated constants.
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
