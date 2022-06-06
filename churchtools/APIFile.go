package churchtools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/bitte-ein-bit/songbeamer-helper/log"
)

// APIFile describes a file as given by the new ChurchTools REST API
type APIFile struct {
	DomainType string `json:"domainType"`
	DomainID   int    `json:"domainID,string"`
	filepath   string
	Name       string `json:"name"`
	Filename   string `json:"filename"`
	FileURL    string `json:"fileUrl,omitempty"`
	uploadName string
}

func (f *APIFile) getID() int {
	if f.FileURL == "" {
		return 0
	}
	params, err := url.ParseQuery(f.FileURL)
	if err != nil {
		log.Fatalf("Cannot parse FileURL into segments: %s", err)
		return 0
	}
	for key, value := range params {
		if key == "id" {
			if n, err := strconv.Atoi(value[0]); err == nil {
				return n
			}
		}
	}
	return 0
}

// NewAPIFile generates a new APIFile struct
func NewAPIFile(path string) *APIFile {
	return &APIFile{
		filepath: path,
		Name:     filepath.Base(path),
	}
}

// NewSongAPIFile returns an song_arrangement type API file
func NewSongAPIFile(path string, domainID int) *APIFile {
	return &APIFile{
		filepath:   path,
		Name:       filepath.Base(path),
		DomainType: "song_arrangement",
		DomainID:   domainID,
	}
}

// SetUploadName overrides the automatically selected name based on the local filename with a defined value
func (f *APIFile) SetUploadName(name string) {
	f.uploadName = name
}

func (f *APIFile) getUploadName() string {
	if f.uploadName != "" {
		return f.uploadName
	}
	return filepath.Base(f.filepath)
}

// Save submits the file to ChurchTools
func (f *APIFile) Save() error {
	// AddSongFile Upload and attach a file to a ChurchTools song
	if f.DomainType == "" {
		return fmt.Errorf("Please set DomainType")
	}
	if f.DomainID == 0 {
		return fmt.Errorf("Please set DomainID")
	}
	currentID := f.getID()
	url := fmt.Sprintf("https://%s/api/files/%s/%d", domain, f.DomainType, f.DomainID)
	request, err := newfileUploadRequest(url, nil, "files[]", f.filepath, "text/plain", f.getUploadName())
	if err != nil {
		return fmt.Errorf("Creating request failed: %w", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Save failed: %w", err)
	}
	var bodyContent []byte
	// log.Infof(resp.StatusCode)
	// log.Infof(resp.Header)
	defer resp.Body.Close()
	resp.Body.Read(bodyContent)
	// log.Infof(bodyContent)
	if currentID != 0 {
		log.Debugf("File has been updated, deleting old version")
		f.Delete(currentID)
	}
	return nil
}

// LoadFromFile sets the name as well as the filepath attribute
func (f *APIFile) LoadFromFile(path string) {
	f.filepath = path
	f.Name = filepath.Base(path)
}

// Delete removes a file by ID from ChurchTools
func (f APIFile) Delete(ID int) error {
	if ID == 0 {
		return fmt.Errorf("Cannot delete file with ID 0")
	}
	params := map[string]string{
		"func": "delFile",
		"id":   fmt.Sprintf("%d", ID),
	}
	resp := postRequest(client, churchServiceAjaxURL, params)
	// log.Println(resp.Status)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%s", err)
	}
	// log.Println(string(data))
	r := songResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		return fmt.Errorf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	if r.Status != "success" {
		return fmt.Errorf("Cannot edit arrangement: %s", r.Message)
	}
	return nil
}
