// Package code defines the opcodes and instructions for use by the VM.
package code

import (
	"encoding/binary"
	"fmt"
)

// Insructions Instructions defines the set of bytecode instructions for the VM to run.
type Insructions []byte

// Opcode directs the VM to push something on to the stack.
type Opcode byte

const (
	// OpConstant acts as an index to hold bytecode instructions.
	// Iota will automatically generate increasing byte values.
	OpConstant Opcode = iota
)

// Definition enables looking up how many operands and opcode has
// and what it's human-readable name is.
type Definition struct {
	Name          string // Human-readable name of the opcode.
	OperandWidths []int  // The number of bytes each operand takes up.
}

// definitions holds the map of opcodes and their definitions.
var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

// Lookup enables looking up opcodes in the definitions map.
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

// Make enables building bytecode instructions.
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// Allocate []byte with the length of instructions.
	instruction := make([]byte, instructionLen)
	// Set the opcode as the first byte in the instructions.
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}