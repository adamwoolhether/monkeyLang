// Package token defines the tokens for use by the lexer.
package token

// TokenType distinguishes the unique token types to represent the source code.
type TokenType string

// Token contains the type of token and its value.
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	
	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1234567890
	
	// Operators
	ASSIGN = "="
	PLUS   = "+"
	
	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	
	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
