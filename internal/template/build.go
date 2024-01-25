package template

import (
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
	cmd.Stderr = os.Stderr
	if commandExists("bicep") {
		cmd = exec.Command("bicep", "build", bicepFile, "--outfile", armFile)
	} else if commandExists("az") {
		cmd = exec.Command("az", "bicep", "build", "--file", bicepFile, "--outfile", armFile)
	} else {
		return "", fmt.Errorf("neither 'bicep' nor 'az' commands were found")
	}

	// Run the command and handle any errors
	if err := runCommand(cmd); err != nil {
		return "", fmt.Errorf("failed to run command: %w", err)
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
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%s: %s", err, exitError.Stderr)
		}
		return err
	}
	return nil
}
