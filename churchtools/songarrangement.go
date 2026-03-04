package churchtools

import "time"

// A SongArrangement describes how the song could be arranged. It has optionally one or more files ataced
type SongArrangement struct {
	ID           int                 `json:"id,string"`
	Bezeichnung  string              `json:"bezeichnung"`
	Name         string              `json:"name"` // Support for new API format
	Default      int                 `json:"default_yn,string"`
	Tonality     string              `json:"tonality"`
	BPM          string              `json:"bpm"`
	Beat         string              `json:"beat"`
	Minutes      int                 `json:"length_min"`
	Seconds      int                 `json:"length_sec"`
	Note         string              `json:"note,omitempty"`
	ModifiedDate string              `json:"modified_date"`
	ModifiedPID  int                 `json:"modified_pid,string"`
	Files        map[string]SongFile `json:"files,omitempty"`
}

// GetModificationDate parses the date string returned by the API into a time struct
func (s *SongArrangement) GetModificationDate() (t time.Time, err error) {
	layout := "2006-01-02 15:04:05"
	t, err = time.Parse(layout, s.ModifiedDate)
	return
}

// GetName returns the arrangement name, supporting both old and new API formats
func (s *SongArrangement) GetName() string {
	if s.Bezeichnung != "" {
		return s.Bezeichnung
	}
	return s.Name
}
