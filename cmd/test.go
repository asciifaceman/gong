/*
Copyright © 2025 <asciifaceman>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	testAppendFile = ""
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//		f := afero.NewOsFs()
		//		if testAppendFile == "" {
		//			panic(errors.New("testAppendFile must be set"))
		//		}
		//
		//		dat, err := afero.ReadFile(f, testAppendFile)
		//		if err != nil {
		//			panic(fmt.Errorf("failed to read file %s: %w", testAppendFile, err))
		//		}
		//
		//		if len(dat) == 0 {
		//			panic(errors.New("file is empty"))
		//		}
		//
		//		a := assets.New(testAppendFile, dat)
		//		spew.Dump(a)
		//
		//		// Write the asset to a new file
		//		outputFile := "output_test_asset.bin"
		//		if err := afero.WriteFile(f, outputFile, a.Bytes(), 0644); err != nil {
		//			panic(fmt.Errorf("failed to write asset to file %s: %w", outputFile, err))
		//		}
		//		fmt.Printf("Asset written to %s\n", outputFile)

	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVarP(&testAppendFile, "file", "f", "", "File to append to the bundle")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
