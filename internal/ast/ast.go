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

// VarStatement represents a variable declaration
type VarStatement struct {
	Token Token // the '变量' token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var out strings.Builder

	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")

	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}

	return out.String()
}

// ConstStatement represents a constant declaration
type ConstStatement struct {
	Token Token // the '常量' token
	Name  *Identifier
	Value Expression
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ConstStatement) String() string {
	var out strings.Builder

	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Name.String())
	out.WriteString(" = ")

	if cs.Value != nil {
		out.WriteString(cs.Value.String())
	}

	return out.String()
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Token       Token // the '返回' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out strings.Builder

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	return out.String()
}

// TypedParam represents a parameter with a type
type TypedParam struct {
	Name *Identifier
	Type *Identifier
}

// FunctionStatement represents a function declaration
type FunctionStatement struct {
	Token      Token // the '數' token
	Name       *Identifier
	Parameters []*TypedParam
	Body       *BlockStatement
	ReturnType *Identifier
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out strings.Builder

	params := []string{}
	for _, p := range fs.Parameters {
		if p.Type != nil {
			params = append(params, p.Name.String()+" "+p.Type.String())
		} else {
			params = append(params, p.Name.String())
		}
	}

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fs.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	if fs.ReturnType != nil {
		out.WriteString(" ")
		out.WriteString(fs.ReturnType.String())
	}

	out.WriteString(" ")
	out.WriteString(fs.Body.String())

	return out.String()
}

// IfStatement represents an if statement
type IfStatement struct {
	Token       Token // the '如果' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out strings.Builder

	out.WriteString("if ")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())

	if is.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(is.Alternative.String())
	}

	return out.String()
}

// ForStatement represents a for loop statement
type ForStatement struct {
	Token     Token // the '循环' token
	Init      Statement
	Condition Expression
	Update    Statement
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out strings.Builder

	out.WriteString("for ")

	if fs.Init != nil {
		out.WriteString(fs.Init.String())
	}
	out.WriteString("; ")

	if fs.Condition != nil {
		out.WriteString(fs.Condition.String())
	}
	out.WriteString("; ")

	if fs.Update != nil {
		out.WriteString(fs.Update.String())
	}

	out.WriteString(" ")
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

// IntegerLiteral represents an integer literal
type IntegerLiteral struct {
	Token Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// StringLiteral represents a string literal
type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	Token Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string {
	if bl.Value {
		return "true"
	}
	return "false"
}

// PrefixExpression represents a prefix expression
type PrefixExpression struct {
	Token    Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

// InfixExpression represents an infix expression
type InfixExpression struct {
	Token    Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

// AssignExpression represents an assignment expression
type AssignExpression struct {
	Token Token // The '=' token
	Left  Expression
	Value Expression
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignExpression) String() string {
	return fmt.Sprintf("%s = %s", ae.Left.String(), ae.Value.String())
}

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
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	IDENT   = "IDENT"
	INT     = "INT"
	STRING  = "STRING"

	// Chinese keywords
	FUNC      = "FUNC"      // 数
	PACKAGE   = "PACKAGE"   // 包
	IMPORT    = "IMPORT"    // 导入
	IF        = "IF"        // 如果
	ELSE      = "ELSE"      // 否则
	FOR       = "FOR"       // 循环
	WHILE     = "WHILE"     // 当
	BREAK     = "BREAK"     // 中断
	CONTINUE  = "CONTINUE"  // 继续
	SWITCH    = "SWITCH"    // 选择
	CASE      = "CASE"      // 情况
	DEFAULT   = "DEFAULT"   // 默认
	RETURN    = "RETURN"    // 返回
	VAR       = "VAR"       // 变量
	CONST     = "CONST"     // 常量
	TRUE      = "TRUE"      // 真
	FALSE     = "FALSE"     // 假
	STRUCT    = "STRUCT"    // 结构
	INTERFACE = "INTERFACE" // 接口
	MAP       = "MAP"       // 映射
	SLICE     = "SLICE"     // 切片
	ARRAY     = "ARRAY"     // 数组
	PUBLIC    = "PUBLIC"    // 公开
	PRIVATE   = "PRIVATE"   // 私有

	// Types
	TYPE_STRING = "TYPE_STRING" // 字符串
	TYPE_INT    = "TYPE_INT"    // 整数
	TYPE_FLOAT  = "TYPE_FLOAT"  // 浮点
	TYPE_BOOL   = "TYPE_BOOL"   // 布尔

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"
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
	LTE    = "<="
	GTE    = ">="
)

// Keywords maps keyword strings to their token types
var Keywords = map[string]TokenType{
	"数":   FUNC,
	"包":   PACKAGE,
	"导入":  IMPORT,
	"如果":  IF,
	"否则":  ELSE,
	"循环":  FOR,
	"当":   WHILE,
	"中断":  BREAK,
	"继续":  CONTINUE,
	"选择":  SWITCH,
	"情况":  CASE,
	"默认":  DEFAULT,
	"返回":  RETURN,
	"变量":  VAR,
	"常量":  CONST,
	"真":   TRUE,
	"假":   FALSE,
	"结构":  STRUCT,
	"接口":  INTERFACE,
	"映射":  MAP,
	"切片":  SLICE,
	"数组":  ARRAY,
	"公开":  PUBLIC,
	"私有":  PRIVATE,
	"字符串": TYPE_STRING,
	"整数":  TYPE_INT,
	"浮点":  TYPE_FLOAT,
	"布尔":  TYPE_BOOL,
}
