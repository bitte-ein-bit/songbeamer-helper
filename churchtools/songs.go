package churchtools

import (
	"encoding/json"
	"fmt"
	"io"
	"math"

	"github.com/bitte-ein-bit/songbeamer-helper/log"
)

const (
	maxSize = 200*1024 // Maximum size of a message before truncating
)

func truncateMessage(message string) string {
    if len(message) <= maxSize {
        return message
    }
    return message[:maxSize] + "... [truncated]"
}

// GetSongs returns the Songs as sent by churchservice/getAllSongs endpoint
func GetSongs() (map[string]Song, error) {
	log.Debugf("Enter GetSongs")
	if client == nil {
		login()
	}
	params := make(map[string]string)
	params["func"] = "getAllSongs"
	resp := getRequest(churchServiceAjaxURL, params)
	log.Debugf("Response Status: %s", resp.Status)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	r := songResponse{}

	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		println(string(data))
		log.Debugf("Data: %s",truncateMessage(string(data)))
		return nil, fmt.Errorf("unable to parse value, error: %s", jsonErr.Error())
	}
	log.Debugf("GetSongs: %d songs found", len(r.Data.Songs))
	return r.Data.Songs, nil
}

// AddSong adds a new song to Churchtools
func AddSong(bezeichnung, author, copyright, ccli, tonality, bpm, beat, songCat string) int {
	if client == nil {
		login()
	}

	params := make(map[string]string)
	params["func"] = "addNewSong"
	params["bezeichnung"] = bezeichnung
	params["author"] = author
	params["copyright"] = copyright
	params["ccli"] = ccli
	params["tonality"] = tonality
	params["bpm"] = bpm
	params["beat"] = beat
	params["songcategory_id"] = songCat
	params["comments[domain_type]"] = "arrangement"
	resp := postRequest(client, churchServiceAjaxURL, params)
	log.Println(resp.Status)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	r := addResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	return r.ID
}

// GetSong loads a song from ChurchTools
func GetSong(songID int) APISong {
	if client == nil {
		login()
	}

	url := fmt.Sprintf("https://%s/api/songs/%d", domain, songID)
	resp := getRequest(url, nil)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
	r := getSongResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	log.Println(r)
	return r.Data
}

// EditArrangement change the details of an arrangement
func EditArrangement(arrangement APISongArrangement, songID int) {
	if client == nil {
		login()
	}
	params := make(map[string]string)
	params["func"] = "editArrangement"
	params["bezeichnung"] = arrangement.Name
	params["length_min"] = fmt.Sprintf("%0.0f", math.Floor(float64(arrangement.Duration)/60))
	params["length_sec"] = fmt.Sprintf("%d", arrangement.Duration%60)
	params["tonality"] = arrangement.KeyOfArrangement
	params["bpm"] = arrangement.BPM
	params["beat"] = arrangement.Beat
	params["note"] = arrangement.Note
	params["song_id"] = fmt.Sprintf("%d", songID)
	params["id"] = fmt.Sprintf("%d", arrangement.ID)
	resp := postRequest(client, churchServiceAjaxURL, params)
	log.Println(resp.Status)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
	r := songResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	if r.Status != "success" {
		log.Fatalf("Cannot edit arrangement: %s", r.Message)
	}
}
