package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"wp2go/internal/restore"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	archivePath string
	dbPath      string
	siteName    string
	outputPath  string
	force       bool
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore WordPress from archive and SQL dump to a DDEV project",
	Long: `wp2go extracts a WordPress .tar.gz archive and SQL dump, creates a DDEV
WordPress project, replaces the old site URL in the database with the DDEV URL,
and imports the database.`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		initConfig()
		return applyConfig(cmd)
	},
	RunE: run,
}

func init() {
	restoreCmd.Flags().StringVarP(&archivePath, "archive", "a", "", "path to WordPress .tar.gz archive")
	restoreCmd.Flags().StringVarP(&dbPath, "db", "d", "", "path to SQL database dump")
	restoreCmd.Flags().StringVarP(&siteName, "siteName", "s", "", "site name (e.g. example → https://example.ddev.site)")
	restoreCmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "directory where to create the project (default: current directory)")
	restoreCmd.Flags().BoolVar(&useUI, "ui", false, "use interactive prompts for missing options")
	restoreCmd.Flags().BoolVarP(&force, "force", "f", false, "force overwrite if project directory already exists")
}

func run(_ *cobra.Command, _ []string) error {
	if useUI {
		if err := promptMissing(); err != nil {
			return err
		}
	}

	if archivePath == "" {
		return errors.New("archive path is required (--archive / -a)")
	}
	if dbPath == "" {
		return errors.New("database path is required (--db / -d)")
	}
	if siteName == "" {
		return errors.New("siteName is required (--siteName / -s)")
	}

	projectDir := outputPath
	if projectDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}
		projectDir = cwd
	} else {
		projectDir = expandPath(projectDir)
	}

	cfg := restore.Config{
		ArchivePath: archivePath,
		DBPath:      dbPath,
		Sitename:    siteName,
		OutputPath:  projectDir,
		Force:       force,
	}
	if useUI {
		if err := maybeSaveConfig(cfg); err != nil {
			return err
		}
	}
	return restore.Run(cfg)
}

func promptMissing() error {
	if archivePath == "" {
		v, err := selectFileOrPath("Archive path (.tar.gz)", []string{"*.tar.gz"})
		if err != nil {
			return fmt.Errorf("archive prompt: %w", err)
		}
		archivePath = strings.TrimSpace(v)
	}
	if dbPath == "" {
		v, err := selectFileOrPath("SQL database file", []string{"*.sql", "*.sql.gz"})
		if err != nil {
			return fmt.Errorf("db prompt: %w", err)
		}
		dbPath = strings.TrimSpace(v)
	}
	if siteName == "" {
		p := promptui.Prompt{
			Label:   "Site name",
			Default: "example",
		}
		v, err := p.Run()
		if err != nil {
			return fmt.Errorf("siteName prompt: %w", err)
		}
		siteName = strings.TrimSpace(v)
	}
	if outputPath == "" {
		p := promptui.Prompt{
			Label:   "Output path (directory for project)",
			Default: ".",
		}
		v, err := p.Run()
		if err != nil {
			return fmt.Errorf("output-path prompt: %w", err)
		}
		v = strings.TrimSpace(v)
		if v != "" && v != "." {
			outputPath = v
		}
	}
	return nil
}

func selectFileOrPath(label string, patterns []string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	files, err := globFiles(cwd, patterns)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return promptPath(label)
	}
	options := append(files, "Enter path...")
	selectPrompt := promptui.Select{
		Label: label,
		Items: options,
		Size:  minimum(12, len(options)),
	}
	_, choice, err := selectPrompt.Run()
	if err != nil {
		return "", err
	}
	if choice == "Enter path..." {
		return promptPath(label)
	}
	return choice, nil
}

func promptPath(label string) (string, error) {
	p := promptui.Prompt{
		Label: label + " (enter path)",
	}
	return p.Run()
}

func globFiles(dir string, patterns []string) ([]string, error) {
	seen := map[string]struct{}{}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			seen[filepath.Base(match)] = struct{}{}
		}
	}
	if len(seen) == 0 {
		return nil, nil
	}
	files := make([]string, 0, len(seen))
	for name := range seen {
		files = append(files, name)
	}
	sort.Strings(files)
	return files, nil
}

func maybeSaveConfig(cfg restore.Config) error {
	p := promptui.Select{
		Label: "Save configuration to wp2go.yaml?",
		Items: []string{"No", "Yes"},
	}
	_, choice, err := p.Run()
	if err != nil {
		return err
	}
	if choice != "Yes" {
		return nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(cwd, "wp2go.yaml")
	data := formatConfigYAML(cfg)
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		return fmt.Errorf("write wp2go.yaml: %w", err)
	}
	return nil
}

func formatConfigYAML(cfg restore.Config) string {
	forceVal := "false"
	if cfg.Force {
		forceVal = "true"
	}

	return fmt.Sprintf(
		"defaults:\n"+
			"  force: %s\n"+
			"\n"+
			"hosts:\n"+
			"  %s:\n"+
			"    archive: %s\n"+
			"    db: %s\n"+
			"    sitename: %s\n",
		forceVal,
		cfg.Sitename,
		strconv.Quote(cfg.ArchivePath),
		strconv.Quote(cfg.DBPath),
		strconv.Quote(cfg.Sitename),
	)
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expandPath(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		if p == "~" {
			return home
		}
		return filepath.Join(home, p[2:])
	}
	return p
}
