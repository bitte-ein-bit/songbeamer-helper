package churchtools

// A SongArrangement describes how the song could be arranged. It has optionally one or more files atached
type SongArrangement struct {
	ID           int                 `json:"id,string"`
	Bezeichnung  string              `json:"bezeichnung"`
	Default      int                 `json:"default_yn,string"`
	Tonality     string              `json:"tonality"`
	BPM          string              `json:"bpm"`
	Beat         string              `json:"beat"`
	Minutes      int                 `json:"length_min,string"`
	Seconds      int                 `json:"length_sec,string"`
	Note         string              `json:"note,omitempty"`
	ModifiedDate string              `json:"modified_date"`
	ModifiedPID  int                 `json:"modified_pid,string"`
	Files        map[string]SongFile `json:"files,omitempty"`
}
