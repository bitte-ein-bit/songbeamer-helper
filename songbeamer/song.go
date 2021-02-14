package songbeamer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/bitte-ein-bit/songbeamer-helper/util"
)

type SongbeamerSong struct {
	ID                       string
	ChurchToolsID            string
	ChurchToolsArrangementID string
	ChurchToolsArrangement   string
	CCLI                     string
	Title                    string
	Author                   string
	Copyright                string
	KeyOfArrangement         string
	BPM                      string
	Beat                     string
	VerseOrder               string
	Filename                 string
	Category                 []string
	Verses                   map[string]string
	Headers                  map[string]string
}

// LoadFromFile loads a sng file into a SongbeamerSong struct
func (s *SongbeamerSong) LoadFromFile(filename string) {
	inHeader := true
	s.Filename = filename
	lines, err := util.File2lines(filename)
	util.CheckForError(err)

	for _, line := range lines {
		if line == "---" {
			inHeader = false
			break
		}
		if inHeader {
			header := strings.Split(line, "=")
			switch strings.Trim(header[0], "#") {
			case "(c)":
				s.Copyright = header[1]
			case "Key":
				s.KeyOfArrangement = header[1]
			case "ID":
				s.ID = header[1]
				tmp := strings.Split(header[1], "-")
				s.ChurchToolsID = tmp[0]
				if len(tmp) > 1 {
					s.ChurchToolsArrangementID = tmp[1]
				}
			case "CCLI":
				s.CCLI = header[1]
			case "Title":
				s.Title = header[1]
			case "Author":
				s.Author = header[1]
			case "VerseOrder":
				s.VerseOrder = header[1]
			case "Categories":
				s.Category = strings.Split(header[1], ",")
			default:
				if s.Headers == nil {
					s.Headers = make(map[string]string)
				}
				s.Headers[strings.Trim(header[0], "#")] = header[1]
			}
			continue
		}
		if line == "--" {
			continue
		}
	}
}

// FixFilename moves files arround according to their name
func (s *SongbeamerSong) FixFilename() error {
	id := ""
	if s.ChurchToolsArrangement != "" {
		id = fmt.Sprintf(" - %s", s.ChurchToolsArrangement)
	}
	filenameByTitle := fmt.Sprintf("%s%s.sng", strings.Replace(s.Title, "/", "_", -1), id)
	if filepath.Base(s.Filename) != filenameByTitle {
		log.Printf("%s should be named %s", s.Filename, filenameByTitle)
		newFilename := fmt.Sprintf("%s/%s", filepath.Dir(s.Filename), filenameByTitle)
		if _, err := os.Stat(newFilename); err == nil {
			return fmt.Errorf("New File already exists: %s", newFilename)
		}
		err := os.Rename(s.Filename, newFilename)
		util.CheckForError(err)
		s.Filename = newFilename
	}
	return nil
}

// AddID adds a ChurchTools Song ID to a Songbeamer File
func (s *SongbeamerSong) AddID(songID int, arrangement churchtools.SongArrangement) {
	s.ChurchToolsID = fmt.Sprintf("%d", songID)
	s.ChurchToolsArrangementID = fmt.Sprintf("%d", arrangement.ID)
	s.ChurchToolsArrangement = arrangement.Bezeichnung
	s.ID = fmt.Sprintf("%s-%s-%s", s.ChurchToolsID, s.ChurchToolsArrangementID, s.ChurchToolsArrangement)

	if s.Filename == "" {
		log.Fatal("Cannot save to non-set file")
	}
	log.Printf("Adding ID to %v", s.Filename)
	line := fmt.Sprintf("#ID=%s\n", s.ID)
	err := util.InsertStringToFile(s.Filename, line, 1)
	util.CheckForError(err)
}


// SetKeyOfArrangement adds a Key to a Songbeamer File
func (s *SongbeamerSong) SetKeyOfArrangement(arrangement churchtools.SongArrangement) {
	if s.KeyOfArrangement == arrangement.Tonality {
		return
	}
	if s.KeyOfArrangement != "" {
		log.Fatalf("missmatch of keys: %s vs. %s", s.KeyOfArrangement, arrangement.Tonality)
	}
	s.KeyOfArrangement = arrangement.Tonality
	if s.Filename == "" {
		log.Fatal("Cannot save to non-set file")
	}
	log.Printf("Adding ID to %v", s.Filename)
	line := fmt.Sprintf("#Key=%s\n", s.KeyOfArrangement)
	err := util.InsertStringToFile(s.Filename, line, 1)
	util.CheckForError(err)
}