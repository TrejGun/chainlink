package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// main manages the process of combining coverage data for Chainlink nodes.
// It identifies "go-coverage" directories within a given root directory,
// merges their data into a "merged" directory for each test, and then
// calculates the overall coverage percentage.
func main() {
	// Check if the user has provided an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run script.go <searchPattern>")
		os.Exit(1)
	}

	// First argument after the program name is the search pattern
	searchPattern := os.Args[1]

	// Glob pattern to find all 'merged' directories in artifact folders
	dirs, err := filepath.Glob(searchPattern)
	if err != nil {
		fmt.Printf("Failed to find directories: %v\n", err)
		os.Exit(1)
	}

	if len(dirs) == 0 {
		fmt.Println("No directories found.")
		return
	}

	// Join the directory paths for input
	dirInput := strings.Join(dirs, ",")

	// Ensure the merged directory exists
	mergedDir := filepath.Join(".covdata", "merged")
	if err := os.MkdirAll(mergedDir, 0755); err != nil {
		fmt.Printf("Failed to create merged directory %s: %v\n", mergedDir, err)
		os.Exit(1)
	}

	// Merge the coverage data from all chainlink nodes
	mergeCmd := exec.Command("go", "tool", "covdata", "merge", "-o", mergedDir, "-i="+dirInput)
	fmt.Printf("Merging coverage for all tests:\n%s\n", mergeCmd.String())
	output, err := mergeCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing merge command: %v, output: %s\n", err, output)
		os.Exit(1)
	}

	// Calculate coverage percentage in the merged directory
	coverageCmd := exec.Command("go", "tool", "covdata", "percent", "-i=.")
	coverageCmd.Dir = mergedDir
	coverageOutput, err := coverageCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error calculating coverage percentage: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Total coverage based on all tests:\n%s\n%s\n", coverageCmd.String(), string(coverageOutput))
}