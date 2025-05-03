package transpiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/saika-m/saika-lang/internal/codegen"
	"github.com/saika-m/saika-lang/internal/lexer"
	"github.com/saika-m/saika-lang/internal/parser"
)

// Transpiler represents a Saika to Go transpiler
type Transpiler struct {
	// Configuration options
	Verbose      bool     // Enable verbose output
	OutputDir    string   // Output directory for generated files
	IncludePaths []string // Include paths for imports
}

// TranspileResult represents the result of a transpilation
type TranspileResult struct {
	GoCode     string   // Generated Go code
	Errors     []string // Errors during transpilation
	Warnings   []string // Warnings during transpilation
	SourceFile string   // Source file
	OutputFile string   // Output file
}

// New creates a new Transpiler
func New() *Transpiler {
	return &Transpiler{
		Verbose:      false,
		OutputDir:    "",
		IncludePaths: []string{},
	}
}

// SetVerbose sets the verbose flag
func (t *Transpiler) SetVerbose(verbose bool) {
	t.Verbose = verbose
}

// SetOutputDir sets the output directory
func (t *Transpiler) SetOutputDir(dir string) {
	t.OutputDir = dir
}

// AddIncludePath adds an include path
func (t *Transpiler) AddIncludePath(path string) {
	t.IncludePaths = append(t.IncludePaths, path)
}

// TranspileFile transpiles a Saika file to Go code
func (t *Transpiler) TranspileFile(saikaFilePath string) (string, error) {
	// Read the Saika file
	saikaCode, err := os.ReadFile(saikaFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read Saika file: %v", err)
	}

	// Transpile the code
	result, err := t.TranspileWithPath(string(saikaCode), saikaFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to transpile Saika code: %v", err)
	}

	return result.GoCode, nil
}

// Transpile transpiles Saika code to Go code
func (t *Transpiler) Transpile(saikaCode string) (string, error) {
	result, err := t.TranspileWithPath(saikaCode, "")
	if err != nil {
		return "", err
	}
	return result.GoCode, nil
}

// TranspileWithPath transpiles Saika code to Go code with file path information
func (t *Transpiler) TranspileWithPath(saikaCode string, filePath string) (*TranspileResult, error) {
	result := &TranspileResult{
		SourceFile: filePath,
		Errors:     []string{},
		Warnings:   []string{},
	}

	// Create a lexer
	l := lexer.NewWithFilename(saikaCode, filePath)

	// Create a parser
	p := parser.New(l)

	// Parse the program
	program := p.ParseProgram()

	// Check for parser errors
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			result.Errors = append(result.Errors, err)
		}
		return result, fmt.Errorf("parser errors: %s", strings.Join(p.Errors(), "\n"))
	}

	// Generate Go code
	g := codegen.New(program)
	goCode := g.Generate()

	// Check for code generation errors
	if len(g.Errors()) > 0 {
		for _, err := range g.Errors() {
			result.Errors = append(result.Errors, err)
		}
		return result, fmt.Errorf("code generation errors: %s", strings.Join(g.Errors(), "\n"))
	}

	result.GoCode = goCode

	// If an output directory is specified, determine the output file
	if t.OutputDir != "" && filePath != "" {
		baseName := filepath.Base(filePath)
		goFileName := strings.TrimSuffix(baseName, filepath.Ext(baseName)) + ".go"
		result.OutputFile = filepath.Join(t.OutputDir, goFileName)
	}

	return result, nil
}

// TranspileProject transpiles a directory of Saika files to Go code
func (t *Transpiler) TranspileProject(saikaDir string) ([]*TranspileResult, error) {
	results := []*TranspileResult{}

	// Create output directory if it doesn't exist
	if t.OutputDir != "" {
		if err := os.MkdirAll(t.OutputDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	// Walk the directory and transpile all .saika files
	err := filepath.Walk(saikaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .saika files
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".saika") {
			return nil
		}

		// Transpile the file
		result, err := t.transpileAndSave(path)
		if err != nil {
			// Continue processing other files on error
			if t.Verbose {
				fmt.Fprintf(os.Stderr, "Error transpiling %s: %v\n", path, err)
			}
		}

		results = append(results, result)
		return nil
	})

	if err != nil {
		return results, fmt.Errorf("error walking directory: %v", err)
	}

	return results, nil
}

// transpileAndSave transpiles a Saika file to Go code and saves it
func (t *Transpiler) transpileAndSave(saikaFilePath string) (*TranspileResult, error) {
	// Transpile the file
	result, err := t.transpileFile(saikaFilePath)
	if err != nil {
		return result, err
	}

	// Determine output file
	outputFile := t.getOutputFilePath(saikaFilePath)
	result.OutputFile = outputFile

	// Save the generated Go code
	if err := t.saveGoCode(result.GoCode, outputFile); err != nil {
		return result, fmt.Errorf("failed to save Go code: %v", err)
	}

	return result, nil
}

// transpileFile transpiles a Saika file to Go code
func (t *Transpiler) transpileFile(saikaFilePath string) (*TranspileResult, error) {
	// Read the Saika file
	saikaCode, err := os.ReadFile(saikaFilePath)
	if err != nil {
		return &TranspileResult{
			SourceFile: saikaFilePath,
			Errors:     []string{fmt.Sprintf("failed to read Saika file: %v", err)},
		}, fmt.Errorf("failed to read Saika file: %v", err)
	}

	// Transpile the code
	return t.TranspileWithPath(string(saikaCode), saikaFilePath)
}

// getOutputFilePath determines the output file path for a Saika file
func (t *Transpiler) getOutputFilePath(saikaFilePath string) string {
	baseDir := filepath.Dir(saikaFilePath)
	baseName := filepath.Base(saikaFilePath)
	goFileName := strings.TrimSuffix(baseName, filepath.Ext(baseName)) + ".go"

	if t.OutputDir != "" {
		// Preserve directory structure if output dir is specified
		relPath, err := filepath.Rel(baseDir, saikaFilePath)
		if err == nil {
			relDir := filepath.Dir(relPath)
			return filepath.Join(t.OutputDir, relDir, goFileName)
		}
		return filepath.Join(t.OutputDir, goFileName)
	}

	return filepath.Join(baseDir, goFileName)
}

// saveGoCode saves the generated Go code to a file
func (t *Transpiler) saveGoCode(goCode string, outputFile string) error {
	// Create directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Write the file
	if err := os.WriteFile(outputFile, []byte(goCode), 0644); err != nil {
		return fmt.Errorf("failed to write Go file: %v", err)
	}

	if t.Verbose {
		fmt.Printf("Generated: %s\n", outputFile)
	}

	return nil
}

// CreateTempGoFile creates a temporary Go file with the given code
func (t *Transpiler) CreateTempGoFile(goCode string) (string, string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "saika-temp")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Create a temporary Go file
	tempFile := filepath.Join(tempDir, "temp.go")
	if err := os.WriteFile(tempFile, []byte(goCode), 0644); err != nil {
		os.RemoveAll(tempDir)
		return "", "", fmt.Errorf("failed to write temp file: %v", err)
	}

	return tempFile, tempDir, nil
}

// GetVersion returns the version of the transpiler
func (t *Transpiler) GetVersion() string {
	return "1.0.0" // Update this version as needed
}
