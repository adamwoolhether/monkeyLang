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
	
	l.skipWhitespace()
	
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
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	
	l.readChar()
	return tok
}

// newToken initializes a token.Token based on the given type.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// readIdentifier reads an identifer's value and advances the
// lexer's position until a non-letter char is encountered.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	
	return l.input[position:l.position]
}

// isLetter checks whether the given argument is a letter or not. It allows
// the char '_' to be treated as a letter, allowing it to be used in
// identifiers and keywords, ex: foo_bar.
// To allow other identifiers like ! or ?, add them here.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// skipWhitespace skips over whitespace, as Monkey does give them meaning.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readNumber reads an int from the char at the position and
// advances the lexer's position until a non-int is encountered.
// It only supports ints, and not floats, hex, or ocatal notions
// for the sake of simplicity.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	
	return l.input[position:l.position]
}

// isDigit checks whether the passed byte is a digit between 0 and 9.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
