package internal

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Transpile converts Saika code to Go code
func Transpile(saikaCode string) string {
	// Replace 數 with func
	goCode := strings.ReplaceAll(saikaCode, "數", "func")

	return goCode
}

// TranspileFile converts a Saika file to a Go file
func TranspileFile(saikaFilePath string) (string, string, error) {
	// Read the Saika file
	content, err := ioutil.ReadFile(saikaFilePath)
	if err != nil {
		return "", "", err
	}

	// Convert the content
	goCode := Transpile(string(content))

	// Determine the output file path
	dir := filepath.Dir(saikaFilePath)
	base := filepath.Base(saikaFilePath)
	goFilePath := filepath.Join(dir, strings.TrimSuffix(base, ".saika")+".go")

	return goCode, goFilePath, nil
}

// SaveTranspiledFile transpiles a Saika file and saves the Go file
func SaveTranspiledFile(saikaFilePath string) (string, error) {
	goCode, goFilePath, err := TranspileFile(saikaFilePath)
	if err != nil {
		return "", err
	}

	// Write the Go file
	err = ioutil.WriteFile(goFilePath, []byte(goCode), 0644)
	if err != nil {
		return "", err
	}

	return goFilePath, nil
}
