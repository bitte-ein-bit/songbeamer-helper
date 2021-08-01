package songbeamer

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/bitte-ein-bit/songbeamer-helper/log"
	"github.com/bitte-ein-bit/songbeamer-helper/util"
)

// Song represents the data stored inside an sng file
type Song struct {
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
	VerseIdentifier          []string
	Filename                 string
	Category                 []string
	Verses                   [][]string
	Headers                  map[string]string
}

func isIdentifier(s string) (status bool, identifier string) {
	identifier = strings.Split(s, " ")[0]
	status = false
	switch strings.ToLower(identifier) {
	case
		"unbekannt",
		"unbenannt",
		"unknown",
		"intro",
		"vers",
		"verse",
		"strophe",
		"pre-bridge",
		"bridge",
		"misc",
		"pre-refrain",
		"refrain",
		"pre-chorus",
		"chorus",
		"pre-coda",
		"zwischenspiel",
		"instrumental",
		"interlude",
		"coda",
		"ending",
		"outro",
		"teil",
		"part",
		"chor",
		"solo",
		"stop":
		status = true
		return
	}
	if strings.HasPrefix(identifier, "$$M=") {
		status = true
		// identifier = strings.SplitN(identifier, "=", 2)[1]
	}
	identifier = ""

	return

}

// LoadFromFile loads a sng file into a Song struct
func (s *Song) LoadFromFile(filename string) {
	inHeader := true
	firstLineOfVerse := false
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
			firstLineOfVerse = true
			continue
		}
		if inHeader {
			header := strings.SplitN(line,
				"=", 2)
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
					log.Infof("Invalid ID field, ignoring")
				}
				s.ChurchToolsID = ID
				if len(tmp) > 1 {
					ID, err = strconv.Atoi(tmp[1])
					if err != nil {
						log.Infof("Invalid ID field, ignoring")
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
				if len(header) < 2 {
					log.Debugf("Header hat ungültige Zeile: %s", line)
					continue
				}
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
		if firstLineOfVerse {
			firstLineOfVerse = false
			status, _ := isIdentifier(line)
			if status {
				s.VerseIdentifier = append(s.VerseIdentifier, line)
			}
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
func (s *Song) FixFilename() error {
	id := ""
	if s.ChurchToolsArrangement != "" {
		id = fmt.Sprintf(" - %s", s.ChurchToolsArrangement)
	}
	filenameByTitle := fmt.Sprintf("%s%s.sng", s.saveFilename(s.Title), id)
	if strings.ToLower(filepath.Base(s.Filename)) != strings.ToLower(filenameByTitle) {
		log.Debugf("%s should be named %s", s.Filename, filenameByTitle)
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

// GetFilenameWithoutArrangement constructs the filename without the arrangement addendum
func (s *Song) GetFilenameWithoutArrangement() string {
	return fmt.Sprintf("%s.sng", s.saveFilename(s.Title))
}

func (s *Song) saveFilename(name string) (n string) {
	n = strings.Replace(s.Title, "/", "_", -1)
	n = strings.Replace(n, "?", "_", -1)
	n = strings.Replace(n, ":", "_", -1)
	return
}

// MoveToDuplicates moves the Songbeamer file out of the way
func (s *Song) MoveToDuplicates(path string) error {
	f, err := os.Open(s.Filename)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("Can't compute MD5: %w", err)
	}
	f.Close()
	newFilename := fmt.Sprintf("%s/%s-%s.sng", path, s.saveFilename(s.Title), fmt.Sprintf("%x", h.Sum(nil)))
	err = os.Rename(s.Filename, newFilename)
	util.CheckForError(err)
	s.Filename = newFilename
	return nil
}

// SetID adds a ChurchTools Song ID to a Songbeamer File
func (s *Song) SetID(songID int, arrangement churchtools.SongArrangement) {
	s.ChurchToolsID = songID
	s.ChurchToolsArrangementID = arrangement.ID
	s.ChurchToolsArrangement = arrangement.Bezeichnung
	s.ID = fmt.Sprintf("%d-%d-%s", s.ChurchToolsID, s.ChurchToolsArrangementID, s.ChurchToolsArrangement)
}

// SetKeyOfArrangement adds a Key to a Songbeamer File
func (s *Song) SetKeyOfArrangement(arrangement churchtools.SongArrangement) {
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
		log.Fatalf("Cannot save to non-set file")
	}
	log.Infof("Adding ID to %v", s.Filename)
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

// Save updates the sng file on disk
func (s *Song) Save() error {
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
		log.Fatalf("Cannot create file: %s", err)
	}
	defer f.Close()

	_, err = fmt.Fprint(f, fileContent)
	return err
}

// ExtractArrangementFromFilename derives the arrangement name from the sng filename
func (s *Song) ExtractArrangementFromFilename() string {
	re := regexp.MustCompile("(.+) - (.+)\\.sng")
	data := re.FindStringSubmatch(s.Filename)
	if len(data) == 0 {
		return ""
	}
	return data[2]
}

// UploadToArrangement attaches a song to a ChurchTools Songbeamer arrangement
func (s *Song) UploadToArrangement(arrangement churchtools.SongArrangement, duplicatesPath string) error {
	s.SetID(s.ChurchToolsID, arrangement)
	s.Save()
	err := s.FixFilename()
	if err != nil {
		log.Warnf("Cannot fix filename %s", err)
		s.MoveToDuplicates(duplicatesPath)
		return nil
	}
	ctAPIFile := churchtools.NewSongAPIFile(s.Filename, arrangement.ID)
	ctAPIFile.SetUploadName(s.GetFilenameWithoutArrangement())
	err = ctAPIFile.Save()
	return err
}

// GetModificationDate returns the modification date of the file
func (s *Song) GetModificationDate() (t time.Time, err error) {
	fi, err := os.Stat(s.Filename)
	if err != nil {
		return
	}
	t = fi.ModTime().Round(time.Second)
	return
}

// Validate makes sure relevant data from CT is embedded
func (s *Song) Validate(apiSong churchtools.APISong, a churchtools.APISongArrangement) (err error) {
	changed := false
	if s.ChurchToolsID != apiSong.ID || s.ChurchToolsArrangementID != a.ID || s.ChurchToolsArrangement != a.Name {
		log.Debugf("Setze ID Feld anhand von ChurchTools")
		s.SetID(apiSong.ID, a.ToArrangement())
		changed = true
	}
	if s.CCLI != apiSong.CCLI {
		log.Infof("Setze CCLI Feld anhand von ChurchTools")
		s.CCLI = apiSong.CCLI
		changed = true
	}
	if s.Author != apiSong.Author {
		log.Infof("Setze Autor Feld anhand von ChurchTools")
		s.Author = apiSong.Author
		changed = true
	}

	if s.Title != apiSong.Bezeichnung {
		log.Infof("Setze Titel Feld anhand von ChurchTools")
		s.Title = apiSong.Bezeichnung
		changed = true
	}

	if s.Copyright != apiSong.Copyright {
		log.Infof("Setze Copyright Feld anhand von ChurchTools")
		s.Copyright = apiSong.Copyright
		changed = true
	}

	if changed {
		err = s.Save()
		if err != nil {
			return fmt.Errorf("Cannot save SNG file: %w", err)
		}
		log.Infof("Datei aktualisiert")
	}
	return
}

func (s *Song) UploadIfNeeded(a *churchtools.APIFile, lastChanged time.Time) {
	uploadNeeded := false
	if a.Name != s.GetFilenameWithoutArrangement() {
		log.Infof("Dateiname auf ChurchTools (%s) stimmt nicht, korrigiere zu %s", a.Name, s.GetFilenameWithoutArrangement())
		a.SetUploadName(s.GetFilenameWithoutArrangement())
		uploadNeeded = true
	}

	ctDate := lastChanged.Round(time.Second)
	sngDate, _ := s.GetModificationDate()
	if sngDate.After(ctDate) {
		log.Debugf("ChurchTools ist älter: %v < %v, lade Datei zu ChurchTools hoch", ctDate, sngDate)
		uploadNeeded = true
	}

	if uploadNeeded {
		a.LoadFromFile(s.Filename)
		a.Save()
	}

}
