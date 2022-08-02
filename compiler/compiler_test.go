package compiler

import (
	"fmt"
	"testing"

	"github.com/adamwoolhether/monkeyLang/ast"
	"github.com/adamwoolhether/monkeyLang/code"
	"github.com/adamwoolhether/monkeyLang/lexer"
	"github.com/adamwoolhether/monkeyLang/object"
	"github.com/adamwoolhether/monkeyLang/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Insructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Insructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}

	runCompilerTests(t, tests)
}

// runCompilerTests takes Monkey code input, parses it, produces and AST, and runs it
// through the compiler before making assertions about the bytecode produced by the
// compiler.
func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		// Lex, parse, and return ast.
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		// Ensure that the bytecode.Instructions produces are correct.
		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testInstructions(expected []code.Insructions, actual code.Insructions) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong isntructions length.\nwant=%q\ngot =%q", concatted, actual)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instructions at %d.\nwant=%q\ngot =%q", i, concatted, actual)
		}
	}

	return nil
}

// concatInstructions helps us turn a slice of slices into a flattened slice.
func concatInstructions(s []code.Insructions) code.Insructions {
	out := code.Insructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}

	return out
}

// testConstants iterates through expected constants, comparing them with actual constants
// that the compiler produced.
func testConstants(t *testing.T, expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got =%d, want=%d", len(actual), len(expected))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T, want=%d", result.Value, expected)
	}

	return nil
}