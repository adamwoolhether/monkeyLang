// Package code defines the opcodes and instructions for use by the VM.
package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Instructions defines the set of bytecode instructions for the VM to run.
type Instructions []byte

// Instructions.String allows pretty printing of Instructions,
// it is essentially a mini-disassembler.
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstructions(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstructions(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// Opcode directs the VM to push something on to the stack.
type Opcode byte

const (
	// OpConstant acts as an index to hold bytecode instructions.
	// Iota will automatically generate increasing byte values.
	OpConstant Opcode = iota
	// OpAdd tells the VM to pop the two leftmost elements off
	// the stack, add them, and push the result back on stack.
	OpAdd
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
	OpAdd:      {"OpAdd", []int{}}, // Empty slice signifies no operands.
}

// Lookup enables looking up opcodes in the definitions map.
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

// Make enables building bytecode instructions by encoding operands.
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

// ReadOperands decodes a given set of encoded instructions and decodes them.
// It exists as the counterpart to Make.
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	// Use *Definition of an opcode to determine how wide operands are.
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	// Range through Instructions to read in and convert as many bytes as defined in definition.
	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

// ReadUint16 enables use by VM to skip looking up definition needed by ReadOperands.
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
