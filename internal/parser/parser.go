package parser

import (
	"fmt"

	"github.com/saika-m/saika-lang/internal/ast"
	"github.com/saika-m/saika-lang/internal/lexer"
)

// Parser represents a parser for Saika
type Parser struct {
	l         *lexer.Lexer
	curToken  ast.Token
	peekToken ast.Token
	errors    []string
}

// New creates a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parses a program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.curToken.Type != ast.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses a statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case ast.PACKAGE:
		return p.parsePackageStatement()
	case ast.IMPORT:
		return p.parseImportStatement()
	case ast.SAIKA_FUNC, ast.FUNCTION:
		return p.parseFunctionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parsePackageStatement parses a package statement
func (p *Parser) parsePackageStatement() *ast.PackageStatement {
	stmt := &ast.PackageStatement{Token: p.curToken}

	if !p.expectPeek(ast.IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	// Expect semicolon or newline
	if p.peekTokenIs(ast.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseImportStatement parses an import statement
func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	if !p.expectPeek(ast.STRING) {
		return nil
	}

	stmt.Path = p.curToken.Literal

	// Expect semicolon or newline
	if p.peekTokenIs(ast.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseFunctionStatement parses a function statement
func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{Token: p.curToken}

	if !p.expectPeek(ast.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(ast.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	// Parse return type (if any)
	if p.peekTokenIs(ast.IDENT) {
		p.nextToken()
		stmt.ReturnType = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(ast.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseFunctionParameters parses function parameters
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(ast.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(ast.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(ast.RPAREN) {
		return nil
	}

	return identifiers
}

// parseBlockStatement parses a block statement
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(ast.RBRACE) && !p.curTokenIs(ast.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseExpressionStatement parses an expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(ast.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Precedence levels
const (
	LOWEST      = 1
	EQUALS      = 2
	LESSGREATER = 3
	SUM         = 4
	PRODUCT     = 5
	PREFIX      = 6
	CALL        = 7
)

// parseExpression parses an expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// This is a simplified implementation that can handle identifiers and function calls
	var expr ast.Expression

	// Handle prefix expressions first
	switch p.curToken.Type {
	case ast.IDENT:
		expr = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case ast.STRING:
		expr = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	default:
		return nil
	}

	// Check for member expressions like fmt.Println
	for p.peekTokenIs(ast.DOT) && precedence < CALL {
		p.nextToken() // consume the dot
		expr = p.parseMemberExpression(expr)
	}

	// Check for function calls
	if p.peekTokenIs(ast.LPAREN) {
		p.nextToken() // consume the left paren
		expr = p.parseCallExpression(expr)
	}

	return expr
}

// parseMemberExpression parses a member expression like fmt.Println
func (p *Parser) parseMemberExpression(object ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{
		Token:  p.curToken,
		Object: object,
	}

	p.nextToken()
	exp.Property = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return exp
}

// parseCallExpression parses a call expression like println("hello")
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}

	exp.Arguments = p.parseExpressionList(ast.RPAREN)

	return exp
}

// parseExpressionList parses a list of expressions
func (p *Parser) parseExpressionList(end ast.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(ast.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// curTokenIs returns whether the current token is of the given type
func (p *Parser) curTokenIs(t ast.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs returns whether the peek token is of the given type
func (p *Parser) peekTokenIs(t ast.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek checks if the peek token is of the given type and advances if it is
func (p *Parser) expectPeek(t ast.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// peekError adds an error when the peek token isn't what was expected
func (p *Parser) peekError(t ast.TokenType) {
	msg := fmt.Sprintf("Line %d:%d expected next token to be %s, got %s instead",
		p.peekToken.Line, p.peekToken.Column, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
