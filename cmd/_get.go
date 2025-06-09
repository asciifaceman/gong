/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/asciifaceman/gong/gong"
	"github.com/asciifaceman/hobocode"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	getId string // ID of the asset to retrieve
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Extract an asset from the Gong bundle by ID but don't alter the Gong file",
	Long:    `Extract an asset from the Gong bundle by ID but don't alter the Gong file.`,
	Example: `gong get --gong bundle.gong --id myfileID`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if getId == "" {
			return fmt.Errorf("the --id flag is required to specify the asset ID")
		}
		if gongFile == "" {
			return fmt.Errorf("the --gong flag is required to specify the Gong bundle file")
		}
		suffix := gongFile[len(gongFile)-5:]
		if suffix != ".gong" {
			return fmt.Errorf("not a valid .gong file: %s", gongFile)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		f := afero.NewOsFs()
		gf, err := gong.LoadGongFile(f, gongFile)
		if err != nil {
			hobocode.Errorf("failed to load gong file %s: %v", gongFile, err)
			return
		}
		defer gf.File.Close()

		asset, err := gf.GetAsset(getId)
		if err != nil {
			hobocode.Errorf("failed to retrieve asset with ID %s: %v", getId, err)
			return
		}

		file, err := f.Create(asset.Filename)
		if err != nil {
			hobocode.Errorf("failed to create file %s: %v", asset.Filename, err)
			return
		}
		defer file.Close()
		if _, err := file.Write(asset.Content); err != nil {
			hobocode.Errorf("failed to write content to file %s: %v", asset.Filename, err)
			return
		}
		hobocode.Infof("Successfully extracted asset with ID %s to file %s", getId, asset.Filename)

		// hobocode.HeaderLeft("Startup")
		// hobocode.Infof("Retrieving asset with ID %s from Gong file %s", getId, gongFile)
		// f := afero.NewOsFs()
		// gf, err := gong.LoadGongfile(f, gongFile)
		// if err != nil {
		// 	hobocode.Errorf("failed to load gong file %s: %v", gongFile, err)
		// 	return
		// }
		// defer gf.File.Close()
		// asset, content, err := gf.GetAsset(getId)
		// if err != nil {
		// 	hobocode.Errorf("failed to retrieve asset with ID %s: %v", getId, err)
		// 	return
		// }

		// file, err := f.Create(asset.Filename)
		// if err != nil {
		// 	hobocode.Errorf("failed to create file %s: %v", asset.Filename, err)
		// 	return
		// }
		// defer file.Close()
		// if _, err := file.Write(content); err != nil {
		// 	hobocode.Errorf("failed to write content to file %s: %v", asset.Filename, err)
		// 	return
		// }
		// hobocode.Infof("Successfully extracted asset with ID %s to file %s", getId, asset.Filename)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&getId, "id", "i", getId, "ID of the asset to retrieve")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
