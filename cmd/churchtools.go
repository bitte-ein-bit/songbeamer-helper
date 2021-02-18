package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

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

	// for _, song := range songs {
	// 	song.Delete()
	// }
	// songs = nil

	songs = processSongbeamerSongs(songs)
	log.Println("These are all songs known at the end of processing:")
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
	log.Println("done")
}

func filterSongs(songs map[string]churchtools.Song, filterField string, search interface{}) (ret churchtools.Song) {
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
			if song.ID == search {
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

	// count := 0

	if songs == nil {
		songs = make(map[string]churchtools.Song)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sng" {
			fullpath := filepath.Join(path, file.Name())
			song := songbeamer.Song{}
			song.LoadFromFile(fullpath)
			log.Printf("Working on %s", song.Title)
			if song.ID != "" {
				// Song has an ID, so it should exists in ChurchTools
				a := song.ExtractArrangementFromFilename()
				if a == "" {
					a = "Standard-Arrangement"
				}
				log.Printf("Song has CT ID of %d, arrangement by Filename is %s, arrangement by ID is %s", song.ChurchToolsID, a, song.ChurchToolsArrangement)

				ctSong := filterSongs(songs, "ID", song.ChurchToolsID)
				if ctSong.ID == 0 {
					song.ID = ""
					song.Save()
					song.FixFilename()
					if err != nil {
						log.Printf("Cannot fix filename %s", err)
						song.MoveToDuplicates(duplicates)
					}
					log.Printf("CT Song ID can't be found, resetting ID and skip")
					continue
				}
				arrangementID := 0
				var arrangement churchtools.SongArrangement
				for _, arrange := range ctSong.Arrangements {
					if arrange.Bezeichnung == a {
						log.Println("Arrangement found, checking if newer")
						arrangementID = arrange.ID
						arrangement = arrange
					}
				}
				if arrangementID == 0 {
					log.Printf("Seems song was renamed to new arrangement: %s", a)
					arrangementID, err = ctSong.AddArrangement(a)
					util.CheckForError(err)
					arrangement = churchtools.SongArrangement{
						ID:          arrangementID,
						Bezeichnung: a,
					}
					err = song.UploadToArrangement(arrangement, duplicates)
					if err != nil {
						log.Fatalf("Error Uploading: %v", err)
					}
					continue
				}
				newer := true
				APIFileToUpdate := *churchtools.NewSongAPIFile(song.Filename, arrangement.ID)

				for _, file := range arrangement.Files {
					log.Printf("Checking %s aka %s", file.Bezeichnung, file.Filename)
					if file.Bezeichnung != song.GetFilenameWithoutArrangement() {
						log.Printf("no match: %s != %s", file.Bezeichnung, song.GetFilenameWithoutArrangement())
						continue
					}
					ctDate, _ := file.GetModificationDate()
					sngDate, _ := song.GetModificationDate()
					if sngDate.After(ctDate) {
						log.Printf("CT is older: %v < %v", ctDate, sngDate)
						APIFileToUpdate = file.ToAPIFile()
						continue
					}
					log.Printf("CT is newer: %v >= %v", ctDate, sngDate)
					newer = false
				}
				if !newer {
					continue
				}
				log.Println("File needs to be updated on CT")
				APIFileToUpdate.SetUploadName(song.GetFilenameWithoutArrangement())
				APIFileToUpdate.LoadFromFile(song.Filename)
				err = APIFileToUpdate.Save()
				if err != nil {
					log.Println(err)
				}
				continue
			}
			if song.CCLI != "" {
				ctSong := filterSongs(songs, "CCLI", song.CCLI)
				if ctSong.ID != 0 {
					// found a matching song on ChurchTools
					log.Println(ctSong)
					song.SetID(ctSong.ID, ctSong.GetDefaultArrangement())
					song.Title = ctSong.Bezeichnung
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
				id := churchtools.AddSong(song.Title, song.Author, song.Copyright, song.CCLI, song.KeyOfArrangement, "", "", cCLISongCat)
				ctAPISong := churchtools.GetSong(id)
				ctSong = ctAPISong.ToSong()
				songs[fmt.Sprintf("%d", ctSong.ID)] = ctSong
				err = song.UploadToArrangement(ctSong.GetDefaultArrangement(), duplicates)
				if err != nil {
					log.Fatalf("Error Uploading: %v", err)
				}
				continue
			}
			if song.Title != "" {
				ctSong := filterSongs(songs, "Title", song.Title)

				if ctSong.ID != 0 && ctSong.Author == song.Author {
					log.Println(ctSong)
					song.SetID(ctSong.ID, ctSong.GetDefaultArrangement())
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
				err = song.UploadToArrangement(ctSong.GetDefaultArrangement(), duplicates)
				if err != nil {
					log.Fatalf("Error Uploading: %v", err)
				}
				continue
			}
			log.Fatal("Shouldn't reach this")
			// err := song.FixFilename()
			// if err != nil {
			// 	log.Printf("Cannot fix filename %s", err)
			// 	song.MoveToDuplicates(duplicates)
			// }
			// count++
			// if count > 10 {
			// 	break
			// }
		}
	}
	return songs
}
