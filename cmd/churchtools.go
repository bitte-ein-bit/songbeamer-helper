package cmd

import (
	"log"

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
	log.Println("listing songs")
	songs, err := churchtools.GetSongs()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(songs)
	// // id := churchtools.AddSong("test1", "jonathan", "", "", "", "", "")
	// id := 72
	// log.Println(id)
	// song := churchtools.GetSong(id)
	// arrangement := song.Arrangements[0]
	// arrangement.Duration = arrangement.Duration + 1
	// churchtools.EditArrangement(arrangement, song.ID)
	// for _, file := range arrangement.Files {
	// 	log.Println(file.Name)
	// 	file.LoadFromFile("songs/All day.sng")
	// 	file.Save()
	// }

	// // songfile, err := churchtools.NewAPIFile("songs/All day.sng")
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// songfile.DomainID = arrangement.ID
	// songfile.DomainType = "song_arrangement"
	// log.Println(songfile)
	// err = songfile.Save()
	// if err != nil {
	// 	log.Fatal(err)
	// }
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
