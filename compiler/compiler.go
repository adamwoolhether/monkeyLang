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

// Bytecode contains compiler-generated instructions and
// compiler-evaluated constants.
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

// EmittedInstruction allows keeping track of an instruction
// and its opcode after being emitted.
type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

// Compiler holds generated bytecode('instruction'), a pool of constants.
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object

	lastInstruction     EmittedInstruction // The very last instruction emitted.
	previousInstruction EmittedInstruction // The instruction emitted immediately before lastInstruction.

	symbolTable *SymbolTable
}

func New() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		symbolTable:         NewSymbolTable(),
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
		// Special case for the `<` operator to reorder compilation
		// and reuse code.OpGreaterThan.
		if n.Operator == "<" {
			if err := c.Compile(n.Right); err != nil {
				return err
			}

			if err := c.Compile(n.Left); err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)

			return nil
		}

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
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", n.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: n.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if n.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.PrefixExpression:
		if err := c.Compile(n.Right); err != nil {
			return err
		}
		switch n.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", n.Operator)
		}
	case *ast.IfExpression:
		if err := c.Compile(n.Condition); err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus value
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		if err := c.Compile(n.Consequence); err != nil {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		// Emit an `OpJump` with a bogus value.
		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.instructions)
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if n.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			if err := c.Compile(n.Alternative); err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.instructions)
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.BlockStatement:
		for _, s := range n.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		if err := c.Compile(n.Value); err != nil {
			return err
		}
		symbol := c.symbolTable.Define(n.Name.Value)

		c.emit(code.OpSetGlobal, symbol.Index)
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(n.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", n.Value)
		}

		c.emit(code.OpGetGlobal, symbol.Index)
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

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstructions := len(c.instructions)
	c.instructions = append(c.instructions, ins...)

	return posNewInstructions
}

// setLastInstruction tracks the last and second-to-last emitted instruction.
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.previousInstruction = previous
	c.lastInstruction = last
}

// lastInstructionIsPop checks if the last instructions opcode is code.OpPop.
func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

// removeLastPop shortens c.instruction to cut off the last instruction.
func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

// changeOperand allows replacing the operand of an instruction.
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

// replaceInstructions replaces an instruction at an arbitrary
// offset in the instructions slice.
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}
