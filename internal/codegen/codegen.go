package codegen

import (
	"fmt"
	"strings"

	"github.com/saika-m/saika-lang/internal/ast"
)

// Generator represents a code generator for Saika
type Generator struct {
	program        *ast.Program
	errors         []string
	indentLevel    int
	currentPackage string
}

// Error represents a code generation error
type Error struct {
	Message  string
	Position ast.Position
}

func (e Error) Error() string {
	if e.Position.File != "" {
		return fmt.Sprintf("%s:%d:%d: %s", e.Position.File, e.Position.Line, e.Position.Column, e.Message)
	}
	return fmt.Sprintf("line %d, column %d: %s", e.Position.Line, e.Position.Column, e.Message)
}

// New creates a new Generator
func New(program *ast.Program) *Generator {
	return &Generator{
		program:     program,
		errors:      []string{},
		indentLevel: 0,
	}
}

// Errors returns code generation errors
func (g *Generator) Errors() []string {
	return g.errors
}

// indent returns the current indentation string
func (g *Generator) indent() string {
	return strings.Repeat("\t", g.indentLevel)
}

// increaseIndent increases the indentation level
func (g *Generator) increaseIndent() {
	g.indentLevel++
}

// decreaseIndent decreases the indentation level
func (g *Generator) decreaseIndent() {
	if g.indentLevel > 0 {
		g.indentLevel--
	}
}

// Generate generates Go code from the AST
func (g *Generator) Generate() string {
	var out strings.Builder

	// Process all statements
	for _, stmt := range g.program.Statements {
		code := g.generateStatement(stmt)
		if code != "" {
			out.WriteString(code)
			out.WriteString("\n")
		}
	}

	return out.String()
}

// generateStatement generates code for a statement
func (g *Generator) generateStatement(stmt ast.Statement) string {
	switch stmt := stmt.(type) {
	case *ast.PackageStatement:
		return g.generatePackageStatement(stmt)
	case *ast.ImportStatement:
		return g.generateImportStatement(stmt)
	case *ast.FunctionStatement:
		return g.generateFunctionStatement(stmt)
	case *ast.VariableStatement:
		return g.generateVariableStatement(stmt)
	case *ast.ReturnStatement:
		return g.generateReturnStatement(stmt)
	case *ast.IfStatement:
		return g.generateIfStatement(stmt)
	case *ast.ForStatement:
		return g.generateForStatement(stmt)
	case *ast.RangeStatement:
		return g.generateRangeStatement(stmt)
	case *ast.ExpressionStatement:
		return g.generateExpressionStatement(stmt)
	case *ast.StructLiteral:
		return g.generateStructStatement(stmt)
	default:
		g.addError(fmt.Sprintf("unsupported statement type: %T", stmt), stmt.GetPosition())
		return ""
	}
}

// generatePackageStatement generates code for a package statement
func (g *Generator) generatePackageStatement(stmt *ast.PackageStatement) string {
	g.currentPackage = stmt.Name
	return fmt.Sprintf("package %s", stmt.Name)
}

// generateImportStatement generates code for an import statement
func (g *Generator) generateImportStatement(stmt *ast.ImportStatement) string {
	// Make sure the path is properly quoted
	path := stmt.Path
	if !strings.HasPrefix(path, "\"") && !strings.HasSuffix(path, "\"") {
		path = fmt.Sprintf("\"%s\"", path)
	}

	if stmt.Alias != "" {
		return fmt.Sprintf("import %s %s", stmt.Alias, path)
	}
	return fmt.Sprintf("import %s", path)
}

// generateFunctionStatement generates code for a function statement
func (g *Generator) generateFunctionStatement(stmt *ast.FunctionStatement) string {
	var out strings.Builder

	// Replace 數 with func
	out.WriteString("func ")
	out.WriteString(stmt.Name.Value)
	out.WriteString("(")

	// Generate parameters
	params := []string{}
	for _, p := range stmt.Parameters {
		paramStr := p.Name.Value
		if p.Type != nil {
			paramStr += " " + g.generateTypeExpression(p.Type)
		}
		params = append(params, paramStr)
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	// Generate return type if any
	if stmt.ReturnType != nil {
		out.WriteString(" ")
		out.WriteString(g.generateTypeExpression(stmt.ReturnType))
	}

	// Generate function body
	out.WriteString(" ")
	out.WriteString(g.generateBlockStatement(stmt.Body))

	return out.String()
}

// generateVariableStatement generates code for a variable declaration
func (g *Generator) generateVariableStatement(stmt *ast.VariableStatement) string {
	var out strings.Builder

	// Use const or var based on the statement type
	if stmt.Const {
		out.WriteString("const ")
	} else {
		out.WriteString("var ")
	}

	out.WriteString(stmt.Name.Value)

	// Add type if specified
	if stmt.Type != nil {
		out.WriteString(" ")
		out.WriteString(g.generateTypeExpression(stmt.Type))
	}

	// Add initialization if present
	if stmt.Value != nil {
		out.WriteString(" = ")
		out.WriteString(g.generateExpression(stmt.Value))
	}

	return out.String()
}

// generateReturnStatement generates code for a return statement
func (g *Generator) generateReturnStatement(stmt *ast.ReturnStatement) string {
	var out strings.Builder

	out.WriteString("return")

	if stmt.Value != nil {
		out.WriteString(" ")
		out.WriteString(g.generateExpression(stmt.Value))
	}

	return out.String()
}

// generateIfStatement generates code for an if statement
func (g *Generator) generateIfStatement(stmt *ast.IfStatement) string {
	var out strings.Builder

	out.WriteString("if ")
	out.WriteString(g.generateExpression(stmt.Condition))
	out.WriteString(" ")
	out.WriteString(g.generateBlockStatement(stmt.Consequence))

	if stmt.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(g.generateBlockStatement(stmt.Alternative))
	}

	return out.String()
}

// generateForStatement generates code for a for statement
func (g *Generator) generateForStatement(stmt *ast.ForStatement) string {
	var out strings.Builder

	out.WriteString("for ")

	if stmt.Init != nil {
		out.WriteString(g.generateStatement(stmt.Init))
	}

	out.WriteString("; ")

	if stmt.Condition != nil {
		out.WriteString(g.generateExpression(stmt.Condition))
	}

	out.WriteString("; ")

	if stmt.Post != nil {
		out.WriteString(g.generateStatement(stmt.Post))
	}

	out.WriteString(" ")
	out.WriteString(g.generateBlockStatement(stmt.Body))

	return out.String()
}

// generateRangeStatement generates code for a for-range statement
func (g *Generator) generateRangeStatement(stmt *ast.RangeStatement) string {
	var out strings.Builder

	out.WriteString("for ")

	if stmt.Key != nil {
		out.WriteString(stmt.Key.Value)

		if stmt.Value != nil {
			out.WriteString(", ")
			out.WriteString(stmt.Value.Value)
		}

		out.WriteString(" := ")
	}

	out.WriteString("range ")
	out.WriteString(g.generateExpression(stmt.Collection))
	out.WriteString(" ")
	out.WriteString(g.generateBlockStatement(stmt.Body))

	return out.String()
}

// generateStructStatement generates code for a struct declaration
func (g *Generator) generateStructStatement(stmt ast.Expression) string {
	// Type assert to get the actual StructLiteral
	structLit, ok := stmt.(*ast.StructLiteral)
	if !ok {
		g.addError(fmt.Sprintf("expected StructLiteral, got %T", stmt), stmt.GetPosition())
		return ""
	}

	var out strings.Builder

	if structLit.Name != nil {
		out.WriteString("type ")
		out.WriteString(structLit.Name.Value)
		out.WriteString(" ")
	}

	out.WriteString("struct {\n")

	g.increaseIndent()

	for name, typ := range structLit.Fields {
		out.WriteString(g.indent())
		out.WriteString(name)
		out.WriteString(" ")
		out.WriteString(g.generateTypeExpression(typ))
		out.WriteString("\n")
	}

	g.decreaseIndent()
	out.WriteString(g.indent())
	out.WriteString("}")

	return out.String()
}

// generateBlockStatement generates code for a block statement
func (g *Generator) generateBlockStatement(stmt *ast.BlockStatement) string {
	var out strings.Builder

	out.WriteString("{\n")

	g.increaseIndent()

	for _, s := range stmt.Statements {
		out.WriteString(g.indent())
		out.WriteString(g.generateStatement(s))
		out.WriteString("\n")
	}

	g.decreaseIndent()
	out.WriteString(g.indent())
	out.WriteString("}")

	return out.String()
}

// generateExpressionStatement generates code for an expression statement
func (g *Generator) generateExpressionStatement(stmt *ast.ExpressionStatement) string {
	return g.generateExpression(stmt.Expression)
}

// generateExpression generates code for an expression
func (g *Generator) generateExpression(expr ast.Expression) string {
	switch expr := expr.(type) {
	case *ast.Identifier:
		return expr.Value
	case *ast.IntegerLiteral:
		return fmt.Sprintf("%d", expr.Value)
	case *ast.FloatLiteral:
		return fmt.Sprintf("%f", expr.Value)
	case *ast.StringLiteral:
		return fmt.Sprintf("\"%s\"", expr.Value)
	case *ast.BooleanLiteral:
		if expr.Value {
			return "true"
		}
		return "false"
	case *ast.ArrayLiteral:
		return g.generateArrayLiteral(expr)
	case *ast.HashLiteral:
		return g.generateHashLiteral(expr)
	case *ast.UnaryExpression:
		return g.generateUnaryExpression(expr)
	case *ast.BinaryExpression:
		return g.generateBinaryExpression(expr)
	case *ast.AssignmentExpression:
		return g.generateAssignmentExpression(expr)
	case *ast.MemberExpression:
		return g.generateMemberExpression(expr)
	case *ast.IndexExpression:
		return g.generateIndexExpression(expr)
	case *ast.CallExpression:
		return g.generateCallExpression(expr)
	default:
		g.addError(fmt.Sprintf("unsupported expression type: %T", expr), expr.GetPosition())
		return ""
	}
}

// generateTypeExpression generates code for a type expression
func (g *Generator) generateTypeExpression(expr ast.TypeExpression) string {
	switch expr := expr.(type) {
	case *ast.IdentifierType:
		return expr.Name
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", g.generateTypeExpression(expr.ElementType))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", g.generateTypeExpression(expr.KeyType), g.generateTypeExpression(expr.ValueType))
	case *ast.StructLiteral:
		return g.generateStructStatement(expr)
	default:
		g.addError(fmt.Sprintf("unsupported type expression: %T", expr), expr.GetPosition())
		return ""
	}
}

// generateArrayLiteral generates code for an array literal
func (g *Generator) generateArrayLiteral(expr *ast.ArrayLiteral) string {
	var out strings.Builder

	elements := []string{}
	for _, el := range expr.Elements {
		elements = append(elements, g.generateExpression(el))
	}

	out.WriteString("[]")
	// TODO: Array type should be inferred or specified
	out.WriteString("interface{}")
	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

// generateHashLiteral generates code for a hash literal
func (g *Generator) generateHashLiteral(expr *ast.HashLiteral) string {
	var out strings.Builder

	// Use map[string]interface{} as default type
	// TODO: Type should be inferred or specified
	out.WriteString("map[string]interface{}{")

	pairs := []string{}
	for key, value := range expr.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", g.generateExpression(key), g.generateExpression(value)))
	}

	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// generateUnaryExpression generates code for a unary expression
func (g *Generator) generateUnaryExpression(expr *ast.UnaryExpression) string {
	if expr.Token.Type == ast.INC || expr.Token.Type == ast.DEC {
		// Postfix expression
		return fmt.Sprintf("%s%s", g.generateExpression(expr.Right), expr.Operator)
	}

	// Prefix expression
	return fmt.Sprintf("%s%s", expr.Operator, g.generateExpression(expr.Right))
}

// generateBinaryExpression generates code for a binary expression
func (g *Generator) generateBinaryExpression(expr *ast.BinaryExpression) string {
	return fmt.Sprintf("(%s %s %s)",
		g.generateExpression(expr.Left),
		expr.Operator,
		g.generateExpression(expr.Right))
}

// generateAssignmentExpression generates code for an assignment expression
func (g *Generator) generateAssignmentExpression(expr *ast.AssignmentExpression) string {
	return fmt.Sprintf("%s %s %s",
		g.generateExpression(expr.Left),
		expr.Operator,
		g.generateExpression(expr.Right))
}

// generateMemberExpression generates code for a member expression
func (g *Generator) generateMemberExpression(expr *ast.MemberExpression) string {
	return fmt.Sprintf("%s.%s",
		g.generateExpression(expr.Object),
		g.generateExpression(expr.Property))
}

// generateIndexExpression generates code for an index expression
func (g *Generator) generateIndexExpression(expr *ast.IndexExpression) string {
	return fmt.Sprintf("%s[%s]",
		g.generateExpression(expr.Left),
		g.generateExpression(expr.Index))
}

// generateCallExpression generates code for a call expression
func (g *Generator) generateCallExpression(expr *ast.CallExpression) string {
	args := []string{}
	for _, arg := range expr.Arguments {
		args = append(args, g.generateExpression(arg))
	}

	return fmt.Sprintf("%s(%s)",
		g.generateExpression(expr.Function),
		strings.Join(args, ", "))
}

// addError adds an error message
func (g *Generator) addError(message string, pos ast.Position) {
	if pos.File != "" {
		g.errors = append(g.errors, fmt.Sprintf("%s:%d:%d: %s", pos.File, pos.Line, pos.Column, message))
	} else {
		g.errors = append(g.errors, fmt.Sprintf("line %d, column %d: %s", pos.Line, pos.Column, message))
	}
}
