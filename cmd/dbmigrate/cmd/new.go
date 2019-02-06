package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var newCmdFileType string

func init() {
	rootCmd.AddCommand(newCmd)
	//
	newCmd.Flags().StringVarP(&newCmdFileType, "type", "t", "", "new migration file type {sql,go} (default: sql)")
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new migration file",
	Long: `Create a new database migration file.
  
Please bear in mind that the contents of the database migration file will be
executed inside a single transaction. You may want to split bigger files
to lighten the load.`,
	Example: "new -t go migration_subject",
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if newCmdFileType == "" {
			newCmdFileType = "sql"
		}
		fmt.Println("new migration file type:", newCmdFileType)
		fmt.Println("TODO: create migration file")
		// step 1 - get the path to the migrations folder
	},
}
