/*
Copyright © 2025 <asciifaceman>
*/
package cmd

import (
	"errors"

	"github.com/asciifaceman/gong/gong"
	"github.com/asciifaceman/hobocode"
	"github.com/spf13/cobra"
)

var (
	addFiles   []string                  // List of files to add to the Gong bundle
	addFileIDs []string                  // ID of the file to add to the Gong bundle
	fileMap    = make(map[string]string) // Map to hold file paths and their corresponding IDs
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a single file to the provided Gong bundle with an ID",
	Example: `gong add --gong bundle.gong --file myfile.txt --id myfileID
gong add -g bundle.gong -i file1,file2 -f file1.txt,file2.txt`,
	Long: `Add a single file with an ID to the given gong bundle.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(addFileIDs) == 0 {
			return errors.New("at least one file ID is required to add to the Gong bundle")
		}
		if len(addFiles) == 0 {
			return errors.New("at least one file is required to add to the Gong bundle")
		}
		if len(addFiles) != len(addFileIDs) {
			return errors.New("the number of files must match the number of IDs")
		}
		if gongFile == "" {
			return errors.New("the --gong flag is required to specify the Gong bundle file")
		}
		suffix := gongFile[len(gongFile)-5:]
		if suffix != ".gong" {
			return errors.New("not a valid .gong file: " + gongFile)
		}

		// check that addFileIDs are unique
		for i, id := range addFileIDs {
			if _, exists := fileMap[id]; exists {
				return errors.New("duplicate file ID found: " + id)
			}
			fileMap[id] = addFiles[i]
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		hobocode.Infof("Adding %d files to Gong bundle: %s", len(fileMap), gongFile)

		a, err := gong.BuildAssetEntries(fileMap)
		if err != nil {
			hobocode.Errorf("Failed to build asset entries: %v", err)
			return
		}
		size := gong.AssignAssetOffsets(a)
		hobocode.Infof("Total size of asset entries: %d bytes", size)

		if err := gong.WriteGongFile(gongFile, a); err != nil {
			hobocode.Errorf("Failed to write Gong file: %v", err)
			return
		}
		hobocode.Infof("Successfully added files to Gong bundle: %s", gongFile)

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	//addCmd.Flags().StringSliceVarP(&addFiles, "files", "f", addFiles, "List of files to add to the Gong bundle")
	addCmd.Flags().StringSliceVarP(&addFiles, "file", "f", addFiles, "File to add to the Gong bundle")
	addCmd.Flags().StringSliceVarP(&addFileIDs, "id", "i", addFileIDs, "ID of the file to add to the Gong bundle")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
