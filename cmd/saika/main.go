// cmd/saika/main.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/saika-m/saika-lang/internal/transpiler"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	saikaFile := os.Args[2]

	// Create a transpiler
	t := transpiler.New()

	switch command {
	case "build":
		buildCommand(t, saikaFile)
	case "run":
		runCommand(t, saikaFile)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  saika build <file.saika>  - Compile the Saika file to an executable")
	fmt.Println("  saika run <file.saika>    - Run the Saika file")
}

func buildCommand(t *transpiler.Transpiler, saikaFile string) {
	// Transpile the Saika file to Go
	goCode, err := t.TranspileFile(saikaFile)
	if err != nil {
		fmt.Printf("Error transpiling file: %v\n", err)
		os.Exit(1)
	}

	// Create a temporary Go file
	tempGoFile, tempDir, err := t.CreateTempGoFile(goCode)
	if err != nil {
		fmt.Printf("Error creating temporary file: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir) // Clean up temporary directory

	// Compile the Go file
	outputFile := strings.TrimSuffix(saikaFile, ".saika")
	cmd := exec.Command("go", "build", "-o", outputFile, tempGoFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error compiling file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully built: %s\n", outputFile)
}

func runCommand(t *transpiler.Transpiler, saikaFile string) {
	// Transpile the Saika file to Go
	goCode, err := t.TranspileFile(saikaFile)
	if err != nil {
		fmt.Printf("Error transpiling file: %v\n", err)
		os.Exit(1)
	}

	// Create a temporary Go file
	tempGoFile, tempDir, err := t.CreateTempGoFile(goCode)
	if err != nil {
		fmt.Printf("Error creating temporary file: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir) // Clean up temporary directory

	// Run the Go file
	cmd := exec.Command("go", "run", tempGoFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running file: %v\n", err)
		os.Exit(1)
	}
}
