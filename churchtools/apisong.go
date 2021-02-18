package churchtools

import (
	"fmt"
	"log"
)

// The APISong is describing the api response of a Song. It is further defined by one or more Arrangements
type APISong struct {
	ID           int                  `json:"id"`
	Bezeichnung  string               `json:"name"`
	Category     APISongCategory      `json:"category"`
	Practice     bool                 `json:"shouldPractice"`
	Author       string               `json:"author"`
	CCLI         string               `json:"ccli"`
	Copyright    string               `json:"copyright"`
	Note         string               `json:"note"`
	Arrangements []APISongArrangement `json:"arrangements"`
}

// GetDefaultArrangement returns the APISongArrangement that is marked as default. It's unsure if the API returns nothing
func (s *APISong) GetDefaultArrangement() (ret APISongArrangement) {
	for _, value := range s.Arrangements {
		if value.Default {
			log.Printf("default arrangement: %d - %s", value.ID, value.Name)
			return value
		}
	}
	log.Print("No default arrangement found")
	return
}

// ToSong converts an APISong to a Song (old to new API)
func (s *APISong) ToSong() (ret Song) {
	a := make(map[string]SongArrangement)
	for _, value := range s.Arrangements {
		a[fmt.Sprintf("%d", value.ID)] = value.ToArrangement()
	}
	ret = Song{
		ID:             s.ID,
		Bezeichnung:    s.Bezeichnung,
		SongcategoryID: s.Category.ID,
		// Practice: s.Practice,
		Author:       s.Author,
		CCLI:         s.CCLI,
		Copyright:    s.Copyright,
		Note:         s.Note,
		Arrangements: a,
	}
	return
}
