package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/bitte-ein-bit/songbeamer-helper/log"
	"github.com/bitte-ein-bit/songbeamer-helper/songbeamer"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
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
		event := selectCTEvent()
		downloadSongsForCTEvent(event)
		createSongbeamerAgenda(event)
	},
}

func selectCTEvent() (event churchtools.Event) {
	events := churchtools.GetEvents(8)
	if len(events) == 0 {
		log.Errorf("In den nächsten 6 Tagen wurden keine Verantstaltungen gefunden.")
		return
	}
	event = ask(events)
	return
}

func downloadSongsForCTEvent(event churchtools.Event) {
	songs := event.GetSongs()
	if len(songs) == 0 {
		log.Errorf("Es sind keine Songs in der Agenda hinterlegt.")
		return
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
			fmt.Printf("%d. %s: %s - %s\n", key+1, value.StartDate.Local().Format("02.01.2006 15:04"), value.Name, value.Description)
		}

		var input string
		fmt.Print("Bitte wähle eine Veranstaltung aus, für die du die Songs herunterladen möchtest: ")

		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Errorf("Fehlerhafte Eingabe: %v\n", err)
			continue
		}

		// Verify we got an integer.
		selected, err := strconv.Atoi(input)
		if err != nil {
			log.Errorf("Ungültige Eingabe '%s'\n", input)
			continue
		}

		// Verify selection is within range.
		if selected < 1 || selected > len(events) {
			log.Errorf("Ungültiger Wert %d. Gültige Werte sind: 1-%d\n", selected, len(events))
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
				log.Warnf("Kann letztes Änderungsdatum nicht auswerten, verwende jetzt")
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
				log.Warnf("Cannot adjust time on file. Ignoring error: %s", err)
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

func createSongbeamerAgenda(event churchtools.Event) {
	a := event.GetAgenda()
	content := "object AblaufPlanItems: TAblaufPlanItems\n  items = <"
	for _, item := range a.Items {
		content += "\n" + item.ToSongbeamerItem()
	}
	content += `
	item
      Color = 33023
    end
    item
      Color = 33023
    end
    item
      Color = 33023
    end
    item
      Color = 33023
    end
    item
      Caption = 'Immer mal n'#252'tzlich'
      Color = 33023
    end
    item
      Caption = 'Gebet - Vater unser'
      Color = clBlue
      FileName = 'Gebet - Vater unser.sng'
    end
    item
      Caption = 'Das Apostolische Glaubensbekenntnis'
      Color = clBlue
      FileName = 'Das Apostolische Glaubensbekenntnis.sng'
    end`
	content += ">\nend"

	encoded := ""
	for _, s := range content {
		format := "%c"
		if s > 200 {
			format = "'#%d'"
		}
		encoded += fmt.Sprintf(format, s)
	}

	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Cannot parse home dir: %s", err)
	}
	filename := fmt.Sprintf("%s/Desktop/Ablaufplan_%s.col", home, event.StartDate.Local().Format("2006-01-02_15-04"))
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		log.Fatalf("Ein Fehler ist beim Erstellen des Ablaufplans aufgetreten: %s", err)
	}

	_, err = fmt.Fprint(f, encoded)

	if err != nil {
		log.Fatalf("Ein Fehler ist beim Schreiben des Ablaufplans aufgetreten: %s", err)
	}

	color.Set(color.FgGreen)
	fmt.Printf("Der Ablaufplan wurde nach %s gespeichert.\n", filename)

	// color.Set(color.FgRed)
	// log.Infof("Das Erstellen des Ablaufplans ist noch nicht implementiert.")
	// log.Infof("Bitte lade den Ablaufplan aus Churchtools herunter")
}
