// Package ast defines the Monkey Abstract Syntax Tree to conduct parsing of source code.
// Statements in Monkey consist of identifiers and expressions. In the following example,
// x, y, and add are identifiers. 10, 15, and the function literal are expressions.
// let x = 10;
// let y = 15;
// let add = fn(a, b) {
//   return a + b;
// };
package ast

import "github.com/adamwoolhether/monkeyLang/token"

// Node defines the ontract for all nodes in the Monkey AST.
// TokenLiteral is used for debugging and testing.
type Node interface {
	TokenLiteral() string
}

// Statement nodes to not produce a value. ex:
// let x = 5
// return 5
type Statement interface {
	Node
	statementNode()
}

// Expression nodes produce a value. ex:
// 5
// add(5, 5)
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root node of every AST produced
// by the Monkey parser. Valid Monkey programs are a
// series of statements.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement represents a let statement in Monkey.
// It's methods satisfy the Statement and Node interfaces.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Identifier represents the identifiers of a binding.
// It satisfies the Expression interface.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// ReturnStatement represents a return statement in Monkey.
// It's methods satisfy the Statement and Node interfaces.
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
