package churchtools

import (
	"fmt"
	"time"
)

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

// ToAPIFile converts an old school file to the new rest api variant
func (s *SongFile) ToAPIFile() APIFile {
	a := APIFile{
		DomainType: s.DomainType,
		DomainID:   s.DomainID,
		Name:       s.Bezeichnung,
		Filename:   s.Filename,
		FileURL:    fmt.Sprintf("https://lkg-pfuhl.church.tools/?q=public/filedownload&id=%d&filename=%s", s.ID, s.Filename),
		uploadName: s.Filename,
	}
	return a
}

// GetModificationDate parses the date string returned by the API into a time struct
func (s *SongFile) GetModificationDate() (t time.Time, err error) {
	layout := "2006-01-02 15:04:05 -0700"
	withTZ := fmt.Sprintf("%s +0100", s.ModifiedDate)
	t, err = time.Parse(layout, withTZ)
	return
}
