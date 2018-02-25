package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/bitte-ein-bit/songbeamer-helper/aws"
	"github.com/bitte-ein-bit/songbeamer-helper/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(cmdChurchSongID)
}

var cmdChurchSongID = &cobra.Command{
	Use:   "churchsongid",
	Short: "Set ChurchSongID in all songs",
	Long:  `Set ChurchSongID in all songs.`,
	Run: func(cmd *cobra.Command, args []string) {
		processSongs()
	},
}

func processSongs() {
	path := viper.GetString("songspath")
	log.Printf("Reading songs from %v", path)
	files, err := ioutil.ReadDir(path)
	util.CheckForError(err)

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sng" {
			fullpath := filepath.Join(path, file.Name())
			processSong(fullpath)
		}
	}
}

func processSong(filename string) {
	found := false
	lines, err := util.File2lines(filename)
	util.CheckForError(err)

	for _, line := range lines {
		if line == "---" {
			break
		}
		header := strings.Split(line, "=")
		if strings.ToLower(header[0]) == "#churchsongid" {
			if header[1] != "1" {
				found = true
				break
			}
		}
	}

	if !found {
		log.Printf("Adding ChurchSongID to %v", filename)
		line := fmt.Sprintf("#ChurchSongID=%v\n", getNextChurchSongID())
		err := util.InsertStringToFile(filename, line, 1)
		util.CheckForError(err)
	}
}

func getNextChurchSongID() string {
	item := aws.GetDynamoDBNumericItem("last_used_ChurchSongID")
	newID := item.Value + 1
	aws.UpdateDynamoDBItem("last_used_ChurchSongID", newID)
	return fmt.Sprintf("LKG%04d", newID)
}
