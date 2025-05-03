#!/bin/bash
# Script to run all Saika tests
# Location: tests/run_tests.sh

# Set colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Print banner
echo -e "${PURPLE}"
echo "===================================================="
echo "       SAIKA LANGUAGE TEST SUITE                     "
echo "===================================================="
echo -e "${NC}"

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Change to root directory to find saika executable
cd "$ROOT_DIR" || exit 1

# Check if saika executable exists
if ! command -v ./saika &> /dev/null; then
    echo -e "${RED}Error: saika executable not found in directory: $ROOT_DIR${NC}"
    echo "Please build the saika transpiler first with: go build -o saika ./cmd/saika"
    exit 1
fi

# Define directories to test (relative to tests directory)
TEST_DIRS=("syntax" "advanced" "edge_cases" "comparison")

# Counters for tests
TOTAL_FILES=0
PASSED_FILES=0
FAILED_FILES=0

# Function to run tests in a directory
run_tests_in_dir() {
    local dir="$SCRIPT_DIR/$1"
    
    echo -e "\n${BLUE}=== Running tests in ${dir} ===${NC}\n"
    
    # Find all .saika files in the directory
    local files=$(find "$dir" -name "*.saika" -type f 2>/dev/null)
    
    if [ -z "$files" ]; then
        echo -e "${YELLOW}No .saika files found in $dir${NC}"
        return
    fi
    
    # Run tests for each file
    for file in $files; do
        TOTAL_FILES=$((TOTAL_FILES + 1))
        
        echo -e "${CYAN}Testing: $file${NC}"
        
        # Run the file
        ./saika run "$file"
        
        # Check the exit code
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ Test passed: $file${NC}"
            PASSED_FILES=$((PASSED_FILES + 1))
        else
            echo -e "${RED}✗ Test failed: $file${NC}"
            FAILED_FILES=$((FAILED_FILES + 1))
        fi
        
        echo "--------------------------------------"
    done
}

# Run tests in each directory
for dir in "${TEST_DIRS[@]}"; do
    if [ -d "$SCRIPT_DIR/$dir" ]; then
        run_tests_in_dir "$dir"
    else
        echo -e "${YELLOW}Directory not found: $SCRIPT_DIR/$dir${NC}"
    fi
done

# Run any individual tests in the root test directory
individual_tests=$(find "$SCRIPT_DIR" -maxdepth 1 -name "*.saika" -type f)
if [ -n "$individual_tests" ]; then
    echo -e "\n${BLUE}=== Running tests in root test directory ===${NC}\n"
    
    for file in $individual_tests; do
        TOTAL_FILES=$((TOTAL_FILES + 1))
        
        echo -e "${CYAN}Testing: $file${NC}"
        
        # Run the file
        ./saika run "$file"
        
        # Check the exit code
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ Test passed: $file${NC}"
            PASSED_FILES=$((PASSED_FILES + 1))
        else
            echo -e "${RED}✗ Test failed: $file${NC}"
            FAILED_FILES=$((FAILED_FILES + 1))
        fi
        
        echo "--------------------------------------"
    done
fi

# Print summary
echo -e "\n${PURPLE}===================================================="
echo "                TEST SUMMARY                        "
echo -e "====================================================${NC}"
echo -e "Total files tested: ${TOTAL_FILES}"
echo -e "${GREEN}Passed: ${PASSED_FILES}${NC}"
echo -e "${RED}Failed: ${FAILED_FILES}${NC}"

# Set exit code based on test results
if [ $FAILED_FILES -gt 0 ]; then
    echo -e "\n${RED}Some tests failed!${NC}"
    exit 1
else
    echo -e "\n${GREEN}All tests passed!${NC}"
    exit 0
fi