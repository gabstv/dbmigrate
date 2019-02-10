package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var applyCmdDBName string
var applyCmdConfirm bool
var applyCmdRootPath string

func init() {
	rootCmd.AddCommand(applyCmd)
	//
	applyCmd.Flags().StringVarP(&applyCmdDBName, "database", "d", "", "database to apply the migrations (default specified in the config file)")
	applyCmd.Flags().BoolVarP(&applyCmdConfirm, "yes", "y", false, "if true, the apply command doesn't wait for confirmation")
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
		if err := os.Setenv("DBMSESS__DATA_SOURCE_NAME", dbcs); err != nil {
			return errors.Wrap(err, "set env DSN")
		}
		os.Unsetenv("DBMSESS__DRIVER_NAME")
		drvname := viper.GetString("driver")
		if drvname == "" {
			return fmt.Errorf("empty driver name")
		}
		if err := os.Setenv("DBMSESS__DRIVER_NAME", drvname); err != nil {
			return errors.Wrap(err, "set env DSN")
		}

		rootp := viper.GetString("migrations.root")
		if rootp == "" && !applyCmdConfirm {
			fmt.Println("The migrations root path is empty. The current directory will be used:")
			wdd, _ := os.Getwd()
			fmt.Println(wdd)
			prompt := promptui.Prompt{
				IsConfirm: true,
				Label:     "Continue",
			}
			result, err := prompt.Run()
			if err != nil {
				return errors.Wrap(err, "get prompt 1")
			}
			if result == "n" {
				return fmt.Errorf("aborted")
			}
		}
		applyCmdRootPath = rootp
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		//viper.Get("default_database")
		fmt.Println("will compare with db", applyCmdDBName)
	},
}
