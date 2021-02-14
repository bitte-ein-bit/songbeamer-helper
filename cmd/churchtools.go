package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/bitte-ein-bit/songbeamer-helper/songbeamer"
	"github.com/bitte-ein-bit/songbeamer-helper/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(cmdCTGet)
}

var cmdCTGet = &cobra.Command{
	Use:   "ct-push",
	Short: "Push songs to ChurchTools",
	Long:  `Sync changes to ChurchTools`,
	Run: func(cmd *cobra.Command, args []string) {
		uploadToChurchTools()
	},
}

func uploadToChurchTools() {
	// churchtools.Login()
	log.Println("syncing to ChurchTools")
	songs, err := churchtools.GetSongs()
	util.CheckForError(err)

	for _, song := range songs {
		song.Delete()
	}
	songs = nil

	songs = processSongbeamerSongs(songs)
	log.Println(songs)
	// id := churchtools.AddSong("test1", "jonathan", "", "", "", "", "")
	// id := 72
	// log.Println(id)
	// s := churchtools.Song{
	// 	ID: 75,
	// }
	// err = s.Delete()
	// util.CheckForError(err)
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

func filterSongs(songs map[string]churchtools.Song, ccliID string) (ret churchtools.Song) {
	for _, song := range songs {
		if song.CCLI == ccliID {
			return song
		}
	}
	return churchtools.Song{}
}

func processSongbeamerSongs(songs map[string]churchtools.Song) map[string]churchtools.Song {
	path := viper.GetString("songspath")
	duplicates := viper.GetString("duplicates")

	log.Printf("Reading songs from %v", path)
	files, err := ioutil.ReadDir(path)
	util.CheckForError(err)

	count := 0

	if songs == nil {
		songs = make(map[string]churchtools.Song)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sng" {
			fullpath := filepath.Join(path, file.Name())
			song := songbeamer.SongbeamerSong{}
			song.LoadFromFile(fullpath)
			log.Println(song)
			if song.ID != "" {
				// check if file is newer -> upload if newer
				continue
			}
			if song.CCLI != "" {
				ctSong := filterSongs(songs, song.CCLI)
				if ctSong.ID != 0 {
					log.Println(ctSong)
					song.AddID(ctSong.ID, ctSong.GetDefaultArrangement())
					song.Title = ctSong.Bezeichnung
					song.SetKeyOfArrangement(ctSong.GetDefaultArrangement())
					err := song.FixFilename()
					if err != nil {
						log.Printf("Cannot fix filename %s", err)
						song.MoveToDuplicates(duplicates)
					}
					continue
				}
				id := churchtools.AddSong(song.Title, song.Author, song.Copyright, song.CCLI, song.KeyOfArrangement, "", "")
				ctAPISong := churchtools.GetSong(id)
				log.Printf("Added new Song: %v", ctAPISong)
				ctSong = ctAPISong.ToSong()
				log.Printf("Converted to %v", ctSong)
				log.Printf("Default Arrangement: %v", ctSong.GetDefaultArrangement())
				songs[fmt.Sprintf("%d", ctSong.ID)] = ctSong
				song.AddID(ctSong.ID, ctSong.GetDefaultArrangement())
				log.Printf("SongID: %s", song.ID)
				err := song.FixFilename()
				if err != nil {
					log.Printf("Cannot fix filename %s", err)
				}
				ctAPIFile, err := churchtools.NewAPIFile(song.Filename)
				util.CheckForError(err)
				ctAPIFile.DomainID = ctAPISong.GetDefaultArrangement().ID
				ctAPIFile.DomainType = "song_arrangement"

				err = ctAPIFile.Save()
				util.CheckForError(err)
			}
			err := song.FixFilename()
			if err != nil {
				log.Printf("Cannot fix filename %s", err)
			}
			count++
			if count > 10 {
				break
			}
			//
			// break
		}
	}
	return songs
}

func processSongbeamerSong(filename string) {
	found := false
	lines, err := util.File2lines(filename)
	util.CheckForError(err)

	for _, line := range lines {
		if line == "---" {
			break
		}
		header := strings.Split(line, "=")
		if strings.ToLower(header[0]) == "#id" {
			if header[1] != "1" {
				found = true
				break
			}
		}
	}

	if !found {
		log.Printf("Adding ID to %v", filename)
		line := fmt.Sprintf("#ID=%v\n", getNextChurchSongID())
		err := util.InsertStringToFile(filename, line, 1)
		util.CheckForError(err)
	}
}
