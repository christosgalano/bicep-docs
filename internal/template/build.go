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

	// Build Bicep template using azure-cli
	var cmd *exec.Cmd
	cmd = exec.Command("az", "bicep", "build", "--file", bicepFile, "--outfile", armFile)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to build Bicep template: %w", err)
	}

	return armFile, nil
}
