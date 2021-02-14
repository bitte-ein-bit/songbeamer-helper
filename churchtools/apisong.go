package churchtools

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

func (s *APISong) GetDefaultArrangement() (ret APISongArrangement) {
	for _, value := range s.Arrangements {
		if value.Default {
			ret = value
		}
	}
	return
}
