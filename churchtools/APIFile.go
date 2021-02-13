package churchtools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
)

type APIFile struct {
	DomainType string `json:"domainType"`
	DomainID   int    `json:"domainID,string"`
	filepath   string
	client     *http.Client
	Name       string `json:"name"`
	Filename   string `json:"filename"`
	FileURL    string `json:"fileUrl,omitempty"`
}

func (f *APIFile) getID() int {
	if f.FileURL == "" {
		return 0
	}
	params, err := url.ParseQuery(f.FileURL)
	if err != nil {
		log.Fatal(err)
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

func NewAPIFile(path string) (*APIFile, error) {
	f := &APIFile{
		filepath: path,
		client:   client,
		Name:     filepath.Base(path),
	}
	return f, nil
}

func (f *APIFile) Save() error {
	// AddSongFile Upload and attach a file to a ChurchTools song
	if f.DomainType == "" {
		return fmt.Errorf("Please set DomainType")
	}
	if f.DomainID == 0 {
		return fmt.Errorf("Please set DomainID")
	}
	if f.client == nil {
		// return fmt.Errorf("Client is uninitialized")
		f.client = client
	}
	currentID := f.getID()
	url := fmt.Sprintf("https://%s/api/files/%s/%d", domain, f.DomainType, f.DomainID)
	request, err := newfileUploadRequest(url, nil, "files[]", f.filepath, "text/plain")

	resp, err := f.client.Do(request)
	if err != nil {
		return fmt.Errorf("Save failed: %w", err)
	}
	var bodyContent []byte
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	resp.Body.Read(bodyContent)
	resp.Body.Close()
	fmt.Println(bodyContent)
	if currentID != 0 {
		f.Delete(currentID)
	}
	return nil
}

func (f *APIFile) LoadFromFile(path string) {
	f.filepath = path
	f.Name = filepath.Base(path)
}

func (f *APIFile) SetHTTPClient(client *http.Client) {
	f.client = client
}

func (f APIFile) Delete(ID int) error {
	if ID == 0 {
		return fmt.Errorf("Cannot delete file with ID 0")
	}
	if f.client == nil {
		// return fmt.Errorf("Client is uninitialized")
		f.client = client
	}
	params := map[string]string{
		"func": "delFile",
		"id":   fmt.Sprintf("%d", ID),
	}
	resp := postRequest(f.client, churchServiceAjaxURL, params)
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
		return fmt.Errorf("Cannot edit arrangement: %s", r.Message)
	}
	return nil
}
