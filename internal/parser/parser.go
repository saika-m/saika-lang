package parser

import (
	"fmt"
	"strconv"

	"github.com/saika-m/saika-lang/internal/ast"
	"github.com/saika-m/saika-lang/internal/lexer"
)

// Parser represents a parser for Saika
type Parser struct {
	l         *lexer.Lexer
	curToken  ast.Token
	peekToken ast.Token
	errors    []string

	prefixParseFns map[ast.TokenType]prefixParseFn
	infixParseFns  map[ast.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Precedence levels
const (
	LOWEST      = 1
	EQUALS      = 2
	LESSGREATER = 3
	SUM         = 4
	PRODUCT     = 5
	PREFIX      = 6
	CALL        = 7
	INDEX       = 8
)

// Precedences maps token types to their precedence levels
var precedences = map[ast.TokenType]int{
	ast.EQ:       EQUALS,
	ast.NOT_EQ:   EQUALS,
	ast.ASSIGN:   EQUALS,
	ast.LT:       LESSGREATER,
	ast.GT:       LESSGREATER,
	ast.LTE:      LESSGREATER,
	ast.GTE:      LESSGREATER,
	ast.PLUS:     SUM,
	ast.MINUS:    SUM,
	ast.SLASH:    PRODUCT,
	ast.ASTERISK: PRODUCT,
	ast.PERCENT:  PRODUCT,
	ast.LPAREN:   CALL,
	ast.DOT:      CALL,
}

// New creates a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Initialize prefix parse functions
	// Register prefix parse functions for all token types that can start an expression
	p.prefixParseFns = make(map[ast.TokenType]prefixParseFn)
	p.registerPrefix(ast.IDENT, p.parseIdentifier)
	p.registerPrefix(ast.INT, p.parseIntegerLiteral)
	p.registerPrefix(ast.STRING, p.parseStringLiteral)
	p.registerPrefix(ast.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(ast.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(ast.BANG, p.parsePrefixExpression)
	p.registerPrefix(ast.MINUS, p.parsePrefixExpression)
	p.registerPrefix(ast.LPAREN, p.parseGroupedExpression)

	// Register infix parse functions
	p.infixParseFns = make(map[ast.TokenType]infixParseFn)
	p.registerInfix(ast.PLUS, p.parseInfixExpression)
	p.registerInfix(ast.MINUS, p.parseInfixExpression)
	p.registerInfix(ast.SLASH, p.parseInfixExpression)
	p.registerInfix(ast.ASTERISK, p.parseInfixExpression)
	p.registerInfix(ast.PERCENT, p.parseInfixExpression) // For %
	p.registerInfix(ast.EQ, p.parseInfixExpression)
	p.registerInfix(ast.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(ast.LT, p.parseInfixExpression)
	p.registerInfix(ast.GT, p.parseInfixExpression)
	p.registerInfix(ast.LTE, p.parseInfixExpression)
	p.registerInfix(ast.GTE, p.parseInfixExpression)
	p.registerInfix(ast.ASSIGN, p.parseAssignExpression)
	p.registerInfix(ast.DOT, p.parseMemberExpression)
	p.registerInfix(ast.LPAREN, p.parseCallExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// registerPrefix registers a prefix parse function
func (p *Parser) registerPrefix(tokenType ast.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registers an infix parse function
func (p *Parser) registerInfix(tokenType ast.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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
	case ast.FUNC:
		return p.parseFunctionStatement()
	case ast.VAR:
		return p.parseVarStatement()
	case ast.CONST:
		return p.parseConstStatement()
	case ast.RETURN:
		return p.parseReturnStatement()
	case ast.IF:
		return p.parseIfStatement()
	case ast.FOR:
		return p.parseForStatement()
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

	// Check if the next token is a left parenthesis
	if p.peekTokenIs(ast.LPAREN) {
		// Parenthesized import
		p.nextToken() // Consume the '('

		// Skip any newlines or whitespace
		p.nextToken()

		// Expect a string literal
		if !p.curTokenIs(ast.STRING) {
			p.errors = append(p.errors, fmt.Sprintf("Line %d:%d expected import path to be a string, got %s",
				p.curToken.Line, p.curToken.Column, p.curToken.Type))
			return nil
		}

		// Get the import path
		stmt.Path = p.curToken.Literal

		// Skip to the closing parenthesis
		for !p.peekTokenIs(ast.RPAREN) && !p.peekTokenIs(ast.EOF) {
			p.nextToken()
		}

		// Expect closing parenthesis
		if !p.expectPeek(ast.RPAREN) {
			return nil
		}
	} else {
		// Simple import
		if !p.expectPeek(ast.STRING) {
			return nil
		}

		stmt.Path = p.curToken.Literal
	}

	// Expect semicolon or newline
	if p.peekTokenIs(ast.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseVarStatement parses a variable declaration
func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(ast.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(ast.ASSIGN) {
		return nil
	}

	p.nextToken() // Skip over the '=' token
	stmt.Value = p.parseExpression(LOWEST)

	// Optional semicolon
	if p.peekTokenIs(ast.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseConstStatement parses a constant declaration
func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}

	if !p.expectPeek(ast.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(ast.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(ast.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

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

	// Handle return type
	if p.peekTokenIs(ast.TYPE_INT) || p.peekTokenIs(ast.TYPE_STRING) ||
		p.peekTokenIs(ast.TYPE_FLOAT) || p.peekTokenIs(ast.TYPE_BOOL) {
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
func (p *Parser) parseFunctionParameters() []*ast.TypedParam {
	typedParams := []*ast.TypedParam{}

	if p.peekTokenIs(ast.RPAREN) {
		p.nextToken()
		return typedParams
	}

	p.nextToken()

	// Create parameter with name
	param := &ast.TypedParam{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}

	// Check if there is a type annotation
	if p.peekTokenIs(ast.TYPE_INT) || p.peekTokenIs(ast.TYPE_STRING) ||
		p.peekTokenIs(ast.TYPE_FLOAT) || p.peekTokenIs(ast.TYPE_BOOL) {
		p.nextToken()
		param.Type = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	typedParams = append(typedParams, param)

	for p.peekTokenIs(ast.COMMA) {
		p.nextToken()
		p.nextToken()

		// Create parameter with name
		param := &ast.TypedParam{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}

		// Check if there is a type annotation
		if p.peekTokenIs(ast.TYPE_INT) || p.peekTokenIs(ast.TYPE_STRING) ||
			p.peekTokenIs(ast.TYPE_FLOAT) || p.peekTokenIs(ast.TYPE_BOOL) {
			p.nextToken()
			param.Type = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}

		typedParams = append(typedParams, param)
	}

	if !p.expectPeek(ast.RPAREN) {
		return nil
	}

	return typedParams
}

// parseIfStatement parses an if statement
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(ast.LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(ast.ELSE) {
		p.nextToken()

		if !p.expectPeek(ast.LBRACE) {
			return nil
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

// parseForStatement parses a for statement
func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	// Skip the "循环" token
	p.nextToken()

	// Parse initialization part
	if !p.curTokenIs(ast.SEMICOLON) {
		if p.curTokenIs(ast.VAR) {
			stmt.Init = p.parseVarStatement()
		} else {
			stmt.Init = p.parseExpressionStatement()
		}
	}

	// Skip semicolon after initialization
	if !p.curTokenIs(ast.SEMICOLON) {
		if !p.expectPeek(ast.SEMICOLON) {
			return nil
		}
	} else {
		p.nextToken() // Skip semicolon
	}

	// Parse condition part
	if !p.curTokenIs(ast.SEMICOLON) {
		stmt.Condition = p.parseExpression(LOWEST)
	}

	// Skip semicolon after condition
	if !p.expectPeek(ast.SEMICOLON) {
		return nil
	}

	// Parse update part
	if !p.peekTokenIs(ast.LBRACE) {
		p.nextToken() // Move past the semicolon
		stmt.Update = p.parseExpressionStatement()
	} else {
		p.nextToken() // Move past the semicolon
	}

	// Expect opening brace for the body
	if !p.expectPeek(ast.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
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

// parseExpression parses an expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(ast.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// parseAssignExpression parses an assignment expression
func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	expr := &ast.AssignExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken() // Skip over the '=' token
	expr.Value = p.parseExpression(LOWEST)

	return expr
}

// parseIdentifier parses an identifier
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral parses an integer literal
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// parseStringLiteral parses a string literal
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseBooleanLiteral parses a boolean literal
func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: p.curTokenIs(ast.TRUE),
	}
}

// parsePrefixExpression parses a prefix expression
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression parses an infix expression
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

// parseGroupedExpression parses a grouped expression
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(ast.RPAREN) {
		return nil
	}

	return exp
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

// noPrefixParseFnError adds an error when no prefix parse function exists for the token type
func (p *Parser) noPrefixParseFnError(t ast.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// peekPrecedence returns the precedence of the peek token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence returns the precedence of the current token
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
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
