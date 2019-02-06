package cmd

import (
	"os"
	"path/filepath"

	"github.com/gabstv/dbmigrate/internal/pkg/util"

	"github.com/spf13/cobra"
)

var initCmdRoot string

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&initCmdRoot, "path", "p", "", "project root path (default is Git's root, or the working directory if git is not available)")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new configuration file",
	Long:  "Create a new configuration file at the project's root folder",
	PreRun: func(cmd *cobra.Command, args []string) {
		if initCmdRoot == "" {
			// get git if available!
			if groot, err := util.GitRoot(); err == nil {
				initCmdRoot = groot
			} else {
				initCmdRoot, _ = os.Getwd()
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "debug" {
			programDebug(cmd)
			return
		}
		if initCmdRoot == "" {
			cmd.Println("root directory could not be obtainer")
			os.Exit(1)
		}
		migdir := "database/migrations/"
		err := util.NewConfig(filepath.Join(initCmdRoot, "dbmigrate.toml"), util.CfgTypeTOML, util.NewConfigDefaults{
			MigrationsPath: migdir,
		})
		if err != nil {
			cmd.Println(err.Error())
			os.Exit(1)
		}
		cmd.Println("created:", filepath.Join(initCmdRoot, "dbmigrate.toml"))
	},
}

func programDebug(cmd *cobra.Command) {
	cmd.Println("DEBUG")
	workingDir, _ := os.Getwd()
	cmd.Println("workdir:", workingDir)
	groot, err := util.GitRoot()
	if err != nil {
		cmd.Println("git root: ERROR (", err.Error(), ")")
	} else {
		cmd.Println("git root:", groot)
	}
}
