package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.ConfigFileUsed() == "" {
			return fmt.Errorf("config file not found/specified")
		}
		if applyCmdDBName == "" {
			applyCmdDBName = viper.GetString("default_database")
		}
		if applyCmdDBName == "" {
			return fmt.Errorf("no database selected")
		}
		// get db credentials
		os.Unsetenv("DBMSESS__DATA_SOURCE_NAME")
		dbcs := getDBFromConfig(applyCmdDBName)
		if dbcs == "" {
			return fmt.Errorf("no connection string")
		}
		os.Setenv("DBMSESS__DATA_SOURCE_NAME", dbcs)
		os.Unsetenv("DBMSESS__DRIVER_NAME")
		drvname := viper.GetString("driver")
		if drvname == "" {
			return fmt.Errorf("empty driver name")
		}
		os.Setenv("DBMSESS__DRIVER_NAME", drvname)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		//viper.Get("default_database")
		fmt.Println("will compare with db", applyCmdDBName)
	},
}
