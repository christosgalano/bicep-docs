package template

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// BuildBicepTemplate compiles a Bicep file into an ARM template, stored in the temporary directory.
// It checks for the 'bicep' or 'az' commands to perform the build.
// The function returns the path to the generated ARM template file, or an error if the build fails or the necessary commands are not found.
func BuildBicepTemplate(bicepFile string) (string, error) {
	// Validate file extension
	basename := filepath.Base(bicepFile)
	filename := strings.TrimSuffix(basename, filepath.Ext(basename))
	extension := filepath.Ext(basename)
	if extension != ".bicep" {
		return "", fmt.Errorf("file extension must be '.bicep'")
	}

	// ARM template file path
	tmp := os.TempDir()
	armFile := filepath.Join(tmp, fmt.Sprintf("%s_%s.json", filename, uuid.New().String()))

	// Build Bicep template into an ARM template
	// 1. If 'bicep' exists, run 'bicep build'
	// 2. If 'bicep' does not exist, check if 'az' exists
	// 3. If 'az' exists, run 'az bicep build'
	// 4. If 'az' does not exist, return an error
	var cmd *exec.Cmd
	switch {
	case commandExists("bicep"):
		cmd = exec.Command("bicep", "build", bicepFile, "--outfile", armFile)
	case commandExists("az"):
		cmd = exec.Command("az", "bicep", "build", "--file", bicepFile, "--outfile", armFile)
	default:
		return "", fmt.Errorf("error processing %s: neither 'bicep' nor 'az' commands were found", bicepFile)
	}

	// Run the command and handle any errors
	if err := runCommand(cmd); err != nil {
		return "", err
	}

	return armFile, nil
}

// commandExists checks if a command exists in the system's PATH.
// It returns true if the command exists, otherwise false.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// runCommand executes the given command and returns an error if the command fails.
// It captures the error message from stderr and returns it as an error.
func runCommand(cmd *exec.Cmd) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Extract the error message from stderr
		errorLines := strings.Split(stderr.String(), "\n")
		for _, line := range errorLines {
			if strings.Contains(line, "Error") {
				return errors.New(line)
			}
		}
		return fmt.Errorf("failed to run command: %w", err)
	}

	return nil
}
