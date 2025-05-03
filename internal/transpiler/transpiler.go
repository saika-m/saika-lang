package transpiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/saika-m/saika-lang/internal/codegen"
	"github.com/saika-m/saika-lang/internal/lexer"
	"github.com/saika-m/saika-lang/internal/parser"
)

// Transpiler represents a Saika to Go transpiler
type Transpiler struct {
	// Configuration options could be added here
}

// New creates a new Transpiler
func New() *Transpiler {
	return &Transpiler{}
}

// TranspileFile transpiles a Saika file to Go code
func (t *Transpiler) TranspileFile(saikaFilePath string) (string, error) {
	// Read the Saika file
	saikaCode, err := ioutil.ReadFile(saikaFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read Saika file: %v", err)
	}

	// Transpile the code
	goCode, err := t.Transpile(string(saikaCode))
	if err != nil {
		return "", fmt.Errorf("failed to transpile Saika code: %v", err)
	}

	return goCode, nil
}

// Transpile transpiles Saika code to Go code
func (t *Transpiler) Transpile(saikaCode string) (string, error) {
	// Create a lexer
	l := lexer.New(saikaCode)

	// Create a parser
	p := parser.New(l)

	// Parse the program
	program := p.ParseProgram()

	// Check for parser errors
	if len(p.Errors()) > 0 {
		return "", fmt.Errorf("parser errors: %v", p.Errors())
	}

	// Generate Go code
	g := codegen.New(program)
	goCode := g.Generate()

	return goCode, nil
}

// CreateTempGoFile creates a temporary Go file with the given code
func (t *Transpiler) CreateTempGoFile(goCode string) (string, string, error) {
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "saika-temp")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Create a temporary Go file
	tempFile := filepath.Join(tempDir, "temp.go")
	if err := ioutil.WriteFile(tempFile, []byte(goCode), 0644); err != nil {
		os.RemoveAll(tempDir)
		return "", "", fmt.Errorf("failed to write temp file: %v", err)
	}

	return tempFile, tempDir, nil
}
