package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/saika-m/saika-lang/internal/ast"
)

// Lexer represents a lexical analyzer for Saika
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
	line         int  // current line
	column       int  // current column
}

// New creates a new Lexer
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position in the input string
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += size
	}

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing the position
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

// NextToken returns the next token
func (l *Lexer) NextToken() ast.Token {
	var tok ast.Token

	l.skipWhitespace()

	// Track token position
	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(ast.PLUS, l.ch)
	case '-':
		tok = newToken(ast.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.BANG, l.ch)
		}
	case '*':
		tok = newToken(ast.ASTERISK, l.ch)
	case '/':
		// Check for comments
		if l.peekChar() == '/' {
			l.skipSingleLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipMultiLineComment()
			return l.NextToken()
		} else {
			tok = newToken(ast.SLASH, l.ch)
		}
	case '.':
		tok = newToken(ast.DOT, l.ch)
	case '<':
		tok = newToken(ast.LT, l.ch)
	case '>':
		tok = newToken(ast.GT, l.ch)
	case ',':
		tok = newToken(ast.COMMA, l.ch)
	case ';':
		tok = newToken(ast.SEMICOLON, l.ch)
	case '(':
		tok = newToken(ast.LPAREN, l.ch)
	case ')':
		tok = newToken(ast.RPAREN, l.ch)
	case '{':
		tok = newToken(ast.LBRACE, l.ch)
	case '}':
		tok = newToken(ast.RBRACE, l.ch)
	case '[':
		tok = newToken(ast.LBRACKET, l.ch)
	case ']':
		tok = newToken(ast.RBRACKET, l.ch)
	case '"':
		tok.Type = ast.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = ast.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = ast.INT
			return tok
		} else {
			tok = newToken(ast.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// skipSingleLineComment skips a single-line comment (// ...)
func (l *Lexer) skipSingleLineComment() {
	l.readChar() // Skip the first '/'
	l.readChar() // Skip the second '/'

	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// skipMultiLineComment skips a multi-line comment (/* ... */)
func (l *Lexer) skipMultiLineComment() {
	l.readChar() // Skip the '/'
	l.readChar() // Skip the '*'

	for {
		if l.ch == 0 { // EOF
			break
		}

		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // Skip the '*'
			l.readChar() // Skip the '/'
			break
		}

		l.readChar()
	}
}

// readIdentifier reads an identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || (l.position != position && isDigit(l.ch)) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString reads a string literal
func (l *Lexer) readString() string {
	l.readChar() // Skip the opening quote
	position := l.position

	for {
		if l.ch == '"' || l.ch == 0 {
			break
		}

		// Handle escape sequences
		if l.ch == '\\' && l.peekChar() == '"' {
			l.readChar() // Skip the backslash
		}

		l.readChar()
	}

	return l.input[position:l.position]
}

// isLetter returns whether the given rune is a letter or underscore
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// isDigit returns whether the given rune is a digit
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// newToken creates a new token
func newToken(tokenType ast.TokenType, ch rune) ast.Token {
	return ast.Token{Type: tokenType, Literal: string(ch)}
}

// LookupIdent looks up an identifier in the keywords map
func LookupIdent(ident string) ast.TokenType {
	if tok, ok := ast.Keywords[ident]; ok {
		return tok
	}
	return ast.IDENT
}
