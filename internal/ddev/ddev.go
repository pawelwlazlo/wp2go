package ddev

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Setup runs ddev config and ddev start in projectDir with the given project name.
func Setup(projectDir, projectName string) error {
	config := exec.Command("ddev", "config", "--project-type=wordpress", "--project-name="+projectName, "--docroot=.")
	config.Dir = projectDir
	config.Stdout = os.Stdout
	config.Stderr = os.Stderr
	if err := config.Run(); err != nil {
		return fmt.Errorf("ddev config: %w", err)
	}

	start := exec.Command("ddev", "start")
	start.Dir = projectDir
	start.Stdout = os.Stdout
	start.Stderr = os.Stderr
	if err := start.Run(); err != nil {
		return fmt.Errorf("ddev start: %w", err)
	}
	return nil
}

// ImportDB runs ddev import-db with the given SQL file path.
func ImportDB(projectDir, sqlPath string) error {
	absPath, err := filepath.Abs(sqlPath)
	if err != nil {
		return err
	}
	cmd := exec.Command("ddev", "import-db", "--file="+absPath)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ddev import-db: %w", err)
	}
	return nil
}
