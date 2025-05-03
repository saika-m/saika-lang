package ast

import (
	"fmt"
	"strings"
)

// Node represents a node in the AST
type Node interface {
	TokenLiteral() string
	String() string
	GetPosition() Position
}

// Position represents a position in the source code
type Position struct {
	Line   int
	Column int
	File   string
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
	Position   Position
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

func (p *Program) GetPosition() Position {
	return p.Position
}

// PackageStatement represents a package declaration
type PackageStatement struct {
	Token    Token
	Name     string
	Position Position
}

func (ps *PackageStatement) statementNode()       {}
func (ps *PackageStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PackageStatement) String() string {
	return fmt.Sprintf("package %s", ps.Name)
}
func (ps *PackageStatement) GetPosition() Position {
	return ps.Position
}

// ImportStatement represents an import declaration
type ImportStatement struct {
	Token    Token
	Path     string
	Alias    string // For aliased imports (e.g., import fmt "fmt")
	Position Position
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	if is.Alias != "" {
		return fmt.Sprintf("import %s %s", is.Alias, is.Path)
	}
	return fmt.Sprintf("import %s", is.Path)
}
func (is *ImportStatement) GetPosition() Position {
	return is.Position
}

// FunctionStatement represents a function declaration
type FunctionStatement struct {
	Token      Token // the '數' token
	Name       *Identifier
	Parameters []*FunctionParameter
	Body       *BlockStatement
	ReturnType TypeExpression
	Position   Position
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
func (fs *FunctionStatement) GetPosition() Position {
	return fs.Position
}

// FunctionParameter represents a function parameter with type
type FunctionParameter struct {
	Name     *Identifier
	Type     TypeExpression
	Position Position
}

func (fp *FunctionParameter) String() string {
	if fp.Type != nil {
		return fmt.Sprintf("%s %s", fp.Name.String(), fp.Type.String())
	}
	return fp.Name.String()
}
func (fp *FunctionParameter) GetPosition() Position {
	return fp.Position
}

// TypeExpression is an interface for types
type TypeExpression interface {
	Expression
	typeExpressionNode()
}

// IdentifierType represents a named type (e.g., int, string)
type IdentifierType struct {
	Token    Token
	Name     string
	Position Position
}

func (it *IdentifierType) expressionNode()       {}
func (it *IdentifierType) typeExpressionNode()   {}
func (it *IdentifierType) TokenLiteral() string  { return it.Token.Literal }
func (it *IdentifierType) String() string        { return it.Name }
func (it *IdentifierType) GetPosition() Position { return it.Position }

// ArrayType represents an array type (e.g., []string)
type ArrayType struct {
	Token       Token
	ElementType TypeExpression
	Position    Position
}

func (at *ArrayType) expressionNode()       {}
func (at *ArrayType) typeExpressionNode()   {}
func (at *ArrayType) TokenLiteral() string  { return at.Token.Literal }
func (at *ArrayType) String() string        { return fmt.Sprintf("[]%s", at.ElementType.String()) }
func (at *ArrayType) GetPosition() Position { return at.Position }

// MapType represents a map type (e.g., map[string]int)
type MapType struct {
	Token     Token
	KeyType   TypeExpression
	ValueType TypeExpression
	Position  Position
}

func (mt *MapType) expressionNode()      {}
func (mt *MapType) typeExpressionNode()  {}
func (mt *MapType) TokenLiteral() string { return mt.Token.Literal }
func (mt *MapType) String() string {
	return fmt.Sprintf("map[%s]%s", mt.KeyType.String(), mt.ValueType.String())
}
func (mt *MapType) GetPosition() Position { return mt.Position }

// BlockStatement represents a block of statements enclosed in { }
type BlockStatement struct {
	Token      Token // the '{' token
	Statements []Statement
	Position   Position
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out strings.Builder

	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	out.WriteString("}")

	return out.String()
}
func (bs *BlockStatement) GetPosition() Position {
	return bs.Position
}

// VariableStatement represents a variable declaration
type VariableStatement struct {
	Token    Token // the 'let' or 'var' token
	Name     *Identifier
	Type     TypeExpression
	Value    Expression
	Const    bool // Whether this is a constant
	Position Position
}

func (vs *VariableStatement) statementNode()       {}
func (vs *VariableStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VariableStatement) String() string {
	var out strings.Builder

	if vs.Const {
		out.WriteString("const ")
	} else {
		out.WriteString("var ")
	}

	out.WriteString(vs.Name.String())

	if vs.Type != nil {
		out.WriteString(" ")
		out.WriteString(vs.Type.String())
	}

	if vs.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vs.Value.String())
	}

	return out.String()
}
func (vs *VariableStatement) GetPosition() Position {
	return vs.Position
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Token    Token // the 'return' token
	Value    Expression
	Position Position
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out strings.Builder

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}

	return out.String()
}
func (rs *ReturnStatement) GetPosition() Position {
	return rs.Position
}

// IfStatement represents an if statement
type IfStatement struct {
	Token       Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	Position    Position
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
func (is *IfStatement) GetPosition() Position {
	return is.Position
}

// ForStatement represents a for loop
type ForStatement struct {
	Token     Token // the 'for' token
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
	Position  Position
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

	if fs.Post != nil {
		out.WriteString(fs.Post.String())
	}

	out.WriteString(" ")
	out.WriteString(fs.Body.String())

	return out.String()
}
func (fs *ForStatement) GetPosition() Position {
	return fs.Position
}

// RangeStatement represents a for-range loop
type RangeStatement struct {
	Token      Token // the 'for' token
	Key        *Identifier
	Value      *Identifier
	Collection Expression
	Body       *BlockStatement
	Position   Position
}

func (rs *RangeStatement) statementNode()       {}
func (rs *RangeStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *RangeStatement) String() string {
	var out strings.Builder

	out.WriteString("for ")

	if rs.Key != nil {
		out.WriteString(rs.Key.String())

		if rs.Value != nil {
			out.WriteString(", ")
			out.WriteString(rs.Value.String())
		}

		out.WriteString(" := ")
	}

	out.WriteString("range ")
	out.WriteString(rs.Collection.String())
	out.WriteString(" ")
	out.WriteString(rs.Body.String())

	return out.String()
}
func (rs *RangeStatement) GetPosition() Position {
	return rs.Position
}

// ExpressionStatement represents an expression used as a statement
type ExpressionStatement struct {
	Token      Token
	Expression Expression
	Position   Position
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
func (es *ExpressionStatement) GetPosition() Position {
	return es.Position
}

// Identifier represents an identifier
type Identifier struct {
	Token    Token
	Value    string
	Position Position
}

func (i *Identifier) expressionNode()       {}
func (i *Identifier) TokenLiteral() string  { return i.Token.Literal }
func (i *Identifier) String() string        { return i.Value }
func (i *Identifier) GetPosition() Position { return i.Position }

// AssignmentExpression represents an assignment operation
type AssignmentExpression struct {
	Token    Token // The '=' token
	Left     Expression
	Operator string // =, +=, -=, etc.
	Right    Expression
	Position Position
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	return fmt.Sprintf("%s %s %s", ae.Left.String(), ae.Operator, ae.Right.String())
}
func (ae *AssignmentExpression) GetPosition() Position {
	return ae.Position
}

// BinaryExpression represents a binary operation
type BinaryExpression struct {
	Token    Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
	Position Position
}

func (be *BinaryExpression) expressionNode()      {}
func (be *BinaryExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Left.String(), be.Operator, be.Right.String())
}
func (be *BinaryExpression) GetPosition() Position {
	return be.Position
}

// UnaryExpression represents a unary operation
type UnaryExpression struct {
	Token    Token // The operator token, e.g. !
	Operator string
	Right    Expression
	Position Position
}

func (ue *UnaryExpression) expressionNode()      {}
func (ue *UnaryExpression) TokenLiteral() string { return ue.Token.Literal }
func (ue *UnaryExpression) String() string {
	return fmt.Sprintf("(%s%s)", ue.Operator, ue.Right.String())
}
func (ue *UnaryExpression) GetPosition() Position {
	return ue.Position
}

// MemberExpression represents a member expression like fmt.Println
type MemberExpression struct {
	Token    Token // the '.' token
	Object   Expression
	Property Expression
	Position Position
}

func (me *MemberExpression) expressionNode()      {}
func (me *MemberExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MemberExpression) String() string {
	return fmt.Sprintf("%s.%s", me.Object.String(), me.Property.String())
}
func (me *MemberExpression) GetPosition() Position {
	return me.Position
}

// IndexExpression represents an index expression like array[0]
type IndexExpression struct {
	Token    Token // the '[' token
	Left     Expression
	Index    Expression
	Position Position
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return fmt.Sprintf("%s[%s]", ie.Left.String(), ie.Index.String())
}
func (ie *IndexExpression) GetPosition() Position {
	return ie.Position
}

// CallExpression represents a function call expression
type CallExpression struct {
	Token     Token // The '(' token
	Function  Expression
	Arguments []Expression
	Position  Position
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
func (ce *CallExpression) GetPosition() Position {
	return ce.Position
}

// StringLiteral represents a string literal
type StringLiteral struct {
	Token    Token
	Value    string
	Position Position
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }
func (sl *StringLiteral) GetPosition() Position {
	return sl.Position
}

// IntegerLiteral represents an integer literal
type IntegerLiteral struct {
	Token    Token
	Value    int64
	Position Position
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
func (il *IntegerLiteral) GetPosition() Position {
	return il.Position
}

// FloatLiteral represents a floating-point literal
type FloatLiteral struct {
	Token    Token
	Value    float64
	Position Position
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }
func (fl *FloatLiteral) GetPosition() Position {
	return fl.Position
}

// BooleanLiteral represents a boolean literal
type BooleanLiteral struct {
	Token    Token
	Value    bool
	Position Position
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }
func (bl *BooleanLiteral) GetPosition() Position {
	return bl.Position
}

// ArrayLiteral represents an array literal
type ArrayLiteral struct {
	Token    Token // the '[' token
	Elements []Expression
	Position Position
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out strings.Builder

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (al *ArrayLiteral) GetPosition() Position {
	return al.Position
}

// HashLiteral represents a hash literal
type HashLiteral struct {
	Token    Token // the '{' token
	Pairs    map[Expression]Expression
	Position Position
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out strings.Builder

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
func (hl *HashLiteral) GetPosition() Position {
	return hl.Position
}

// StructLiteral represents a struct literal
type StructLiteral struct {
	Token    Token // the 'struct' token
	Name     *Identifier
	Fields   map[string]TypeExpression
	Position Position
}

func (sl *StructLiteral) expressionNode()      {}
func (sl *StructLiteral) typeExpressionNode()  {}
func (sl *StructLiteral) statementNode()       {}
func (sl *StructLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StructLiteral) String() string {
	var out strings.Builder

	out.WriteString("struct ")
	if sl.Name != nil {
		out.WriteString(sl.Name.String() + " ")
	}
	out.WriteString("{\n")

	for name, typ := range sl.Fields {
		out.WriteString(fmt.Sprintf("\t%s %s\n", name, typ.String()))
	}

	out.WriteString("}")

	return out.String()
}
func (sl *StructLiteral) GetPosition() Position {
	return sl.Position
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

	// Identifiers + literals
	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"
	CHAR   = "CHAR"

	// Keywords
	SAIKA_FUNC = "SAIKA_FUNC" // 數
	FUNCTION   = "FUNCTION"   // func
	LET        = "LET"        // let
	VAR        = "VAR"        // var
	CONST      = "CONST"      // const
	TRUE       = "TRUE"       // true
	FALSE      = "FALSE"      // false
	IF         = "IF"         // if
	ELSE       = "ELSE"       // else
	RETURN     = "RETURN"     // return
	PACKAGE    = "PACKAGE"    // package
	IMPORT     = "IMPORT"     // import
	FOR        = "FOR"        // for
	RANGE      = "RANGE"      // range
	BREAK      = "BREAK"      // break
	CONTINUE   = "CONTINUE"   // continue
	STRUCT     = "STRUCT"     // struct
	INTERFACE  = "INTERFACE"  // interface
	MAP        = "MAP"        // map
	CHAN       = "CHAN"       // chan
	GO         = "GO"         // go
	SELECT     = "SELECT"     // select
	CASE       = "CASE"       // case
	DEFAULT    = "DEFAULT"    // default
	SWITCH     = "SWITCH"     // switch
	TYPE       = "TYPE"       // type

	// Operators
	ASSIGN      = "="
	PLUS        = "+"
	PLUS_EQ     = "+="
	MINUS       = "-"
	MINUS_EQ    = "-="
	BANG        = "!"
	ASTERISK    = "*"
	ASTERISK_EQ = "*="
	SLASH       = "/"
	SLASH_EQ    = "/="
	PERCENT     = "%"
	PERCENT_EQ  = "%="
	DOT         = "."
	AND         = "&"
	AND_EQ      = "&="
	OR          = "|"
	OR_EQ       = "|="
	XOR         = "^"
	XOR_EQ      = "^="
	SHL         = "<<"
	SHL_EQ      = "<<="
	SHR         = ">>"
	SHR_EQ      = ">>="
	AND_AND     = "&&"
	OR_OR       = "||"
	ARROW       = "<-"
	INC         = "++"
	DEC         = "--"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
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
	LT_EQ  = "<="
	GT     = ">"
	GT_EQ  = ">="
)

// Keywords maps keyword strings to their token types
var Keywords = map[string]TokenType{
	"func":      FUNCTION,
	"數":         SAIKA_FUNC,
	"let":       LET,
	"var":       VAR,
	"const":     CONST,
	"true":      TRUE,
	"false":     FALSE,
	"if":        IF,
	"else":      ELSE,
	"return":    RETURN,
	"package":   PACKAGE,
	"import":    IMPORT,
	"for":       FOR,
	"range":     RANGE,
	"break":     BREAK,
	"continue":  CONTINUE,
	"struct":    STRUCT,
	"interface": INTERFACE,
	"map":       MAP,
	"chan":      CHAN,
	"go":        GO,
	"select":    SELECT,
	"case":      CASE,
	"default":   DEFAULT,
	"switch":    SWITCH,
	"type":      TYPE,
}
