// Package evaluator defines the logic to
// evaluate Monkey AST expressions.
package evaluator

import (
	"github.com/adamwoolhether/monkeyLang/ast"
	"github.com/adamwoolhether/monkeyLang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval taks an ast.Node and returns an object.Object. Any node
// that fulfills the ast.Node interface can be evaluated. Integer
// and Boolean literals evaluate themselves.
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
		
		// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}
	
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	
	for _, statement := range stmts {
		result = Eval(statement)
	}
	
	return result
}

// nativeBoolToBooleanObject returns one of the predefined TRUE or FALSE
// vars to prevent instantiating a new object.Boolean every time.
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	
	return FALSE
}
