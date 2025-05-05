package codegen

import (
	"fmt"
	"strings"

	"github.com/saika-m/saika-lang/internal/ast"
)

// Generator represents a code generator for Saika
type Generator struct {
	program *ast.Program
}

// New creates a new Generator
func New(program *ast.Program) *Generator {
	return &Generator{
		program: program,
	}
}

// Generate generates Go code from the AST
func (g *Generator) Generate() string {
	var out strings.Builder

	// Process all statements
	for _, stmt := range g.program.Statements {
		out.WriteString(g.generateStatement(stmt))
		out.WriteString("\n")
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
	case *ast.VarStatement:
		return g.generateVarStatement(stmt)
	case *ast.ConstStatement:
		return g.generateConstStatement(stmt)
	case *ast.ReturnStatement:
		return g.generateReturnStatement(stmt)
	case *ast.IfStatement:
		return g.generateIfStatement(stmt)
	case *ast.ForStatement:
		return g.generateForStatement(stmt)
	case *ast.ExpressionStatement:
		return g.generateExpressionStatement(stmt)
	default:
		return ""
	}
}

// generatePackageStatement generates code for a package statement
func (g *Generator) generatePackageStatement(stmt *ast.PackageStatement) string {
	return fmt.Sprintf("package %s", stmt.Name)
}

// generateImportStatement generates code for an import statement
func (g *Generator) generateImportStatement(stmt *ast.ImportStatement) string {
	// Make sure the path has quotes around it
	// The Path field might already contain quotes from the parser
	path := stmt.Path
	if !strings.HasPrefix(path, "\"") {
		path = "\"" + path + "\""
	}
	return fmt.Sprintf("import %s", path)
}

// translateTypeName translates a Chinese type name to its Go equivalent
func (g *Generator) translateTypeName(typeName string) string {
	switch typeName {
	case "整数":
		return "int"
	case "字符串":
		return "string"
	case "浮点":
		return "float64"
	case "布尔":
		return "bool"
	default:
		return typeName
	}
}

// generateVarStatement generates code for a variable statement
func (g *Generator) generateVarStatement(stmt *ast.VarStatement) string {
	return fmt.Sprintf("var %s = %s",
		stmt.Name.Value,
		g.generateExpression(stmt.Value))
}

// generateConstStatement generates code for a constant statement
func (g *Generator) generateConstStatement(stmt *ast.ConstStatement) string {
	return fmt.Sprintf("const %s = %s",
		stmt.Name.Value,
		g.generateExpression(stmt.Value))
}

// generateReturnStatement generates code for a return statement
func (g *Generator) generateReturnStatement(stmt *ast.ReturnStatement) string {
	if stmt.ReturnValue != nil {
		return fmt.Sprintf("return %s", g.generateExpression(stmt.ReturnValue))
	}
	return "return"
}

// generateFunctionStatement generates code for a function statement
func (g *Generator) generateFunctionStatement(stmt *ast.FunctionStatement) string {
	var out strings.Builder

	// Replace 數 with func
	out.WriteString("func ")

	// Special case for main function (入口 -> main)
	if stmt.Name.Value == "入口" {
		out.WriteString("main")
	} else {
		out.WriteString(stmt.Name.Value)
	}

	out.WriteString("(")

	// Generate parameters
	params := []string{}
	for _, p := range stmt.Parameters {
		if p.Type != nil {
			params = append(params, fmt.Sprintf("%s %s",
				p.Name.Value,
				g.translateTypeName(p.Type.Value)))
		} else {
			params = append(params, p.Name.Value)
		}
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	// Generate return type if any
	if stmt.ReturnType != nil {
		out.WriteString(" ")
		out.WriteString(g.translateTypeName(stmt.ReturnType.Value))
	}

	// Generate function body
	out.WriteString(" ")
	out.WriteString(g.generateBlockStatement(stmt.Body))

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

	// Special handling for variable declarations in the initializer
	if stmt.Init != nil {
		if varStmt, ok := stmt.Init.(*ast.VarStatement); ok {
			// Use short declaration (:=) syntax instead of var
			out.WriteString(fmt.Sprintf("%s := %s",
				varStmt.Name.Value,
				g.generateExpression(varStmt.Value)))
		} else {
			// For other statement types, generate normally
			out.WriteString(g.generateStatement(stmt.Init))
		}
	}

	out.WriteString("; ")

	if stmt.Condition != nil {
		out.WriteString(g.generateExpression(stmt.Condition))
	}

	out.WriteString("; ")

	if stmt.Update != nil {
		// Strip the trailing semicolon from the update statement
		updateStmt := g.generateStatement(stmt.Update)
		if strings.HasSuffix(updateStmt, ";") {
			updateStmt = updateStmt[:len(updateStmt)-1]
		}
		out.WriteString(updateStmt)
	}

	out.WriteString(" ")
	out.WriteString(g.generateBlockStatement(stmt.Body))

	return out.String()
}

// generateBlockStatement generates code for a block statement
func (g *Generator) generateBlockStatement(stmt *ast.BlockStatement) string {
	var out strings.Builder

	out.WriteString("{\n")

	for _, s := range stmt.Statements {
		out.WriteString(g.generateStatement(s))

		// Add semicolon for certain statement types
		switch s.(type) {
		case *ast.ExpressionStatement, *ast.VarStatement, *ast.ConstStatement:
			if !strings.HasSuffix(out.String(), ";") {
				out.WriteString(";")
			}
		}

		out.WriteString("\n")
	}

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
	case *ast.StringLiteral:
		return fmt.Sprintf("\"%s\"", expr.Value)
	case *ast.BooleanLiteral:
		if expr.Value {
			return "true"
		}
		return "false"
	case *ast.PrefixExpression:
		return fmt.Sprintf("%s%s",
			expr.Operator,
			g.generateExpression(expr.Right))
	case *ast.InfixExpression:
		// Special case for modulo operator (% -> %)
		operator := expr.Operator

		return fmt.Sprintf("%s %s %s",
			g.generateExpression(expr.Left),
			operator,
			g.generateExpression(expr.Right))
	case *ast.AssignExpression:
		return fmt.Sprintf("%s = %s",
			g.generateExpression(expr.Left),
			g.generateExpression(expr.Value))
	case *ast.MemberExpression:
		return fmt.Sprintf("%s.%s",
			g.generateExpression(expr.Object),
			g.generateExpression(expr.Property))
	case *ast.CallExpression:
		args := []string{}
		for _, arg := range expr.Arguments {
			args = append(args, g.generateExpression(arg))
		}
		return fmt.Sprintf("%s(%s)",
			g.generateExpression(expr.Function),
			strings.Join(args, ", "))
	default:
		return ""
	}
}
