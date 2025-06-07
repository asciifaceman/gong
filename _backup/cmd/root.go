/*
Copyright © 2025 <asciifaceman>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	gongFile = ""
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gong",
	Short: "A tool for managing Gong files",
	Long: `Gong is a pure-go asset packing/management system for
bundling files into a single file that can either be extracted
or accessed from within a go project.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&gongFile, "gong", "g", "", "Path to the Gong file")
}

func initConfig() {
}
