package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	useUI bool
)

var rootCmd = &cobra.Command{
	Use:   "wp2go",
	Short: "Commands to backup/restore WordPress environments",
	Long: `wp2go extracts a WordPress .tar.gz archive and SQL dump, creates a DDEV
WordPress project, replaces the old site URL in the database with the DDEV URL,
and imports the database.`,
	RunE: run,
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
