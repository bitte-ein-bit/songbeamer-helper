package churchtools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	r := apiResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		return fmt.Errorf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	if r.Status != "success" {
		return fmt.Errorf("Cannot delete song: %s", r.Message)
	}
	return nil
}

func (s *Song) GetDefaultArrangement() (ret SongArrangement) {
	for _, value := range s.Arrangements {
		if value.Default == 1 {
			ret = value
		}
	}
	return
}
