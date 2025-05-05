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
	return fmt.Sprintf("import \"%s\"", stmt.Path)
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
		params = append(params, p.Value)
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	// Generate return type if any
	if stmt.ReturnType != nil {
		out.WriteString(" ")
		out.WriteString(stmt.ReturnType.Value)
	}

	// Generate function body
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
	case *ast.StringLiteral:
		return fmt.Sprintf("\"%s\"", expr.Value)
	default:
		return ""
	}
}
