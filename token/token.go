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
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	
	LT = "<"
	GT = ">"
	
	EQ     = "=="
	NOT_EQ = "!="
	
	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	
	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	
	// Data Types
	STRING = "STRING"
)

// keywords holds our language keywords, to separate them
// from user-defined identifiers.
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent checks keywords to see if the user-given identifier is a language
// keyword, returning the TokenType constant if so. If not, it returns
// token.IDENT if not, which is the token type for user-defined identifiers.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	
	return IDENT
}
