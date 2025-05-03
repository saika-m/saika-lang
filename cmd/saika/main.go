package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/saika-m/saika-lang/internal/transpiler"
)

const (
	VERSION = "1.0.0"
)

func main() {
	// If no arguments are provided, print usage and exit
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Create a transpiler
	t := transpiler.New()

	// Handle commands
	switch command {
	case "build":
		if len(os.Args) < 3 {
			fmt.Println("Error: No input file specified")
			printUsage()
			os.Exit(1)
		}
		processFiles(t, os.Args[2:], false)
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Error: No input file specified")
			printUsage()
			os.Exit(1)
		}
		processFiles(t, os.Args[2:], true)
	case "version":
		fmt.Printf("Saika Transpiler v%s\n", VERSION)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Saika Language Transpiler")
	fmt.Println("=========================")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  saika build <file.saika>...  - Compile Saika file(s) to executable(s)")
	fmt.Println("  saika run <file.saika>       - Run a Saika file")
	fmt.Println("  saika version                - Print version information")
	fmt.Println("  saika help                   - Print this help message")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -o, --output <dir>           - Specify output directory")
	fmt.Println("  -v, --verbose                - Enable verbose output")
	fmt.Println("  -I, --include <dir>          - Add include path for imports")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  saika build examples/hello.saika")
	fmt.Println("  saika run examples/hello.saika")
	fmt.Println("  saika build -o build examples/*.saika")
}

// Parse command-line arguments and options
func parseArgs(args []string) ([]string, map[string]string, error) {
	files := []string{}
	options := map[string]string{}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Handle options
		if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-o", "--output":
				if i+1 >= len(args) {
					return nil, nil, fmt.Errorf("missing output directory")
				}
				options["output"] = args[i+1]
				i++ // Skip the next arg
			case "-v", "--verbose":
				options["verbose"] = "true"
			case "-I", "--include":
				if i+1 >= len(args) {
					return nil, nil, fmt.Errorf("missing include path")
				}
				includeDir := args[i+1]
				if existingPaths, ok := options["include"]; ok {
					options["include"] = existingPaths + "," + includeDir
				} else {
					options["include"] = includeDir
				}
				i++ // Skip the next arg
			default:
				return nil, nil, fmt.Errorf("unknown option: %s", arg)
			}
		} else {
			// Handle files
			if filepath.Ext(arg) == ".saika" {
				files = append(files, arg)
			} else {
				// Check if it's a directory or a glob pattern
				matches, err := filepath.Glob(arg)
				if err != nil {
					return nil, nil, fmt.Errorf("invalid pattern: %s", arg)
				}

				for _, match := range matches {
					if filepath.Ext(match) == ".saika" {
						files = append(files, match)
					}
				}
			}
		}
	}

	return files, options, nil
}

// Process Saika files
func processFiles(t *transpiler.Transpiler, args []string, run bool) {
	// Parse command-line arguments
	files, options, err := parseArgs(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Apply options to the transpiler
	if outputDir, ok := options["output"]; ok {
		t.SetOutputDir(outputDir)
	}

	if _, ok := options["verbose"]; ok {
		t.SetVerbose(true)
		fmt.Println("Verbose mode enabled")
	}

	if includePaths, ok := options["include"]; ok {
		for _, path := range strings.Split(includePaths, ",") {
			t.AddIncludePath(path)
			if t.Verbose {
				fmt.Printf("Added include path: %s\n", path)
			}
		}
	}

	// Process each file
	for _, file := range files {
		if run {
			runCommand(t, file)
		} else {
			buildCommand(t, file)
		}
	}
}

// buildCommand handles the 'build' command
func buildCommand(t *transpiler.Transpiler, saikaFile string) {
	if t.Verbose {
		fmt.Printf("Building %s...\n", saikaFile)
	}

	// Transpile the Saika file to Go
	goCode, err := t.TranspileFile(saikaFile)
	if err != nil {
		fmt.Printf("Error transpiling file %s: %v\n", saikaFile, err)
		os.Exit(1)
	}

	// Create a temporary Go file
	tempGoFile, tempDir, err := t.CreateTempGoFile(goCode)
	if err != nil {
		fmt.Printf("Error creating temporary file: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir) // Clean up temporary directory

	// Determine output file
	outputFile := t.OutputDir
	if outputFile == "" {
		outputFile = strings.TrimSuffix(saikaFile, ".saika")
	} else {
		baseName := filepath.Base(strings.TrimSuffix(saikaFile, ".saika"))
		outputFile = filepath.Join(outputFile, baseName)
	}

	// Compile the Go file
	cmd := exec.Command("go", "build", "-o", outputFile, tempGoFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error compiling file %s: %v\n", saikaFile, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully built: %s\n", outputFile)
}

// runCommand handles the 'run' command
func runCommand(t *transpiler.Transpiler, saikaFile string) {
	if t.Verbose {
		fmt.Printf("Running %s...\n", saikaFile)
	}

	// Transpile the Saika file to Go
	goCode, err := t.TranspileFile(saikaFile)
	if err != nil {
		fmt.Printf("Error transpiling file %s: %v\n", saikaFile, err)
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
		fmt.Printf("Error running file %s: %v\n", saikaFile, err)
		os.Exit(1)
	}
}
