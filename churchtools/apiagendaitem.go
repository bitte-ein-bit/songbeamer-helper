package churchtools

import "fmt"

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
	caption = i.Title
	switch i.Type {
	case "song":
		color = "clBlue"
		caption = fmt.Sprintf("%s - %s", i.Song.Title, i.Song.Arrangement)
		extra += fmt.Sprintf("      FileName = '%s'\n", i.Song.ToFilename())
	case "header":
		color = "16711680"
	case "normal":
		color = "33023"
	case "intro":
		color = "clBlack"

		extra += fmt.Sprint("      FileName = 'C:\\Program Files (x86)\\SongBeamer\\Intro.wav'\n")
	}
	text += fmt.Sprintf("      Caption = '%s'\n      Color = %s\n%s    end", caption, color, extra)
	return
}
