/*
Copyright © 2025 <asciifaceman>
*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/asciifaceman/gong/gong/bundle"
	"github.com/asciifaceman/hobocode"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new empty Gong file",
	Example: `gong create --gong myfile.gong`,
	Long:    `Create and initialize an empty gong bundle.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if gongFile == "" {
			return fmt.Errorf("the --gong flag is required")
		}
		suffix := filepath.Ext(gongFile)
		if suffix != ".gong" {
			return fmt.Errorf("not a valid .gong file: %s", gongFile)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		hobocode.Infof("Creating Gong file: %s", gongFile)
		err := bundle.CreateEmptyGong(gongFile)
		if err != nil {
			hobocode.Errorf("Failed to create Gong file: %v", err)
			return
		}
		hobocode.Successf("Successfully created Gong file: %s", gongFile)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
