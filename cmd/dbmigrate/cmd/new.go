package cmd

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/gabstv/dbmigrate/pkg/dbmigrate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmdFileType string
var newCmdMigrationName string

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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if newCmdFileType == "" {
			newCmdFileType = viper.GetString("migrations.default_type")
		}
		if newCmdFileType == "" {
			newCmdFileType = "sql"
		}
		if u, err := user.Current(); u != nil {
			if len(args) > 0 {
				newCmdMigrationName = time.Now().Format("2006_01_02T150405_") + u.Username + "_" + args[0]
			} else {
				newCmdMigrationName = time.Now().Format("2006_01_02T150405_") + u.Username
			}
		} else if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new migration file type:", newCmdFileType)
		// step 1 - get the path to the migrations folder
		rootp, err := getMigrationsRootPath()
		if err != nil {
			fmt.Println("ERROR GETTING MIGRATIONS ROOT FOLDER:", err.Error())
			os.Exit(1)
		}
		mig, err := dbmigrate.New(newCmdMigrationName, dbmigrate.FileType(newCmdFileType), rootp)
		if err != nil {
			fmt.Println("ERROR:", err.Error())
			os.Exit(2)
		}
		fmt.Println(mig)
	},
}
