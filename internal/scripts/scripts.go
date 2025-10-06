package scripts

import (
	"github.com/spf13/cobra"

	"mobile-backend-boilerplate/internal/scripts/add_admin"
)

var rootCmd = &cobra.Command{
	Use:   "scripts",
	Short: "Utility CLI for the backend project",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(add_admin.Command)
}
