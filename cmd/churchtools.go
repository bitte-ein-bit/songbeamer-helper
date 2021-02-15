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

const (
	noCCLISongCat = "2" // Songs ohne CCLI Nummer
	cCLISongCat   = "1" // inactive Songs
)

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
}

func filterSongs(songs map[string]churchtools.Song, filterField, search string) (ret churchtools.Song) {
	for _, song := range songs {
		switch filterField {
		case "CCLI":
			if song.CCLI == search {
				return song
			}
		case "Title":
			if song.Bezeichnung == search {
				return song
			}
		case "ID":
			if fmt.Sprintf("%d", song.ID) == search {
				return song
			}
		default:
			log.Fatal("No default filter")
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
				ctSong := filterSongs(songs, "CCLI", song.CCLI)
				if ctSong.ID != 0 {
					log.Println(ctSong)
					song.ID = fmt.Sprintf("%d", ctSong.ID)
					song.Title = ctSong.Bezeichnung
					song.Save()
					song.SetKeyOfArrangement(ctSong.GetDefaultArrangement())

					err := song.FixFilename()
					if err != nil {
						log.Printf("Cannot fix filename %s", err)
						song.MoveToDuplicates(duplicates)
					}

					continue
				}
				id := churchtools.AddSong(song.Title, song.Author, song.Copyright, song.CCLI, song.KeyOfArrangement, "", "", cCLISongCat)
				ctAPISong := churchtools.GetSong(id)
				ctSong = ctAPISong.ToSong()
				songs[fmt.Sprintf("%d", ctSong.ID)] = ctSong
				song.ID = fmt.Sprintf("%d", ctSong.ID)

				log.Printf("Added new Song: %v", ctAPISong)
				log.Printf("Converted to %v", ctSong)
				log.Printf("Default Arrangement: %v", ctSong.GetDefaultArrangement())
				log.Printf("SongID: %s", song.ID)
				song.Save()
				err := song.FixFilename()
				if err != nil {
					log.Printf("Cannot fix filename %s", err)
					song.MoveToDuplicates(duplicates)
				}

				ctAPIFile, err := churchtools.NewAPIFile(song.Filename)
				util.CheckForError(err)
				ctAPIFile.DomainID = ctAPISong.GetDefaultArrangement().ID
				ctAPIFile.DomainType = "song_arrangement"

				err = ctAPIFile.Save()
				util.CheckForError(err)
				continue
			}
			if song.Title != "" {
				ctSong := filterSongs(songs, "Title", song.Title)

				if ctSong.ID != 0 && ctSong.Author == song.Author {
					log.Println(ctSong)
					song.ID = fmt.Sprintf("%d", ctSong.ID)
					song.Save()
					song.SetKeyOfArrangement(ctSong.GetDefaultArrangement())
					// TODO arrangement erkennen
					err := song.FixFilename()
					if err != nil {
						log.Printf("Cannot fix filename %s", err)
						song.MoveToDuplicates(duplicates)
					}
					continue
				}

				id := churchtools.AddSong(song.Title, song.Author, song.Copyright, song.CCLI, song.KeyOfArrangement, "", "", noCCLISongCat)
				ctAPISong := churchtools.GetSong(id)
				ctSong = ctAPISong.ToSong()
				songs[fmt.Sprintf("%d", ctSong.ID)] = ctSong
				song.ID = fmt.Sprintf("%d", ctSong.ID)

				log.Printf("Added new Song: %v", ctAPISong)
				log.Printf("Converted to %v", ctSong)
				log.Printf("Default Arrangement: %v", ctSong.GetDefaultArrangement())
				log.Printf("SongID: %s", song.ID)
				song.Save()
				err := song.FixFilename()
				if err != nil {
					log.Printf("Cannot fix filename %s", err)
					song.MoveToDuplicates(duplicates)
				}

				ctAPIFile, err := churchtools.NewAPIFile(song.Filename)
				util.CheckForError(err)
				ctAPIFile.DomainID = ctAPISong.GetDefaultArrangement().ID
				ctAPIFile.DomainType = "song_arrangement"

				err = ctAPIFile.Save()
				util.CheckForError(err)
				continue
			}
			err := song.FixFilename()
			if err != nil {
				log.Printf("Cannot fix filename %s", err)
				song.MoveToDuplicates(duplicates)
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
