package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/bitte-ein-bit/songbeamer-helper/songbeamer"
	"github.com/google/martian/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(cmdCTDownload)
}

var cmdCTDownload = &cobra.Command{
	Use:   "ct-download",
	Short: "Download songs for event from Churchtools",
	Long:  `Downloads songs from CT that are listed on the selected event.`,
	Run: func(cmd *cobra.Command, args []string) {
		downloadFromCTEvent()
	},
}

func downloadFromCTEvent() {
	events := churchtools.GetEvents(8)
	if len(events) == 0 {
		fmt.Printf("In den nächsten 6 Tagen wurden keine Verantstaltungen gefunden.")
		return
	}
	event := ask(events)
	songs := event.GetSongs()
	if len(songs) == 0 {
		fmt.Print("Es sind keine Songs in der Agenda hinterlegt.")
	}

	c := churchtools.CTClient{}
	c.Login()

	path := viper.GetString("songspath")
	for _, song := range songs {
		_, err := DownloadSongbeamerFile(c, song, path)
		if err != nil {
			log.Errorf("Cannot download song: %v", err)
		}

	}
}

func ask(events []churchtools.Event) (event churchtools.Event) {
	for {
		for key, value := range events {
			fmt.Printf("%d. %s: %s - %s\n", key+1, value.StartDate.Format("02.01.2006 15:04"), value.Name, value.Description)
		}

		var input string
		fmt.Print("Bitte wähle eine Veranstaltung aus, für die du die Songs herunterladen möchtest: ")

		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Printf("Fehlerhafte Eingabe: %v\n", err)
			continue
		}

		// Verify we got an integer.
		selected, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("Ungültige Eingabe '%s'\n", input)
			continue
		}

		// Verify selection is within range.
		if selected < 1 || selected > len(events) {
			fmt.Printf("Ungültiger Wert %d. Gültige Werte sind: 1-%d\n", selected, len(events))
			continue
		}

		// Translate user-selected index back to zero-based index.
		event = events[selected-1]
		break
	}

	return
}

// DownloadSongbeamerFile downloads a file from Churchtools so it can be used with Songbeamer
func DownloadSongbeamerFile(c churchtools.CTClient, s churchtools.APISong, path string) (files []string, err error) {
	if path == "" {
		return nil, fmt.Errorf("Cannot save to an empty path: %s", path)
	}
	for _, a := range s.Arrangements {
		for _, f := range a.Files {
			if filepath.Ext(f.Name) != ".sng" {
				continue
			}
			resp := c.GetRequest(f.FileURL, nil)
			defer resp.Body.Close()

			last := resp.Header.Get("Last-Modified")
			lastTime, err := time.Parse(time.RFC1123, last)
			if err != nil {
				fmt.Println(err)
				lastTime = time.Now()
			}

			filename := fmt.Sprintf("%s/%s - %s.sng", path, s.Bezeichnung, a.Name)
			out, err := os.Create(filename)
			if err != nil {
				return nil, fmt.Errorf("Cannot create file: %w", err)
			}
			defer out.Close()
			io.Copy(out, resp.Body)
			err = os.Chtimes(filename, lastTime, lastTime)
			if err != nil {
				fmt.Printf("Cannot adjust time on file. Ignoring error: %s", err)
			}
			sng := songbeamer.Song{}
			sng.LoadFromFile(filename)
			sng.Validate(s, a)
			sng.UploadIfNeeded(&f, lastTime)
			files = append(files, filename)
		}
	}
	return
}
