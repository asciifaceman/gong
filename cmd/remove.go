/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/asciifaceman/gong/gong"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		d := &gong.DirectoryEntryHeader{
			ID:          []byte("exampleID"),
			FNAME:       []byte("exampleFile"),
			FileType:    []byte("txt"),
			Offset:      1234,
			Size:        5678,
			Compression: 0,
		}
		encoded := d.Encode()
		spew.Dump(encoded)

		f := afero.NewOsFs()
		file, err := f.OpenFile("test.gong", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		if _, err := file.Write(encoded); err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
