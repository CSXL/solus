package cmd

import (
	"fmt"
	"os"

	"github.com/CSXL/solus/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "solus",
	Short: "An AI-assisted project generator.",
	Long:  `Solus is an AI-assisted project generator by CSX Labs.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// TODO: Fix global logging configuration to work with Cobra.
	},
	Run: func(cmd *cobra.Command, args []string) {
		_, err := tui.Run()
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.solus.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
