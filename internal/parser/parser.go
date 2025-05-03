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

	// Maps for precedence
	prefixParseFns map[ast.TokenType]prefixParseFn
	infixParseFns  map[ast.TokenType]infixParseFn
}

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

// Precedence levels
const (
	LOWEST      = 1
	ASSIGN      = 2  // =, +=, -=, etc.
	LOGICAL_OR  = 3  // ||
	LOGICAL_AND = 4  // &&
	EQUALS      = 5  // ==, !=
	LESSGREATER = 6  // >, <, >=, <=
	BIT_OR      = 7  // |
	BIT_XOR     = 8  // ^
	BIT_AND     = 9  // &
	SHIFT       = 10 // <<, >>
	SUM         = 11 // +, -
	PRODUCT     = 12 // *, /, %
	PREFIX      = 13 // -X, !X, etc.
	POSTFIX     = 14 // ++, --
	CALL        = 15 // myFunction(X)
	INDEX       = 16 // array[index]
	MEMBER      = 17 // foo.bar
)

// Precedence table for tokens
var precedences = map[ast.TokenType]int{
	ast.ASSIGN:      ASSIGN,
	ast.PLUS_EQ:     ASSIGN,
	ast.MINUS_EQ:    ASSIGN,
	ast.ASTERISK_EQ: ASSIGN,
	ast.SLASH_EQ:    ASSIGN,
	ast.PERCENT_EQ:  ASSIGN,
	ast.AND_EQ:      ASSIGN,
	ast.OR_EQ:       ASSIGN,
	ast.XOR_EQ:      ASSIGN,
	ast.SHL_EQ:      ASSIGN,
	ast.SHR_EQ:      ASSIGN,

	ast.OR_OR:   LOGICAL_OR,
	ast.AND_AND: LOGICAL_AND,

	ast.EQ:     EQUALS,
	ast.NOT_EQ: EQUALS,

	ast.LT:    LESSGREATER,
	ast.GT:    LESSGREATER,
	ast.LT_EQ: LESSGREATER,
	ast.GT_EQ: LESSGREATER,

	ast.OR:  BIT_OR,
	ast.XOR: BIT_XOR,
	ast.AND: BIT_AND,

	ast.SHL: SHIFT,
	ast.SHR: SHIFT,

	ast.PLUS:  SUM,
	ast.MINUS: SUM,

	ast.ASTERISK: PRODUCT,
	ast.SLASH:    PRODUCT,
	ast.PERCENT:  PRODUCT,

	ast.INC: POSTFIX,
	ast.DEC: POSTFIX,

	ast.LPAREN:   CALL,
	ast.LBRACKET: INDEX,
	ast.DOT:      MEMBER,
}

// Error represents a parser error
type Error struct {
	Message  string
	Line     int
	Column   int
	FileName string
}

func (e Error) Error() string {
	if e.FileName != "" {
		return fmt.Sprintf("%s:%d:%d: %s", e.FileName, e.Line, e.Column, e.Message)
	}
	return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// New creates a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Initialize maps for prefix and infix parse functions
	p.prefixParseFns = make(map[ast.TokenType]prefixParseFn)
	p.infixParseFns = make(map[ast.TokenType]infixParseFn)

	// Register prefix parse functions
	p.registerPrefix(ast.IDENT, p.parseIdentifier)
	p.registerPrefix(ast.INT, p.parseIntegerLiteral)
	p.registerPrefix(ast.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(ast.STRING, p.parseStringLiteral)
	p.registerPrefix(ast.CHAR, p.parseCharLiteral)
	p.registerPrefix(ast.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(ast.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(ast.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(ast.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(ast.LBRACE, p.parseHashLiteral)
	p.registerPrefix(ast.BANG, p.parsePrefixExpression)
	p.registerPrefix(ast.MINUS, p.parsePrefixExpression)
	p.registerPrefix(ast.PLUS, p.parsePrefixExpression) // Unary plus

	// Register infix parse functions
	p.registerInfix(ast.PLUS, p.parseInfixExpression)
	p.registerInfix(ast.MINUS, p.parseInfixExpression)
	p.registerInfix(ast.ASTERISK, p.parseInfixExpression)
	p.registerInfix(ast.SLASH, p.parseInfixExpression)
	p.registerInfix(ast.PERCENT, p.parseInfixExpression)
	p.registerInfix(ast.EQ, p.parseInfixExpression)
	p.registerInfix(ast.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(ast.LT, p.parseInfixExpression)
	p.registerInfix(ast.GT, p.parseInfixExpression)
	p.registerInfix(ast.LT_EQ, p.parseInfixExpression)
	p.registerInfix(ast.GT_EQ, p.parseInfixExpression)
	p.registerInfix(ast.AND, p.parseInfixExpression)
	p.registerInfix(ast.OR, p.parseInfixExpression)
	p.registerInfix(ast.XOR, p.parseInfixExpression)
	p.registerInfix(ast.SHL, p.parseInfixExpression)
	p.registerInfix(ast.SHR, p.parseInfixExpression)
	p.registerInfix(ast.AND_AND, p.parseInfixExpression)
	p.registerInfix(ast.OR_OR, p.parseInfixExpression)

	p.registerInfix(ast.ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(ast.PLUS_EQ, p.parseAssignmentExpression)
	p.registerInfix(ast.MINUS_EQ, p.parseAssignmentExpression)
	p.registerInfix(ast.ASTERISK_EQ, p.parseAssignmentExpression)
	p.registerInfix(ast.SLASH_EQ, p.parseAssignmentExpression)
	p.registerInfix(ast.PERCENT_EQ, p.parseAssignmentExpression)
	p.registerInfix(ast.AND_EQ, p.parseAssignmentExpression)
	p.registerInfix(ast.OR_EQ, p.parseAssignmentExpression)
	p.registerInfix(ast.XOR_EQ, p.parseAssignmentExpression)
	p.registerInfix(ast.SHL_EQ, p.parseAssignmentExpression)
	p.registerInfix(ast.SHR_EQ, p.parseAssignmentExpression)

	p.registerInfix(ast.LPAREN, p.parseCallExpression)
	p.registerInfix(ast.LBRACKET, p.parseIndexExpression)
	p.registerInfix(ast.DOT, p.parseMemberExpression)
	p.registerInfix(ast.INC, p.parsePostfixExpression)
	p.registerInfix(ast.DEC, p.parsePostfixExpression)

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
		Position:   p.getPosition(),
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
	case ast.LET, ast.VAR, ast.CONST:
		return p.parseVariableStatement()
	case ast.RETURN:
		return p.parseReturnStatement()
	case ast.IF:
		return p.parseIfStatement()
	case ast.FOR:
		return p.parseForStatement()
	case ast.STRUCT:
		return p.parseStructStatement()
	case ast.TYPE:
		return p.parseTypeStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parsePackageStatement parses a package statement
func (p *Parser) parsePackageStatement() *ast.PackageStatement {
	stmt := &ast.PackageStatement{
		Token:    p.curToken,
		Position: p.getPosition(),
	}

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
	stmt := &ast.ImportStatement{
		Token:    p.curToken,
		Position: p.getPosition(),
	}

	// Check for alias (e.g., import fmt "fmt")
	if p.peekTokenIs(ast.IDENT) {
		p.nextToken()
		stmt.Alias = p.curToken.Literal

		if !p.expectPeek(ast.STRING) {
			return nil
		}
	} else if !p.expectPeek(ast.STRING) {
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
	stmt := &ast.FunctionStatement{
		Token:    p.curToken,
		Position: p.getPosition(),
	}

	if !p.expectPeek(ast.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token:    p.curToken,
		Value:    p.curToken.Literal,
		Position: p.getPosition(),
	}

	if !p.expectPeek(ast.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	// Parse return type (if any)
	if p.peekTokenIs(ast.IDENT) {
		p.nextToken()
		stmt.ReturnType = &ast.IdentifierType{
			Token:    p.curToken,
			Name:     p.curToken.Literal,
			Position: p.getPosition(),
		}
	}

	if !p.expectPeek(ast.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseFunctionParameters parses function parameters
func (p *Parser) parseFunctionParameters() []*ast.FunctionParameter {
	params := []*ast.FunctionParameter{}

	if p.peekTokenIs(ast.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := &ast.FunctionParameter{
		Name: &ast.Identifier{
			Token:    p.curToken,
			Value:    p.curToken.Literal,
			Position: p.getPosition(),
		},
		Position: p.getPosition(),
	}

	// Check for type annotation
	if p.peekTokenIs(ast.IDENT) {
		p.nextToken()
		param.Type = &ast.IdentifierType{
			Token:    p.curToken,
			Name:     p.curToken.Literal,
			Position: p.getPosition(),
		}
	}

	params = append(params, param)

	for p.peekTokenIs(ast.COMMA) {
		p.nextToken()
		p.nextToken()

		param := &ast.FunctionParameter{
			Name: &ast.Identifier{
				Token:    p.curToken,
				Value:    p.curToken.Literal,
				Position: p.getPosition(),
			},
			Position: p.getPosition(),
		}

		// Check for type annotation
		if p.peekTokenIs(ast.IDENT) {
			p.nextToken()
			param.Type = &ast.IdentifierType{
				Token:    p.curToken,
				Name:     p.curToken.Literal,
				Position: p.getPosition(),
			}
		}

		params = append(params, param)
	}

	if !p.expectPeek(ast.RPAREN) {
		return nil
	}

	return params
}

// parseVariableStatement parses a variable declaration statement
func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	stmt := &ast.VariableStatement{
		Token:    p.curToken,
		Position: p.getPosition(),
		Const:    p.curToken.Type == ast.CONST,
	}

	if !p.expectPeek(ast.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token:    p.curToken,
		Value:    p.curToken.Literal,
		Position: p.getPosition(),
	}

	// Check for type annotation
	if p.peekTokenIs(ast.IDENT) {
		p.nextToken()
		stmt.Type = &ast.IdentifierType{
			Token:    p.curToken,
			Name:     p.curToken.Literal,
			Position: p.getPosition(),
		}
	}

	// Check for assignment
	if p.peekTokenIs(ast.ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	// Expect semicolon
	if p.peekTokenIs(ast.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token:    p.curToken,
		Position: p.getPosition(),
	}

	p.nextToken()

	// Parse return value if present
	if !p.curTokenIs(ast.SEMICOLON) && !p.curTokenIs(ast.RBRACE) {
		stmt.Value = p.parseExpression(LOWEST)
	}

	// Expect semicolon
	if p.peekTokenIs(ast.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseIfStatement parses an if statement
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{
		Token:    p.curToken,
		Position: p.getPosition(),
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(ast.LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	// Check for else clause
	if p.peekTokenIs(ast.ELSE) {
		p.nextToken()

		if !p.expectPeek(ast.LBRACE) {
			return nil
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

// parseForStatement parses a for statement or for-range statement
func (p *Parser) parseForStatement() ast.Statement {
	token := p.curToken
	position := p.getPosition()

	p.nextToken()

	// Check for range-based for loop
	if p.curTokenIs(ast.IDENT) && p.peekTokenIs(ast.COMMA) {
		// This is likely a range loop with key and value
		key := &ast.Identifier{
			Token:    p.curToken,
			Value:    p.curToken.Literal,
			Position: p.getPosition(),
		}

		p.nextToken() // consume comma
		p.nextToken() // move to value ident

		value := &ast.Identifier{
			Token:    p.curToken,
			Value:    p.curToken.Literal,
			Position: p.getPosition(),
		}

		if !p.expectPeek(ast.ASSIGN) || !p.peekTokenIs(ast.RANGE) {
			p.errors = append(p.errors, fmt.Sprintf("Expected := range at %d:%d", p.peekToken.Line, p.peekToken.Column))
			return nil
		}

		p.nextToken() // consume :=
		p.nextToken() // consume range

		p.nextToken() // move to collection expression
		collection := p.parseExpression(LOWEST)

		if !p.expectPeek(ast.LBRACE) {
			return nil
		}

		body := p.parseBlockStatement()

		return &ast.RangeStatement{
			Token:      token,
			Key:        key,
			Value:      value,
			Collection: collection,
			Body:       body,
			Position:   position,
		}
	} else if p.curTokenIs(ast.IDENT) && p.peekTokenIs(ast.ASSIGN) && p.peekTokenIs(ast.RANGE) {
		// This is a range loop with just value
		value := &ast.Identifier{
			Token:    p.curToken,
			Value:    p.curToken.Literal,
			Position: p.getPosition(),
		}

		p.nextToken() // consume :=
		p.nextToken() // consume range

		p.nextToken() // move to collection expression
		collection := p.parseExpression(LOWEST)

		if !p.expectPeek(ast.LBRACE) {
			return nil
		}

		body := p.parseBlockStatement()

		return &ast.RangeStatement{
			Token:      token,
			Value:      value,
			Collection: collection,
			Body:       body,
			Position:   position,
		}
	} else {
		// This is a regular for loop
		var init ast.Statement
		var condition ast.Expression
		var post ast.Statement

		// Parse initialization statement
		if !p.curTokenIs(ast.SEMICOLON) {
			init = p.parseStatement()

			if !p.expectPeek(ast.SEMICOLON) {
				return nil
			}
		} else {
			p.nextToken() // consume the semicolon
		}

		// Parse condition expression
		if !p.curTokenIs(ast.SEMICOLON) {
			condition = p.parseExpression(LOWEST)

			if !p.expectPeek(ast.SEMICOLON) {
				return nil
			}
		} else {
			p.nextToken() // consume the semicolon
		}

		// Parse post statement
		if !p.curTokenIs(ast.LBRACE) {
			post = p.parseStatement()

			if !p.expectPeek(ast.LBRACE) {
				return nil
			}
		}

		body := p.parseBlockStatement()

		return &ast.ForStatement{
			Token:     token,
			Init:      init,
			Condition: condition,
			Post:      post,
			Body:      body,
			Position:  position,
		}
	}
}

// parseStructStatement parses a struct declaration
func (p *Parser) parseStructStatement() *ast.StructLiteral {
	stmt := &ast.StructLiteral{
		Token:    p.curToken,
		Fields:   make(map[string]ast.TypeExpression),
		Position: p.getPosition(),
	}

	// Check for struct name
	if p.peekTokenIs(ast.IDENT) {
		p.nextToken()
		stmt.Name = &ast.Identifier{
			Token:    p.curToken,
			Value:    p.curToken.Literal,
			Position: p.getPosition(),
		}
	}

	if !p.expectPeek(ast.LBRACE) {
		return nil
	}

	p.nextToken() // consume opening brace

	// Parse fields
	for !p.curTokenIs(ast.RBRACE) && !p.curTokenIs(ast.EOF) {
		if !p.curTokenIs(ast.IDENT) {
			p.errors = append(p.errors, fmt.Sprintf("Expected field name at %d:%d", p.curToken.Line, p.curToken.Column))
			return nil
		}

		fieldName := p.curToken.Literal

		if !p.expectPeek(ast.IDENT) {
			return nil
		}

		fieldType := &ast.IdentifierType{
			Token:    p.curToken,
			Name:     p.curToken.Literal,
			Position: p.getPosition(),
		}

		stmt.Fields[fieldName] = fieldType

		p.nextToken() // move to next field or closing brace
	}

	return stmt
}

// parseTypeStatement parses a type declaration
func (p *Parser) parseTypeStatement() ast.Statement {
	// TODO: Implement type declarations
	p.errors = append(p.errors, fmt.Sprintf("Type declarations not yet implemented at %d:%d", p.curToken.Line, p.curToken.Column))
	return nil
}

// parseBlockStatement parses a block statement
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Statement{},
		Position:   p.getPosition(),
	}

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
	stmt := &ast.ExpressionStatement{
		Token:    p.curToken,
		Position: p.getPosition(),
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(ast.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpression parses an expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Check for prefix function
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	// Process any infix expressions
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

// parseIdentifier parses an identifier
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token:    p.curToken,
		Value:    p.curToken.Literal,
		Position: p.getPosition(),
	}
}

// parseIntegerLiteral parses an integer literal
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token:    p.curToken,
		Position: p.getPosition(),
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// parseFloatLiteral parses a floating-point literal
func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{
		Token:    p.curToken,
		Position: p.getPosition(),
	}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// parseStringLiteral parses a string literal
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token:    p.curToken,
		Value:    p.curToken.Literal,
		Position: p.getPosition(),
	}
}

// parseCharLiteral parses a character literal
func (p *Parser) parseCharLiteral() ast.Expression {
	// TODO: Implement character literals
	p.errors = append(p.errors, fmt.Sprintf("Character literals not yet implemented at %d:%d", p.curToken.Line, p.curToken.Column))
	return nil
}

// parseBooleanLiteral parses a boolean literal
func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token:    p.curToken,
		Value:    p.curTokenIs(ast.TRUE),
		Position: p.getPosition(),
	}
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

// parseArrayLiteral parses an array literal
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{
		Token:    p.curToken,
		Position: p.getPosition(),
	}

	array.Elements = p.parseExpressionList(ast.RBRACKET)

	return array
}

// parseHashLiteral parses a hash literal
func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{
		Token:    p.curToken,
		Pairs:    make(map[ast.Expression]ast.Expression),
		Position: p.getPosition(),
	}

	for !p.peekTokenIs(ast.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(ast.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(ast.RBRACE) && !p.expectPeek(ast.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(ast.RBRACE) {
		return nil
	}

	return hash
}

// parsePrefixExpression parses a prefix expression
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.UnaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Position: p.getPosition(),
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression parses an infix expression
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.BinaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
		Position: p.getPosition(),
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseAssignmentExpression parses an assignment expression
func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	expression := &ast.AssignmentExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
		Position: p.getPosition(),
	}

	p.nextToken()
	expression.Right = p.parseExpression(LOWEST)

	return expression
}

// parseCallExpression parses a call expression
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
		Position: p.getPosition(),
	}

	expression.Arguments = p.parseExpressionList(ast.RPAREN)

	return expression
}

// parseIndexExpression parses an index expression
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{
		Token:    p.curToken,
		Left:     left,
		Position: p.getPosition(),
	}

	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(ast.RBRACKET) {
		return nil
	}

	return expression
}

// parseMemberExpression parses a member expression
func (p *Parser) parseMemberExpression(object ast.Expression) ast.Expression {
	expression := &ast.MemberExpression{
		Token:    p.curToken,
		Object:   object,
		Position: p.getPosition(),
	}

	p.nextToken()

	if !p.curTokenIs(ast.IDENT) {
		msg := fmt.Sprintf("expected identifier after dot, got %s", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	expression.Property = &ast.Identifier{
		Token:    p.curToken,
		Value:    p.curToken.Literal,
		Position: p.getPosition(),
	}

	return expression
}

// parsePostfixExpression parses a postfix expression like i++, i--
func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.UnaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Right:    left, // In postfix, the operand is on the left
		Position: p.getPosition(),
	}

	return expression
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
	}

	p.peekError(t)
	return false
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

// noPrefixParseFnError adds an error when no prefix parse function is found
func (p *Parser) noPrefixParseFnError(t ast.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// peekError adds an error when the peek token isn't what was expected
func (p *Parser) peekError(t ast.TokenType) {
	msg := fmt.Sprintf("Line %d:%d expected next token to be %s, got %s instead",
		p.peekToken.Line, p.peekToken.Column, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// getPosition returns the current position in the source file
func (p *Parser) getPosition() ast.Position {
	return ast.Position{
		Line:   p.curToken.Line,
		Column: p.curToken.Column,
		File:   p.l.GetFileName(),
	}
}
