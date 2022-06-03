package churchtools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/log"
)

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

// Delete deletes the song from ChurchTools
func (s *Song) Delete() error {
	if s.ID == 0 {
		return fmt.Errorf("Cannot delete file with ID 0")
	}
	params := map[string]string{
		"func": "deleteSong",
		"id":   fmt.Sprintf("%d", s.ID),
	}
	resp := postRequest(client, churchServiceAjaxURL, params)
	log.Println(resp.Status)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
	r := songResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		return fmt.Errorf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	if r.Status != "success" {
		return fmt.Errorf("Cannot delete song: %s", r.Message)
	}
	return nil
}

// GetDefaultArrangement retrieves the arrangement marked as default
func (s *Song) GetDefaultArrangement() (ret SongArrangement) {
	for _, value := range s.Arrangements {
		if value.Default == 1 {
			ret = value
		}
	}
	return
}

// AddArrangement adds an additional arrangement to an existing song on ChurchTools. It does return the ID of the new arrangement
func (s *Song) AddArrangement(name string) (int, error) {
	params := map[string]string{
		"func":        "addArrangement",
		"bezeichnung": name,
		"song_id":     fmt.Sprintf("%d", s.ID),
	}
	resp := postRequest(client, churchServiceAjaxURL, params)
	log.Println(resp.Status)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
	r := addResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		return 0, fmt.Errorf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	if r.Status != "success" {
		return 0, fmt.Errorf("Cannot add arrangement: %s", r.Message)
	}
	s.Arrangements[fmt.Sprintf("%d", r.ID)] = SongArrangement{
		Bezeichnung: name,
		ID:          r.ID,
	}
	return r.ID, nil
}

// GetModificationDate parses the date string returned by the API into a time struct
func (s *Song) GetModificationDate() (t time.Time, err error) {
	layout := "2006-01-02 15:04:05"
	t, err = time.Parse(layout, s.ModifiedDate)
	return
}
