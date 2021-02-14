package churchtools

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
