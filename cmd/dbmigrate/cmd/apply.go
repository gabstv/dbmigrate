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

var insertQDefaults = map[string]string{
	"sqlite3": "INSERT INTO db_migrations (migration_id,author_name,created_at) VALUES (?,?,?);",
	"mysql":   "INSERT INTO db_migrations (migration_id,author_name,created_at) VALUES (?,?,?);",
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
		os.Unsetenv("DBMSESS__INSERT_QUERY")
		os.Setenv("DBMSESS__INSERT_QUERY", insertQDefaults[drvname])
		if viper.IsSet("migrations.insert_query") {
			if v := viper.GetString("migrations.insert_query"); v != "" {
				os.Setenv("DBMSESS__INSERT_QUERY", v)
			}
		}

		rootp, err := getMigrationsRootPath()
		if err != nil {
			return err
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
		napplied := 0
		for _, v := range newmigs {
			if verboseMode {
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "START:", v.Name)
			}
			os.Unsetenv("DBMSESS__UUID")
			os.Unsetenv("DBMSESS__AUTHOR")
			os.Unsetenv("DBMSESS__CREATED")
			os.Setenv("DBMSESS__UUID", v.UUID)
			os.Setenv("DBMSESS__AUTHOR", v.Author)
			os.Setenv("DBMSESS__CREATED", v.T.Format("2006-02-01 15:04:05"))
			if err := dbmigrate.Apply(v, verboseMode); err != nil {
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "FAILED:", v.Name)
				fmt.Println(err.Error())
				fmt.Printf("%v of %v migrations applied.\n", napplied, len(newmigs))
				os.Exit(2)
			}
			napplied++
			if verboseMode {
				fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "SUCCESS:", v.Name)
			}
		}
		os.Unsetenv("DBMSESS__UUID")
		os.Unsetenv("DBMSESS__AUTHOR")
		os.Unsetenv("DBMSESS__CREATED")
		fmt.Printf("%v of %v migrations applied.\n", napplied, len(newmigs))
	},
}
