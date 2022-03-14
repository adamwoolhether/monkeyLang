// Package parser implements a recursive descent parser for the Monkey language.
// It makes use of "Pratt Parsing", aka Top Down Operator Precedence.
package parser

import (
	"fmt"
	"strconv"
	
	"github.com/adamwoolhether/monkeyLang/ast"
	"github.com/adamwoolhether/monkeyLang/lexer"
	"github.com/adamwoolhether/monkeyLang/token"
)

// Define the precedences of the Monkey programming language.
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// precedences defines a table for our precedences,
// associating a token.Type with its precedence const..
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type (
	// prefixParseFn gets called when a token
	// type in prefix position is encountered.
	prefixParseFn func() ast.Expression
	// infixParseFn gets caleld when a token
	// type in infix position is encountered.
	infixParseFn func(ast.Expression) ast.Expression
)

// Parser represents the information necessary to parse a Monkey program.
// curToken allows parsing of the token at the current position, peekToken
// allows the parser to make a decision based on the next token. Token
// types can have up to two parsing funcs asociated with them, depending
// on whether the token is found in a prefix or infix position.
type Parser struct {
	l      *lexer.Lexer
	errors []string
	
	curToken  token.Token
	peekToken token.Token
	
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// registerPrefix adds entries to the Parser's respective function map.
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfixi adds entries to the Parser's respective function map.
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// New returns a pointer to a new parser with the prefixParseFns map
// initialized and registered with the correct parsing function to the
// respective token type.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	
	// Read two tokens, setting curToken and peekToken
	p.nextToken()
	p.nextToken()
	
	return p
}

// parseIdentifier returns an *ast.Identifier with the current token
// in the Token field and literal value of the token in Value. It
// does not advance tokens.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// Errors returns a slice of error strings that the parser may encounter.
func (p *Parser) Errors() []string {
	return p.errors
}

// peekErrors appends an error to p.errors when the type of peekToken
// doesn't match the expectation.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// nextToken is a helper func that advances both curToken and peekToken.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram construct the AST's root node, iterates over
// each token and parses the statement until EOF is reached,
// adds the statement to program.Statements, and returns the
// program's node.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	
	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	
	return program
}

// parseStatement decides how to handle the current token based on its type.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement constructs an *ast.LetStatement node with current
// token.Let token. It expects an identifier token followed by an
// assignment token.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	
	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	
	// TODO: handle expressions. We skip for now until we
	// encounter a semicolon.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	
	return stmt
}

// expectPeek is an assertion function that enforces the correctness
// of token ordering by checking the next token's type.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// parseReturnStatment constructs an *ast.Statement node with
// the current token, skipping over until it encounters a
// semicolon.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	
	p.nextToken()
	
	// TODO: we're skipping expressions until we encounter a semicolon.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	
	return stmt
}

// parseExpressionStatement constructs an *ast.Statement node with
// the current token, skipping over until it encounters a
// semicolon. The semicolon is optional, allowing expression
// statements to accept things like '5 + 5' into the REPL. The
// lowest possible precedence is passed to parseExpression(), as
// nothing has been parsed yet, meaning we can't compare precedences.
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	
	return stmt
}

// noPrefixParseFnError appends a formatted message to
// p.errors when there is no suitable parsing function
// available for the given token.
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse functions for %s found", t)
	p.errors = append(p.errors, msg)
}

// parseExpression checks if the parsing func associated with
// p.curToken.Type is available in the prefix position. Its loop
// tried to find infixParseFns for the enxt token, calling if
// found until a lower-precedence token is encounterd.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		
		p.nextToken()
		
		leftExp = infix(leftExp)
	}
	
	return leftExp
}

// peekPrecedence returns the precedence of the token for
// p.peekToken, returning LOWEST if one not found.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	
	return LOWEST
}

// curPrecedence returns the precedence of the token for
// p.curToken, returning LOWEST if none is found.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	
	return LOWEST
}

// parseIntegerLiteral returns a *ast.IntegerLiteral
// node after it converts p.curToken.Literal to an int64
// saved to the Value field.
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("couldn't parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	
	lit.Value = value
	
	return lit
}

// parsePrefixExpression builds and returns a PrefixExpresion
// AST node for tokens of type token.BANG or token.MINUS. It
// also advances to the next token and calls p.parseExpression
// to correctly parse an express like '-5', as more than one
// token must be consumed, setting the expression to expression.Right.
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	
	p.nextToken()
	
	expression.Right = p.parseExpression(PREFIX)
	
	return expression
}

// parseInfixExpressions uses the given left arguments to construct
// an *ast.InfixExpression node assigning the Left field and the
// current token's precedence (the infix expressions operator) to
// the local var before advancing and assigning the Right field.
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	
	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	
	exp := p.parseExpression(LOWEST)
	
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	
	return exp
}
