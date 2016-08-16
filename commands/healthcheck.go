package commands

import (
	"os"

	"github.com/spf13/cobra"
)

// HealthcheckCmd is the main (root) command for the CLI.
var HealthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "healthcheck monitors components of installable software",
	Long:  "healthcheck monitors installable software components as part of the premkit toolchain",

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := InitializeConfig(); err != nil {
			return err
		}

		cmd.Usage()
		return nil
	},
}

// Execute adds all child commands to the room command HealthcheckCmd and sets flags
func Execute() {
	AddCommands()

	if _, err := HealthcheckCmd.ExecuteC(); err != nil {
		os.Exit(-1)
	}
}

// AddCommands will add all child commands to the HealthcheckCmd
func AddCommands() {
	HealthcheckCmd.AddCommand(serverCmd)
}

// InitializeConfig initializes the config environment with defaults.
func InitializeConfig(subCmdVs ...*cobra.Command) error {
	return nil
}
