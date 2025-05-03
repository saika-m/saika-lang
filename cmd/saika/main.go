package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/saika-m/saika/internal"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	saikaFile := os.Args[2]

	switch command {
	case "build":
		buildCommand(saikaFile)
	case "run":
		runCommand(saikaFile)
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

func buildCommand(saikaFile string) {
	// Transpile the Saika file to Go
	goFile, err := internal.SaveTranspiledFile(saikaFile)
	if err != nil {
		fmt.Printf("Error transpiling file: %v\n", err)
		os.Exit(1)
	}

	// Compile the Go file
	outputFile := strings.TrimSuffix(saikaFile, ".saika")
	cmd := exec.Command("go", "build", "-o", outputFile, goFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error compiling file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully built: %s\n", outputFile)
}

func runCommand(saikaFile string) {
	// Transpile the Saika file to Go
	goFile, err := internal.SaveTranspiledFile(saikaFile)
	if err != nil {
		fmt.Printf("Error transpiling file: %v\n", err)
		os.Exit(1)
	}

	// Run the Go file
	cmd := exec.Command("go", "run", goFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running file: %v\n", err)
		os.Exit(1)
	}
}
