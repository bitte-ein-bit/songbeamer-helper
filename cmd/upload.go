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
	RootCmd.AddCommand(cmdCTUpload)
}

var cmdCTUpload = &cobra.Command{
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
	log.Println("syncing to ChurchTools")
	songs, err := churchtools.GetSongs()
	util.CheckForError(err)

	// for _, song := range songs {
	// 	song.Delete()
	// }
	// songs = nil

	_ = processSongbeamerSongs(songs)
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
		if filepath.Ext(file.Name()) != ".sng" {
			continue
		}
		// time.Sleep(time.Millisecond * 100)
		fullpath := filepath.Join(path, file.Name())
		song := songbeamer.Song{}
		song.LoadFromFile(fullpath)
		log.Printf("Working on %s", song.Title)
		if song.ID != "" {
			// Song has an ID, so it should exists in ChurchTools
			ctSong := filterSongs(songs, "ID", song.ChurchToolsID)
			if ctSong.ID == 0 {
				log.Printf("CT Song ID (%d) can't be found, resetting ID and skip", song.ChurchToolsID)
				song.ID = ""
				song.Save()
				err = song.FixFilename()
				if err != nil {
					log.Printf("Cannot fix filename %s", err)
					song.MoveToDuplicates(duplicates)
				}
				continue
			}

			a := song.ExtractArrangementFromFilename()
			if a == "" {
				a = "Standard-Arrangement"
			}
			log.Printf("Song has CT ID of %d, arrangement by Filename is %s, arrangement by ID is %s", song.ChurchToolsID, a, song.ChurchToolsArrangement)

			var arrangement churchtools.SongArrangement
			for _, arrange := range ctSong.Arrangements {
				if arrange.Bezeichnung == a {
					log.Println("Arrangement found, checking if newer")
					arrangement = arrange
				}
			}
			if arrangement.ID == 0 {
				log.Printf("Seems song was renamed to new arrangement: %s", a)
				arrangementID, err := ctSong.AddArrangement(a)
				util.CheckForError(err)
				arrangement = churchtools.SongArrangement{
					ID:          arrangementID,
					Bezeichnung: a,
				}
				err = song.UploadToArrangement(arrangement, duplicates)
				if err != nil {
					log.Fatalf("Error Uploading: %v", err)
				}
				song.SetID(ctSong.ID, arrangement)
				err = song.FixFilename()
				if err != nil {
					log.Printf("Cannot fix filename %s", err)
					song.MoveToDuplicates(duplicates)
				}
				continue
			}
			newer := true

			err = song.FixFilename()
			if err != nil {
				log.Printf("Cannot fix filename %s", err)
				song.MoveToDuplicates(duplicates)
				continue
			}
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

			song.SetID(ctSong.ID, ctSong.GetDefaultArrangement())
			song.Save()
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
			song.SetID(ctSong.ID, ctSong.GetDefaultArrangement())
			song.Save()
			err = song.UploadToArrangement(ctSong.GetDefaultArrangement(), duplicates)
			if err != nil {
				log.Fatalf("Error Uploading: %v", err)
			}
			continue
		}
		song.MoveToDuplicates(duplicates)
		log.Fatal("Shouldn't reach this")
	}
	return songs
}
