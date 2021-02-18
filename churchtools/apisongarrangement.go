package churchtools

import (
	"math"
)

// APISongArrangement represents a song arrangement using the REST API
type APISongArrangement struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Default          bool      `json:"isDefault"`
	KeyOfArrangement string    `json:"keyOfArrangement"`
	BPM              string    `json:"bpm"`
	Beat             string    `json:"beat"`
	Duration         int       `json:"duration"`
	Note             string    `json:"note,omitempty"`
	Files            []APIFile `json:"files,omitempty"`
	Links            []APIFile `json:"links,omitempty"`
}

// ToArrangement converts to an old style API arrangement
func (a *APISongArrangement) ToArrangement() (ret SongArrangement) {
	d := 0
	if a.Default {
		d = 1
	}
	ret = SongArrangement{
		ID:          a.ID,
		Bezeichnung: a.Name,
		Default:     d,
		Tonality:    a.KeyOfArrangement,
		BPM:         a.BPM,
		Beat:        a.Beat,
		Seconds:     int(a.Duration % 60),
		Minutes:     int(math.Floor(float64(a.Duration) / 60)),
		Note:        a.Note,
		// Files: map[string]SongFile,
	}
	return
}
