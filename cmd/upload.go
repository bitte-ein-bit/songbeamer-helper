package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/bitte-ein-bit/songbeamer-helper/log"
	"github.com/bitte-ein-bit/songbeamer-helper/songbeamer"
	"github.com/bitte-ein-bit/songbeamer-helper/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(cmdCTUpload)
}

var cmdCTUpload = &cobra.Command{
	Use:   "ct-upload",
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
	log.Infof("syncing to ChurchTools")
	songs, err := churchtools.GetSongs()
	util.CheckForError(err)

	// for _, song := range songs {
	// 	song.Delete()
	// }
	// songs = nil

	_ = processSongbeamerSongs(songs)
	log.Infof("done")
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
			log.Fatalf("No default filter")
		}

	}
	return churchtools.Song{}
}

func processSongbeamerSongs(songs map[string]churchtools.Song) map[string]churchtools.Song {
	path := viper.GetString("songspath")
	duplicates := viper.GetString("duplicates")

	log.Infof("Lese Songbeamer Songs in %v", path)
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
		log.Infof("Working on %s", song.Title)
		if song.ID != "" {
			// Song has an ID, so it should exists in ChurchTools
			ctSong := filterSongs(songs, "ID", song.ChurchToolsID)
			if ctSong.ID == 0 {
				log.Infof("Die Song ID (%d) konnte auf ChurchTools nicht gefunden werden. Lösche ID und überspringe", song.ChurchToolsID)
				song.ID = ""
				song.Save()
				err = song.FixFilename()
				if err != nil {
					log.Warnf("Cannot fix filename %s", err)
					song.MoveToDuplicates(duplicates)
				}
				continue
			}

			a := song.ExtractArrangementFromFilename()
			if a == "" {
				a = "Standard-Arrangement"
			}
			log.Infof("Lied hat ChurchTools ID %d, Arrangement anhand Dateiname ist %s, Arrangement anhand ID ist %s", song.ChurchToolsID, a, song.ChurchToolsArrangement)

			var arrangement churchtools.SongArrangement
			for _, arrange := range ctSong.Arrangements {
				if arrange.Bezeichnung == a {
					log.Debugf("Arrangement found, checking if newer")
					arrangement = arrange
				}
			}
			if arrangement.ID == 0 {
				log.Infof("Scheinbar wurde die Datei zu einem neuen Arrangement %s umbenannt", a)
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
					log.Infof("Cannot fix filename %s", err)
					song.MoveToDuplicates(duplicates)
				}
				continue
			}
			newer := true

			err = song.FixFilename()
			if err != nil {
				log.Infof("Cannot fix filename %s", err)
				song.MoveToDuplicates(duplicates)
				continue
			}
			APIFileToUpdate := *churchtools.NewSongAPIFile(song.Filename, arrangement.ID)

			for _, file := range arrangement.Files {
				log.Infof("Checking %s aka %s", file.Bezeichnung, file.Filename)
				if file.Bezeichnung != song.GetFilenameWithoutArrangement() {
					log.Infof("no match: %s != %s", file.Bezeichnung, song.GetFilenameWithoutArrangement())
					continue
				}
				ctDate, _ := file.GetModificationDate()
				sngDate, _ := song.GetModificationDate()
				if sngDate.After(ctDate) {
					log.Infof("ChurchTools ist älter: %v < %v", ctDate, sngDate)
					APIFileToUpdate = file.ToAPIFile()
					continue
				}
				log.Infof("ChurchTools ist neuer: %v >= %v", ctDate, sngDate)
				newer = false
			}
			if !newer {
				continue
			}
			log.Debugf("Datei muss auf ChurchTools aktualisiert werden")
			APIFileToUpdate.SetUploadName(song.GetFilenameWithoutArrangement())
			APIFileToUpdate.LoadFromFile(song.Filename)
			err = APIFileToUpdate.Save()
			if err != nil {
				log.Debugf("%v", err)
			}
			continue
		}
		if song.CCLI != "" {
			ctSong := filterSongs(songs, "CCLI", song.CCLI)
			if ctSong.ID != 0 {
				// found a matching song on ChurchTools
				log.Debugf("%v", ctSong)
				song.SetID(ctSong.ID, ctSong.GetDefaultArrangement())
				song.Title = ctSong.Bezeichnung
				song.Save()
				song.SetKeyOfArrangement(ctSong.GetDefaultArrangement())
				// TODO arrangement erkennen
				err := song.FixFilename()
				if err != nil {
					log.Infof("Cannot fix filename %s", err)
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
				log.Debugf("%v", ctSong)
				song.SetID(ctSong.ID, ctSong.GetDefaultArrangement())
				song.Save()
				song.SetKeyOfArrangement(ctSong.GetDefaultArrangement())
				// TODO arrangement erkennen
				err := song.FixFilename()
				if err != nil {
					log.Infof("Cannot fix filename %s", err)
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
		log.Fatalf("Shouldn't reach this")
	}
	return songs
}
