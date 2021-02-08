package cmd

import (
	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(cmdCTGet)
}

var cmdCTGet = &cobra.Command{
	Use:   "ct-get",
	Short: "Get songs from ChurchTools",
	Long:  `Set ChurchSongID in all songs.`,
	Run: func(cmd *cobra.Command, args []string) {
		listSongs()
	},
}

func listSongs() {
	// churchtools.Login()
	// churchtools.Songs()
	// churchtools.AddSong("test1", "jonathan", "", "", "", "", "")
	churchtools.AddSongFile(73, "songs/All day.sng")
	// path := viper.GetString("songspath")
	// log.Printf("Reading songs from %v", path)
	// files, err := ioutil.ReadDir(path)
	// util.CheckForError(err)

	// for _, file := range files {
	// 	if filepath.Ext(file.Name()) == ".sng" {
	// 		fullpath := filepath.Join(path, file.Name())
	// 		processSong(fullpath)
	// 	}
	// }
}
