package churchtools

type songsdata struct {
	Songs map[string]Song `json:"songs"`
}

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
type apiResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
	Data    songsdata `json:"data"`
}

type addResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	ID      int    `json:"data,string"`
}

type getSongResponse struct {
	Data APISong `json:"data"`
}
