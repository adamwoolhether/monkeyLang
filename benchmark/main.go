package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/adamwoolhether/monkeyLang/compiler"
	"github.com/adamwoolhether/monkeyLang/evaluator"
	"github.com/adamwoolhether/monkeyLang/lexer"
	"github.com/adamwoolhether/monkeyLang/object"
	"github.com/adamwoolhether/monkeyLang/parser"
	"github.com/adamwoolhether/monkeyLang/vm"
)

var input = `
let fibonacci = fn(x) {
	if (x == 0) { 
		0
	} else {
		if (x == 1) {
			return 1;
		} else {
			fibonacci(x - 1) + fibonacci(x - 2);
		}
	} 
};
   fibonacci(35);
   `

func main() {
	var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")
	flag.Parse()

	var duration time.Duration
	var result object.Object

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if *engine == "vm" {
		comp := compiler.New()
		if err := comp.Compile(program); err != nil {
			fmt.Printf("compiler error: %s", err)
			return
		}

		machine := vm.New(comp.Bytecode())

		start := time.Now()

		if err := machine.Run(); err != nil {
			fmt.Printf("vm error: %s", err)
			return
		}

		duration = time.Since(start)
		result = machine.LastPoppedStackElem()
	} else {
		env := object.NewEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	}

	fmt.Printf("engine=%s, result=%s, duration=%s\n", *engine, result.Inspect(), duration)
}
