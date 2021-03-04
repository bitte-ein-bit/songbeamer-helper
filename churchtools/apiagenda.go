package churchtools

import (
	"fmt"
	"strings"
)

type apiAgendaResponse struct {
	Agenda APIAgenda `json:"data"`
}

// APIAgenda describes the agenda as returned by the REST API
type APIAgenda struct {
	ID      int             `json:"id"`
	Name    string          `json:"name"`
	IsFinal bool            `json:"isFinal"`
	Items   []APIAgendaItem `json:"items"`
}

// APIAgendaSong describes a song as shown in an agenda item
type APIAgendaSong struct {
	SongID        int    `json:"songId"`
	ArrangementID int    `json:"arrangementId"`
	Title         string `json:"title"`
	Arrangement   string `json:"arrangement"`
	Category      string `json:"category"`
	Key           string `json:"key"`
	BPM           string `json:"bpm"`
	IsDefault     bool   `json:"isDefault"`
}

func (s *APIAgendaSong) ToFilename() (filename string) {
	filename = fmt.Sprintf("%s - %s.sng", strings.Replace(s.Title, "/", "_", -1), s.Arrangement)
	return
}
