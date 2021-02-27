package churchtools

type apiAgendaResponse struct {
	Agenda APIAgenda `json:"data"`
}

// APIAgenda describes the agenda as returned by the REST API
type APIAgenda struct {
	ID      int             `json:"id"`
	Name    string          `json:"name"`
	IsFinal bool            `json:"isFinal"`
	Items   []APIAgendaItem `json:"items"`
}

// APIAgendaItem describes an item in APIAgenda
type APIAgendaItem struct {
	ID            int    `json:"id"`
	Position      int    `json:"position"`
	Type          string `json:"type"`
	Title         string `json:"title"`
	IsBeforeEvent bool   `json:"isBeforeEvent"`
	SongID        int    `json:"songId,omitempty"`
	ArrangementID int    `json:"arrangementId,omitempty"`
}
