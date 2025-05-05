package ast

import (
	"fmt"
	"strings"
)

// Node represents a node in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents a statement node in the AST
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node in the AST
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// PackageStatement represents a package declaration
type PackageStatement struct {
	Token Token
	Name  string
}

func (ps *PackageStatement) statementNode()       {}
func (ps *PackageStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PackageStatement) String() string {
	return fmt.Sprintf("package %s", ps.Name)
}

// ImportStatement represents an import declaration
type ImportStatement struct {
	Token Token
	Path  string
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	return fmt.Sprintf("import %s", is.Path)
}

// FunctionStatement represents a function declaration
type FunctionStatement struct {
	Token      Token // the '數' token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
	ReturnType *Identifier
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out strings.Builder

	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fs.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")

	if fs.ReturnType != nil {
		out.WriteString(fs.ReturnType.String())
		out.WriteString(" ")
	}

	out.WriteString(fs.Body.String())

	return out.String()
}

// BlockStatement represents a block of statements enclosed in { }
type BlockStatement struct {
	Token      Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out strings.Builder

	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")

	return out.String()
}

// ExpressionStatement represents an expression used as a statement
type ExpressionStatement struct {
	Token      Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Identifier represents an identifier
type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// MemberExpression represents a member expression like fmt.Println
type MemberExpression struct {
	Token    Token // the '.' token
	Object   Expression
	Property Expression
}

func (me *MemberExpression) expressionNode()      {}
func (me *MemberExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MemberExpression) String() string {
	return fmt.Sprintf("%s.%s", me.Object.String(), me.Property.String())
}

// CallExpression represents a function call expression
type CallExpression struct {
	Token     Token // The '(' token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out strings.Builder

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// StringLiteral represents a string literal
type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// Token represents a token produced by the lexer
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// TokenType represents the type of a token
type TokenType string

// Token types
const (
	ILLEGAL    = "ILLEGAL"
	EOF        = "EOF"
	IDENT      = "IDENT"
	INT        = "INT"
	STRING     = "STRING"
	SAIKA_FUNC = "SAIKA_FUNC"
	FUNCTION   = "FUNCTION"
	LET        = "LET"
	TRUE       = "TRUE"
	FALSE      = "FALSE"
	IF         = "IF"
	ELSE       = "ELSE"
	RETURN     = "RETURN"
	PACKAGE    = "PACKAGE"
	IMPORT     = "IMPORT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	DOT      = "."

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// Comparisons
	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"
)

// Keywords maps keyword strings to their token types
var Keywords = map[string]TokenType{
	"func":    FUNCTION,
	"數":       SAIKA_FUNC,
	"let":     LET,
	"true":    TRUE,
	"false":   FALSE,
	"if":      IF,
	"else":    ELSE,
	"return":  RETURN,
	"package": PACKAGE,
	"import":  IMPORT,
}
