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

// BuildBicepTemplate builds a Bicep template into an ARM template and stores it in the OS temp directory.
//
// The path to the newly created ARM template is returned.
//
// Bicep or Azure CLI must be installed (https://learn.microsoft.com/en-us/azure/azure-resource-manager/bicep/install).
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

// commandExists checks if a command exists.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// runCommand runs a command and returns an error if it fails.
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
