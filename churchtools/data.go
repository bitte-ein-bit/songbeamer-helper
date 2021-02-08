package churchtools

// A SongFile describes files attached to an SongArrangement
type SongFile struct {
	ID                     int    `json:"id,string"`
	DomainType             string `json:"domain_type"`
	DomainID               int    `json:"domain_id,string"`
	Bezeichnung            string `json:"bezeichnung"`
	Filename               string `json:"filename"`
	ShowonlywheneditableYN int    `json:"showonlywheneditable_yn,string"`
	SecuritylevelID        int    `json:"securitylevel_id,omitempty"`
	ImageOptions           string `json:"image_options"`
	ModifiedDate           string `json:"modified_date"`
	ModifiedPID            int    `json:"modified_pid,string"`
	DeletionDate           string `json:"deletion_date,omitempty"`
	ModifiedUsername       string `json:"modified_username"`
}

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
	Note         string              `json:"note"`
	ModifiedDate string              `json:"modified_date"`
	ModifiedPID  int                 `json:"modified_pid,string"`
	Files        map[string]SongFile `json:"files"`
}

// The Song is describing the api response of a Song. It is further defined by one or more Arrangements
type Song struct {
	ID             int                        `json:"id,string"`
	Bezeichnung    string                     `json:"bezeichnung"`
	SongcategoryID int                        `json:"songcategory_id,string"`
	Practice       int                        `json:"practice_yn,string"`
	Author         string                     `json:"author"`
	CCLI           int                        `json:"ccli,string"`
	Copyright      string                     `json:"copyright"`
	Note           string                     `json:"note"`
	ModifiedDate   string                     `json:"modified_date"`
	ModifiedPid    int                        `json:"modified_pid,string"`
	Tag            []string                   `json:"tag"`
	Arrangements   map[string]SongArrangement `json:"arrangement"`
}

type songsdata struct {
	Songs map[string]Song `json:"songs"`
}

type apiResponse struct {
	Status string    `json:"status"`
	Data   songsdata `json:"data"`
}