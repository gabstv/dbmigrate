package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var applyCmdDBName string

func init() {
	rootCmd.AddCommand(applyCmd)
	//
	applyCmd.Flags().StringVarP(&applyCmdDBName, "database", "d", "", "database to apply the migrations (default specified in the config file)")
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the migrations to a database.",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}
