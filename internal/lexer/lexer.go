package lexer

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/saika-m/saika-lang/internal/ast"
)

// Lexer represents a lexical analyzer for Saika
type Lexer struct {
	input        string
	position     int    // current position in input (points to current char)
	readPosition int    // current reading position in input (after current char)
	ch           rune   // current char under examination
	line         int    // current line
	column       int    // current column
	fileName     string // source file name for error reporting
}

// New creates a new Lexer
func New(input string) *Lexer {
	return NewWithFilename(input, "")
}

// NewWithFilename creates a new Lexer with the given file name
func NewWithFilename(input string, fileName string) *Lexer {
	l := &Lexer{
		input:    input,
		fileName: fileName,
		line:     1,
		column:   0,
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

// peekCharN returns the Nth next character without advancing the position
func (l *Lexer) peekCharN(n int) rune {
	pos := l.readPosition
	for i := 0; i < n; i++ {
		if pos >= len(l.input) {
			return 0
		}
		r, size := utf8.DecodeRuneInString(l.input[pos:])
		pos += size
		if i == n-1 {
			return r
		}
	}
	return 0
}

// NextToken returns the next token
func (l *Lexer) NextToken() ast.Token {
	var tok ast.Token

	l.skipWhitespace()

	// Store token position for error reporting
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
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.PLUS_EQ, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.INC, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.MINUS_EQ, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.DEC, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.ARROW, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.MINUS, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.BANG, l.ch)
		}
	case '*':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.ASTERISK_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.ASTERISK, l.ch)
		}
	case '/':
		// Check for comments
		if l.peekChar() == '/' {
			l.skipSingleLineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipMultiLineComment()
			return l.NextToken()
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.SLASH_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.SLASH, l.ch)
		}
	case '%':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.PERCENT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.PERCENT, l.ch)
		}
	case '&':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.AND_EQ, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.AND_AND, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.AND, l.ch)
		}
	case '|':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.OR_EQ, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.OR_OR, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.OR, l.ch)
		}
	case '^':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.XOR_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.XOR, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.LT_EQ, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '<' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				literal := string(ch) + string(l.ch)
				l.readChar()
				literal += string(l.ch)
				tok = ast.Token{Type: ast.SHL_EQ, Literal: literal}
			} else {
				tok = ast.Token{Type: ast.SHL, Literal: string(ch) + string(l.ch)}
			}
		} else if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.ARROW, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ast.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = ast.Token{Type: ast.GT_EQ, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				literal := string(ch) + string(l.ch)
				l.readChar()
				literal += string(l.ch)
				tok = ast.Token{Type: ast.SHR_EQ, Literal: literal}
			} else {
				tok = ast.Token{Type: ast.SHR, Literal: string(ch) + string(l.ch)}
			}
		} else {
			tok = newToken(ast.GT, l.ch)
		}
	case '.':
		tok = newToken(ast.DOT, l.ch)
	case ',':
		tok = newToken(ast.COMMA, l.ch)
	case ';':
		tok = newToken(ast.SEMICOLON, l.ch)
	case ':':
		tok = newToken(ast.COLON, l.ch)
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
	case '\'':
		tok.Type = ast.CHAR
		tok.Literal = l.readCharLiteral()
	case 0:
		tok.Literal = ""
		tok.Type = ast.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			return l.readNumber()
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
			// Report error: unclosed comment
			fmt.Printf("Error at line %d, column %d: unclosed comment\n", l.line, l.column)
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

// readNumber reads a number (integer or float)
func (l *Lexer) readNumber() ast.Token {
	position := l.position
	isFloat := false

	// Read integer part
	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // consume the '.'

		// Read fractional part
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Check for scientific notation (e.g., 1e10, 1.5e-5)
	if (l.ch == 'e' || l.ch == 'E') && (isDigit(l.peekChar()) ||
		((l.peekChar() == '+' || l.peekChar() == '-') && isDigit(l.peekCharN(2)))) {
		isFloat = true
		l.readChar() // consume the 'e' or 'E'

		// Handle optional sign
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}

		// Read exponent
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	literal := l.input[position:l.position]

	// Validate the number
	if isFloat {
		if _, err := strconv.ParseFloat(literal, 64); err != nil {
			return ast.Token{
				Type:    ast.ILLEGAL,
				Literal: literal,
				Line:    l.line,
				Column:  l.column - len(literal),
			}
		}
		return ast.Token{
			Type:    ast.FLOAT,
			Literal: literal,
			Line:    l.line,
			Column:  l.column - len(literal),
		}
	} else {
		if _, err := strconv.ParseInt(literal, 10, 64); err != nil {
			return ast.Token{
				Type:    ast.ILLEGAL,
				Literal: literal,
				Line:    l.line,
				Column:  l.column - len(literal),
			}
		}
		return ast.Token{
			Type:    ast.INT,
			Literal: literal,
			Line:    l.line,
			Column:  l.column - len(literal),
		}
	}
}

// readString reads a string literal
func (l *Lexer) readString() string {
	startLine := l.line
	startColumn := l.column

	l.readChar() // Skip the opening quote
	position := l.position

	for {
		if l.ch == '"' || l.ch == 0 {
			break
		}

		// Handle escape sequences
		if l.ch == '\\' {
			l.readChar() // Skip the backslash

			// Handle various escape sequences
			switch l.ch {
			case 'n', 'r', 't', '\\', '"', '\'':
				// Valid escape sequences
			default:
				// Invalid escape sequence
				fmt.Printf("Warning at line %d, column %d: unknown escape sequence \\%c\n",
					l.line, l.column, l.ch)
			}
		}

		// Check for unterminated string at end of file
		if l.ch == 0 {
			fmt.Printf("Error at line %d, column %d: unterminated string literal\n",
				startLine, startColumn)
			break
		}

		l.readChar()
	}

	return l.input[position:l.position]
}

// readChar reads a character literal
func (l *Lexer) readCharLiteral() string {
	startLine := l.line
	startColumn := l.column

	l.readChar() // Skip the opening quote
	position := l.position

	// Handle escape sequences
	if l.ch == '\\' {
		l.readChar() // Skip the backslash

		// Handle various escape sequences
		switch l.ch {
		case 'n', 'r', 't', '\\', '"', '\'':
			// Valid escape sequences
		default:
			// Invalid escape sequence
			fmt.Printf("Warning at line %d, column %d: unknown escape sequence \\%c\n",
				l.line, l.column, l.ch)
		}

		l.readChar()
	} else if l.ch != 0 && l.ch != '\'' {
		l.readChar()
	}

	// Check for proper closing quote
	if l.ch != '\'' {
		fmt.Printf("Error at line %d, column %d: unterminated character literal\n",
			startLine, startColumn)
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

// lookupIdent looks up an identifier in the keywords map
func lookupIdent(ident string) ast.TokenType {
	if tok, ok := ast.Keywords[ident]; ok {
		return tok
	}
	return ast.IDENT
}

// GetFileName returns the name of the file being lexed
func (l *Lexer) GetFileName() string {
	return l.fileName
}

// GetPosition returns the current position in the source file
func (l *Lexer) GetPosition() ast.Position {
	return ast.Position{
		Line:   l.line,
		Column: l.column,
		File:   l.fileName,
	}
}

// GetPositionAt returns the position at the given offset
func (l *Lexer) GetPositionAt(offset int) ast.Position {
	// This is a simple implementation that doesn't account for newlines
	// A complete implementation would need to track line and column for each character
	line := 1
	column := 1

	for i := 0; i < offset && i < len(l.input); i++ {
		if l.input[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}

	return ast.Position{
		Line:   line,
		Column: column,
		File:   l.fileName,
	}
}
