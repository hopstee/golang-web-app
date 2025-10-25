package scripts

import (
	"github.com/spf13/cobra"

	addadmin "mobile-backend-boilerplate/internal/scripts/add_admin"
	generateschemas "mobile-backend-boilerplate/internal/scripts/generate_schema"
)

var rootCmd = &cobra.Command{
	Use:   "scripts",
	Short: "Utility CLI for the backend project",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(addadmin.Command)
	rootCmd.AddCommand(generateschemas.Command)
}
