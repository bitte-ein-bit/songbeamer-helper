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
	Note         string              `json:"note,omitempty"`
	ModifiedDate string              `json:"modified_date"`
	ModifiedPID  int                 `json:"modified_pid,string"`
	Files        map[string]SongFile `json:"files,omitempty"`
}

// The Song is describing the api response of a Song. It is further defined by one or more Arrangements
type Song struct {
	ID             int                        `json:"id,string"`
	Bezeichnung    string                     `json:"bezeichnung"`
	SongcategoryID int                        `json:"songcategory_id,string"`
	Practice       int                        `json:"practice_yn,string"`
	Author         string                     `json:"author"`
	CCLI           string                     `json:"ccli,omitempty"`
	Copyright      string                     `json:"copyright"`
	Note           string                     `json:"note"`
	ModifiedDate   string                     `json:"modified_date"`
	ModifiedPid    int                        `json:"modified_pid,string"`
	Tag            []string                   `json:"tag"`
	Arrangements   map[string]SongArrangement `json:"arrangement"`
}

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

// APISongCategory represents a categorie
type APISongCategory struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	NameTranslated string `json:"nameTranslated"`
	SortKey        int    `json:"sortKey"`
	CampusID       int    `json:"campusID"`
}

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

type songsdata struct {
	Songs map[string]Song `json:"songs"`
}

// data := []byte(`{
//     "id": "45",
//     "bezeichnung": "Wie weit würd ich gehn",
//     "songcategory_id": "0",
//     "practice_yn": "0",
//     "author": "Arne Kopfermann, Benjamin Heinrich",
//     "ccli": "7096862",
//     "copyright": "2017 SCM Hänssler, Holzgerlingen (Verwaltet von SCM Hänssler)",
//     "note": "",
//     "modified_date": "2021-01-31 11:18:06",
//     "modified_pid": "279",
//     "arrangement": {
//       "48": {
//         "id": "48",
//         "bezeichnung": "Standard-Arrangement",
//         "default_yn": "1",
//         "tonality": "",
//         "bpm": "",
//         "beat": "",
//         "length_min": "0",
//         "length_sec": "0",
//         "note": null,
//         "modified_date": "2021-01-31 11:18:06",
//         "modified_pid": "279",
//         "files": {
//           "1884": {
//             "id": "1884",
//             "domain_type": "song_arrangement",
//             "domain_id": "48",
//             "bezeichnung": "Wie weit würd ich gehn.txt",
//             "filename": "04fa8dc5201c3b7c7860e6d946f6b9be.txt",
//             "showonlywheneditable_yn": "0",
//             "securitylevel_id": null,
//             "image_options": null,
//             "modified_date": "2021-01-31 11:18:11",
//             "modified_pid": "279",
//             "deletion_date": null,
//             "modified_username": "Benjamin Böttinger Admin"
//           },
//           "1887": {
//             "id": "1887",
//             "domain_type": "song_arrangement",
//             "domain_id": "48",
//             "bezeichnung": "Wie weit würd ich gehn.sng",
//             "filename": "fa028ad85c298e0efade2bad6991dee9.sng",
//             "showonlywheneditable_yn": "0",
//             "securitylevel_id": null,
//             "image_options": null,
//             "modified_date": "2021-01-31 11:18:11",
//             "modified_pid": "279",
//             "deletion_date": null,
//             "modified_username": "Benjamin Böttinger Admin"
//           }
//         }
//       }
//     },
//     "tags": []
//   }
// `)
type apiResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
	Data    songsdata `json:"data"`
}

type addResponse struct {
	Status string `json:"status"`
	ID     int    `json:"data,string"`
}

type getSongResponse struct {
	Data APISong `json:"data"`
}
