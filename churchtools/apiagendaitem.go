package churchtools

import (
	"fmt"
	"strings"
)

// APIAgendaItem describes an item in APIAgenda
type APIAgendaItem struct {
	ID            int           `json:"id"`
	Position      int           `json:"position"`
	Type          string        `json:"type"`
	Title         string        `json:"title"`
	IsBeforeEvent bool          `json:"isBeforeEvent"`
	Song          APIAgendaSong `json:"song,omitempty"`
	ArrangementID int           `json:"arrangementId,omitempty"`
}

func (i *APIAgendaItem) ToSongbeamerItem() (text string) {
	text = "    item\n"
	var color string
	var caption string
	var extra string
	if i.Title == "Count-Down / Jingle" {
		i.Type = "intro"
	}
	// if title contains Kinderprogramm, set MP3
	if strings.Contains(i.Title, "Kinderprogramm") {
		i.Type = "kinder"
	}
	caption = i.Title
	switch i.Type {
	case "song":
		color = "clBlue"
		caption = fmt.Sprintf("%s - %s", i.Song.Title, i.Song.Arrangement)
		extra += fmt.Sprintf("      FileName = '%s'\n", strings.Replace(i.Song.ToFilename(), "'", "'#39'", -1))
	case "header":
		color = "clBlack"
	case "normal":
		color = "33023"
	case "intro":
		color = "clBlack"
		extra += fmt.Sprint("      FileName = 'C:\\Program Files (x86)\\SongBeamer\\Intro.mp3'\n")
	case "kinder":
		color = "clBlack"
		extra += fmt.Sprint("      FileName = 'C:\\Program Files (x86)\\SongBeamer\\Kinder.mp3'\n")
    default:
		color = "clBlue"
	}
	text += fmt.Sprintf("      Caption = '%s'\n      Color = %s\n%s    end", strings.Replace(caption, "'", "'#39'", -1), color, extra)
	return
}
