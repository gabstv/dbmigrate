package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gabstv/dbmigrate/pkg/dbmigrate"

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
		tableexists, err := dbmigrate.MigrationTableExists()
		if err != nil {
			fmt.Println("error while checking if the db_migrations table exists:", err.Error())
			os.Exit(1)
		}
		if !tableexists {
			if !applyCmdConfirm {
				fmt.Println("The table db_migrations does not exist and will be created.")
				prompt := promptui.Prompt{
					IsConfirm: true,
					Label:     "Continue",
				}
				result, err := prompt.Run()
				if err != nil {
					fmt.Println("Error:", err.Error())
					os.Exit(1)
				}
				if result == "n" {
					os.Exit(2)
				}
			}
			if err := dbmigrate.CreateMigrationTable(); err != nil {
				fmt.Println("Error:", err.Error())
				os.Exit(1)
			}
		}
		newmigs, oldmigs, err := dbmigrate.ListMigrations(applyCmdRootPath)
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
		if verboseMode {
			fmt.Println("Installed Migrations:")
			for _, v := range oldmigs {
				fmt.Println(v.Name, v.T.Format("2006-02-01 15:04:05"))
			}
		}
		if len(newmigs) == 0 {
			fmt.Println("No migrations to apply.")
			os.Exit(0)
		}
		fmt.Printf("%v migrations to apply:\n", len(newmigs))
		for _, v := range newmigs {
			fmt.Println(v.Name, v.T.Format("2006-02-01 15:04:05"))
		}
		if !applyCmdConfirm {
			prompt := promptui.Prompt{
				IsConfirm: true,
				Label:     "Continue",
			}
			result, err := prompt.Run()
			if err != nil {
				fmt.Println("Error:", err.Error())
				os.Exit(1)
			}
			if result == "n" {
				os.Exit(2)
			}
		}
		for _, v := range newmigs {
			fmt.Println("TODO: mig this", v.Name)
			// TODO: run migration (go run for go files)
			if verboseMode {
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "START:", v.Name)
			}
			if err := dbmigrate.Apply(v, verboseMode); err != nil {
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "FAILED:", v.Name)
				fmt.Println(err.Error())
			}
		}
	},
}
