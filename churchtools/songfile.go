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

// ToAPIFile converts an old school file to the new rest api variant
func (s *SongFile) ToAPIFile() APIFile {
	a := APIFile{
		DomainType: s.DomainType,
		DomainID:   s.DomainID,
		Name:       s.Bezeichnung,
		Filename:   s.Filename,
		FileURL:    "",
		uploadName: s.Filename,
	}
	return a
}
