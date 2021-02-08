package churchtools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func getSongs() error {
	if client == nil {
		login()
	}
	params := make(map[string]string)
	params["func"] = "getAllSongs"
	resp := getRequest(churchServiceAjaxURL, params)
	log.Println(resp.Status)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	r := apiResponse{}
	// data := []byte(`{
	//     "id": "45",
	//     "bezeichnung": "Wie weit würd ich gehn",
	//     "songcategory_id": "0",
	//     "practice_yn": "0",
	//     "author": "Arne Kopfermann, Benjamin Heinrich",
	//     "ccli": "7096862",
	//     "copyright": "2017 SCM Hänssler, Holzgerlingen (Verwaltet von SCM Hänssler)",
	//     "note": "",
	//     "modified_date": "2021-01-31 11:18:06",
	//     "modified_pid": "279",
	//     "arrangement": {
	//       "48": {
	//         "id": "48",
	//         "bezeichnung": "Standard-Arrangement",
	//         "default_yn": "1",
	//         "tonality": "",
	//         "bpm": "",
	//         "beat": "",
	//         "length_min": "0",
	//         "length_sec": "0",
	//         "note": null,
	//         "modified_date": "2021-01-31 11:18:06",
	//         "modified_pid": "279",
	//         "files": {
	//           "1884": {
	//             "id": "1884",
	//             "domain_type": "song_arrangement",
	//             "domain_id": "48",
	//             "bezeichnung": "Wie weit würd ich gehn.txt",
	//             "filename": "04fa8dc5201c3b7c7860e6d946f6b9be.txt",
	//             "showonlywheneditable_yn": "0",
	//             "securitylevel_id": null,
	//             "image_options": null,
	//             "modified_date": "2021-01-31 11:18:11",
	//             "modified_pid": "279",
	//             "deletion_date": null,
	//             "modified_username": "Benjamin Böttinger Admin"
	//           },
	//           "1887": {
	//             "id": "1887",
	//             "domain_type": "song_arrangement",
	//             "domain_id": "48",
	//             "bezeichnung": "Wie weit würd ich gehn.sng",
	//             "filename": "fa028ad85c298e0efade2bad6991dee9.sng",
	//             "showonlywheneditable_yn": "0",
	//             "securitylevel_id": null,
	//             "image_options": null,
	//             "modified_date": "2021-01-31 11:18:11",
	//             "modified_pid": "279",
	//             "deletion_date": null,
	//             "modified_username": "Benjamin Böttinger Admin"
	//           }
	//         }
	//       }
	//     },
	//     "tags": []
	//   }
	// `)
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	// log.Println(r)
	for s := range r.Data.Songs {
		log.Printf("[%05d] %s", r.Data.Songs[s].ID, string(r.Data.Songs[s].Bezeichnung))
	}

	return nil
}

// Songs returns the Songs as sent by churchservice/getAllSongs endpoint
func Songs() {
	getSongs()
}

// AddSong adds a new song to Churchtools
func AddSong(bezeichnung, author, copyright, ccli, tonality, bpm, beat string) int {
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
	params["songcategory_id"] = "1"
	params["comments[domain_type]"] = "arrangement"
	resp := postRequest(churchServiceAjaxURL, params)
	log.Println(resp.Status)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

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

// AddSongFile Upload and attach a file to a ChurchTools song
func AddSongFile(arrangementID int, filepath string) {
	if client == nil {
		login()
	}
	url := fmt.Sprintf("https://%s/api/files/song_arrangement/%d", domain, arrangementID)
	request, err := newfileUploadRequest(url, nil, "files[]", filepath, "text/plain")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		var bodyContent []byte
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		resp.Body.Read(bodyContent)
		resp.Body.Close()
		fmt.Println(bodyContent)
	}
}
