package restore

import (
	"fmt"
	"os"
	"path/filepath"

	"wp2go/internal/archive"
	"wp2go/internal/ddev"
	"wp2go/internal/domain"
	"wp2go/internal/sql"
	"wp2go/internal/wpconfig"
)

// Config holds input paths and options for the restore.
type Config struct {
	ArchivePath string
	DBPath      string
	Sitename    string
	OutputPath  string
	Force       bool
}

// Run performs the full restore: normalize sitename, create project dir,
// extract archive, patch wp-config for DDEV, configure DDEV, replace domain in SQL, import DB.
func Run(cfg Config) error {
	shortName, fqdn := domain.Normalize(cfg.Sitename)
	projectDir := filepath.Join(cfg.OutputPath, shortName)

	// Check if project directory already exists
	if _, err := os.Stat(projectDir); err == nil {
		if !cfg.Force {
			return fmt.Errorf("project directory already exists: %s\nUse --force to overwrite", projectDir)
		}
	}

	if err := archive.Extract(cfg.ArchivePath, projectDir); err != nil {
		return fmt.Errorf("extract archive: %w", err)
	}

	if err := wpconfig.PatchForDDEV(projectDir); err != nil {
		return fmt.Errorf("patch wp-config for DDEV: %w", err)
	}

	if err := ddev.Setup(projectDir, shortName); err != nil {
		return fmt.Errorf("ddev setup: %w", err)
	}

	sqlPath, err := sql.Prepare(cfg.DBPath, fqdn)
	if err != nil {
		return fmt.Errorf("prepare SQL: %w", err)
	}
	defer sql.Cleanup(sqlPath)

	if err := ddev.ImportDB(projectDir, sqlPath); err != nil {
		return fmt.Errorf("import DB: %w", err)
	}

	return nil
}
