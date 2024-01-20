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
	if commandExists("bicep") {
		cmd = exec.Command("bicep", "build", bicepFile, "--outfile", armFile)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to run 'bicep build': %w", err)
		}
	} else if commandExists("az") {
		cmd = exec.Command("az", "bicep", "build", "--file", bicepFile, "--outfile", armFile)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to run 'az bicep build': %w", err)
		}
	} else {
		return "", fmt.Errorf("neither 'bicep' nor 'az' commands were found")
	}
	return armFile, nil
}

// commandExists checks if a command exists.
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
