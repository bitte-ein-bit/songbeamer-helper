package songbeamer

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/bitte-ein-bit/songbeamer-helper/util"
)

type SongbeamerSong struct {
	ID                       string
	ChurchToolsID            int
	ChurchToolsArrangementID int
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
	Verses                   [][]string
	Headers                  map[string]string
}

// LoadFromFile loads a sng file into a SongbeamerSong struct
func (s *SongbeamerSong) LoadFromFile(filename string) {
	inHeader := true
	var verse []string
	s.Filename = filename
	lines, err := util.File2lines(filename)
	util.CheckForError(err)

	for _, line := range lines {
		if line == "---" {
			inHeader = false
			if len(verse) > 0 {
				s.Verses = append(s.Verses, verse)
				verse = []string{}
			}
			continue
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
				tmp := strings.SplitN(header[1], "-", 3)
				ID, err := strconv.Atoi(tmp[0])
				if err != nil {
					fmt.Println("Invalid ID field, ignoring")
				}
				s.ChurchToolsID = ID
				if len(tmp) > 1 {
					ID, err = strconv.Atoi(tmp[1])
					if err != nil {
						fmt.Println("Invalid ID field, ignoring")
					}
					s.ChurchToolsArrangementID = ID
					s.ChurchToolsArrangement = tmp[2]
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
				name := strings.Trim(header[0], "#")
				s.Headers[name] = header[1]
				if name == "Melody" && s.Author == "" {
					s.Author = header[1]
				}
			}
			continue
		}
		verse = append(verse, line)
		// if line == "--" {
		// 	continue
		// }
	}
	if len(verse) > 0 {
		s.Verses = append(s.Verses, verse)
	}
}

// FixFilename moves files arround according to their name
func (s *SongbeamerSong) FixFilename() error {
	id := ""
	if s.ChurchToolsArrangement != "" {
		id = fmt.Sprintf(" - %s", s.ChurchToolsArrangement)
	}
	filenameByTitle := fmt.Sprintf("%s%s.sng", strings.Replace(s.Title, "/", "_", -1), id)
	if strings.ToLower(filepath.Base(s.Filename)) != strings.ToLower(filenameByTitle) {
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

// MoveToDuplicates moves the Songbeamer file out of the way
func (s *SongbeamerSong) MoveToDuplicates(path string) error {
	f, err := os.Open(s.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("Can't compute MD5: %w", err)
	}
	newFilename := fmt.Sprintf("%s/%s-%s.sng", path, strings.Replace(s.Title, "/", "_", -1), fmt.Sprintf("%x", h.Sum(nil)))
	err = os.Rename(s.Filename, newFilename)
	util.CheckForError(err)
	s.Filename = newFilename
	return nil
}

// SetID adds a ChurchTools Song ID to a Songbeamer File
func (s *SongbeamerSong) SetID(songID int, arrangement churchtools.SongArrangement) {
	s.ChurchToolsID = songID
	s.ChurchToolsArrangementID = arrangement.ID
	s.ChurchToolsArrangement = arrangement.Bezeichnung
	s.ID = fmt.Sprintf("%d-%d-%s", s.ChurchToolsID, s.ChurchToolsArrangementID, s.ChurchToolsArrangement)
}

// SetKeyOfArrangement adds a Key to a Songbeamer File
func (s *SongbeamerSong) SetKeyOfArrangement(arrangement churchtools.SongArrangement) {
	nT := strings.TrimSpace(arrangement.Tonality)
	if nT == "" {
		return
	}
	if strings.TrimSpace(s.KeyOfArrangement) == nT {
		return
	}
	if s.KeyOfArrangement != "" {
		log.Fatalf("missmatch of keys: %s vs. %s", s.KeyOfArrangement, nT)
		return
	}
	s.KeyOfArrangement = nT
	if s.Filename == "" {
		log.Fatal("Cannot save to non-set file")
	}
	log.Printf("Adding ID to %v", s.Filename)
	line := fmt.Sprintf("#Key=%s\n", s.KeyOfArrangement)
	err := util.InsertStringToFile(s.Filename, line, 1)
	util.CheckForError(err)
}

func nonEmptyHeader(name, value string) string {
	if value != "" {
		return fmt.Sprintf("#%s=%s\r\n", name, value)
	}
	return ""
}

func (s *SongbeamerSong) Save() error {
	fileContent := string('\uFEFF') // Add BOM for Songbeamer

	fileContent += nonEmptyHeader("Author", s.Author)
	fileContent += nonEmptyHeader("(c)", s.Copyright)
	fileContent += nonEmptyHeader("Key", s.KeyOfArrangement)
	fileContent += nonEmptyHeader("ID", s.ID)
	fileContent += nonEmptyHeader("CCLI", s.CCLI)
	fileContent += nonEmptyHeader("Title", s.Title)
	fileContent += nonEmptyHeader("VerseOrder", s.VerseOrder)
	fileContent += nonEmptyHeader("Categories", strings.Join(s.Category, ","))

	for name, value := range s.Headers {
		fileContent += nonEmptyHeader(name, value)
	}

	for _, lines := range s.Verses {
		fileContent += "---\r\n"
		for _, line := range lines {
			fileContent += fmt.Sprintf("%s\r\n", line)
		}
	}

	f, err := os.Create(s.Filename)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprint(f, fileContent)

	if err != nil {
		return err
	}

	err = f.Close()
	return err
}

func (s *SongbeamerSong) ExtractArrangementFromFilename() string {
	re := regexp.MustCompile("(.+) - (.+)\\.sng")
	data := re.FindStringSubmatch(s.Filename)
	if len(data) == 0 {
		return ""
	}
	return data[2]
}
