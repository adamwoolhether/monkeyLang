// Package lexer takes source code as input and outputs the tokens that represent the source code.
package lexer

import "github.com/adamwoolhether/monkeyLang/token"

// Lexer contains the inputted source code and defines methods
// to obtain information about the input's characters.
type Lexer struct {
	input        string
	position     int  // current position in input (points to the current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New returns a new Lexer with l.ch, l.position, and l.readPosition already initialized.
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	
	return l
}

// readChar gives the next character and advances to the next position in the input string.
// If the end of input is reached, ch is set to the ASCII code for "NUL", 0.
// Currently only ASCII chars are supported. Unicode & UTF-8 support require conversion of
// l.ch from a byte to a rune, as well as changing how the next char is read, as it could
// be multiple bytes. // TODO: Implement full Unicode support for Monkey.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken determines which token corresponds to the character
// being examined and advances to the next position.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}
	
	l.readChar()
	return tok
}

// newToken intializes a token.Token based on the given type.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
