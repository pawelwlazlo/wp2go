package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"wp2go/internal/restore"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	archivePath string
	dbPath      string
	sitename    string
	outputPath  string
	useUI       bool
)

var rootCmd = &cobra.Command{
	Use:   "wp2go",
	Short: "Restore WordPress from archive and SQL dump to a DDEV project",
	Long: `wp2go extracts a WordPress .tar.gz archive and SQL dump, creates a DDEV
WordPress project, replaces the old site URL in the database with the DDEV URL,
and imports the database.`,
	RunE: run,
}

func init() {
	rootCmd.Flags().StringVarP(&archivePath, "archive", "a", "", "path to WordPress .tar.gz archive")
	rootCmd.Flags().StringVarP(&dbPath, "db", "d", "", "path to SQL database dump")
	rootCmd.Flags().StringVarP(&sitename, "sitename", "s", "", "site name (e.g. delfinki → https://delfinki.ddev.site)")
	rootCmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "directory where to create the project (default: current directory)")
	rootCmd.Flags().BoolVar(&useUI, "ui", false, "use interactive prompts for missing options")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
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
	if sitename == "" {
		return errors.New("sitename is required (--sitename / -s)")
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
		Sitename:    sitename,
		OutputPath:  projectDir,
	}
	return restore.Run(cfg)
}

func promptMissing() error {
	if archivePath == "" {
		p := promptui.Prompt{
			Label: "Archive path (.tar.gz)",
		}
		v, err := p.Run()
		if err != nil {
			return fmt.Errorf("archive prompt: %w", err)
		}
		archivePath = strings.TrimSpace(v)
	}
	if dbPath == "" {
		p := promptui.Prompt{
			Label: "SQL database file",
		}
		v, err := p.Run()
		if err != nil {
			return fmt.Errorf("db prompt: %w", err)
		}
		dbPath = strings.TrimSpace(v)
	}
	if sitename == "" {
		p := promptui.Prompt{
			Label:   "Site name",
			Default: "delfinki",
		}
		v, err := p.Run()
		if err != nil {
			return fmt.Errorf("sitename prompt: %w", err)
		}
		sitename = strings.TrimSpace(v)
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
