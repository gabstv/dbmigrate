package cmd

import (
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(utilCmd)
	//
	utilCmd.AddCommand(utilUUIDCmd)
}

var utilCmd = &cobra.Command{
	Use:   "util",
	Short: "Utilities",
	Long:  `Utilities command`,
}

var utilUUIDCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate a new UUID",
	Long:  `Generate a new UUID V4`,
	Run: func(cmd *cobra.Command, args []string) {
		id := uuid.New()
		print(id.String())
	},
}
