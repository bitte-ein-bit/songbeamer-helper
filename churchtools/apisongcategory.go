package churchtools

// APISongCategory represents a categorie
type APISongCategory struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	NameTranslated string `json:"nameTranslated"`
	SortKey        int    `json:"sortKey"`
	CampusID       int    `json:"campusID"`
}
