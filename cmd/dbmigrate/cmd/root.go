// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabstv/dbmigrate/internal/pkg/util"
	"github.com/gabstv/dbmigrate/pkg/dbmigrate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var projectRoot string
var migrationsRoot string
var verboseMode bool

//
var configDir string

//

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbmigrate",
	Short: "Create and run database migrations using Go or pure SQL",
	Long: `DBMigrate is a CLI that creates and runs 
database migrations using Go or pure SQL.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $GIT_ROOT/dbmigrate.toml)")
	rootCmd.PersistentFlags().BoolVarP(&verboseMode, "verbose", "v", false, "verbose mode")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		p2, _ := filepath.Abs(cfgFile)
		configDir, _ = filepath.Split(p2)
		fmt.Println("path is", configDir)
	} else {
		if gitroot, err := util.GitRoot(); err == nil {
			viper.AddConfigPath(gitroot)
			viper.SetConfigName("dbmigrate")
			if !filepath.IsAbs(gitroot) {
				p2, _ := filepath.Abs(gitroot)
				configDir = p2
			} else {
				configDir = gitroot
			}
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verboseMode {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		if verboseMode {
			fmt.Println("not using config!", err.Error())
		}
	}
	// parsepaths
	if configDir != "" {
		if viper.IsSet("migrations.root") {
			viper.Set("migrations.root", rpath(configDir, viper.GetString("migrations.root")))
		}
	}
}

func rpath(base, v string) string {
	if !strings.HasPrefix(v, "./") {
		return v
	}
	return filepath.Join(base, v[1:])
}

func getcs(cs, cpath interface{}) string {
	css := ""
	cpaths := ""
	if cs != nil {
		css = cs.(string)
	}
	if cpath != nil {
		cpaths = cpath.(string)
	}
	if css != "" {
		return rpath(configDir, css)
	}
	if ff, err := ioutil.ReadFile(cpaths); err != nil {
		fmt.Println("Could not read", cpaths, "->", err.Error())
	} else {
		css = string(ff)
		return rpath(configDir, css)
	}
	return ""
}

func getDBFromConfig(name string) (connectionString string) {
	dbs := make([]dbmigrate.Database, 0)
	if err := viper.UnmarshalKey("databases", &dbs); err == nil {
		for _, v := range dbs {
			if v.Name == name {
				return getcs(v.CS, v.File)
			}
		}
	} else {
		fmt.Println("ERR", err.Error())
	}
	if verboseMode {
		fmt.Printf("databases.%v not set\n", name)
	}
	return
}
