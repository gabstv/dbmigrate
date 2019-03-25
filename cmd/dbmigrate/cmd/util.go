package cmd

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func getMigrationsRootPath() (string, error) {
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
			return "", errors.Wrap(err, "get prompt 1")
		}
		if result == "n" {
			return "", fmt.Errorf("aborted")
		}
	}
	return rootp, nil
}
