package churchtools

type apiAgendaResponse struct {
	Agenda Agenda `json:"data"`
}

type Agenda struct {
	ID      int          `json:"id"`
	Name    string       `json:"name"`
	IsFinal bool         `json:"isFinal"`
	Items   []AgendaItem `json:"items"`
}

type AgendaItem struct {
	ID            int    `json:"id"`
	Position      int    `json:"position"`
	Type          string `json:"type"`
	Title         string `json:"title"`
	IsBeforeEvent bool   `json:"isBeforeEvent"`
	SongID        int    `json:"songId,omitempty"`
	ArrangementID int    `json:"arrangementId,omitempty"`
}
