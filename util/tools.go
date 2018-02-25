package util

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/igungor/chardet"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func CheckForError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

/**
 * Insert sting to n-th line of file.
 * If you want to insert a line, append newline '\n' to the end of the string.
 */
func InsertStringToFile(path, str string, index int) error {
	lines, err := File2lines(path)
	if err != nil {
		return err
	}

	fileContent := string('\uFEFF') // Add BOM for Songbeamer
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\r\n"
	}

	f, err := os.Create(path)
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

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	enc := DetectEncoding(filePath)
	return LinesFromReader(f, enc)
}

func LinesFromReader(r io.Reader, c *charmap.Charmap) ([]string, error) {
	var lines []string
	if c != nil {
		r = transform.NewReader(r, c.NewDecoder())
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func DetectEncoding(path string) *charmap.Charmap {
	dat, err := ioutil.ReadFile(path)
	CheckForError(err)
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(dat)
	CheckForError(err)

	log.Printf("Detected charset for %s is %s", path, result.Charset)
	switch result.Charset {
	case "ISO-8859-1":
		return charmap.Windows1252
	case "ISO-8859-2":
		return charmap.Windows1252
	case "ISO-8859-9":
		return charmap.Windows1252
	case "ISO-8859-15":
		return charmap.Windows1252
	case "windows-1252":
		return charmap.Windows1252
	case "UTF-8":
		return nil
	default:
		log.Fatalf("Encoding %s not supported", result.Charset)
		return nil
	}

}
