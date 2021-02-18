package churchtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const domain = "lkg-pfuhl.church.tools"

var basisURL = fmt.Sprintf("https://%s/?q=", domain)
var churchServiceAjaxURL = "https://lkg-pfuhl.church.tools/?q=churchservice/ajax"
var churchServiceFiledownloadURL = "https://lkg-pfuhl.church.tools/?q=churchservice/filedownload"

const userid = "2392"
const token = "23bwRElUXrBXmriaIrMP8vrJAxoIcJH9KJGfTLsEHpusNqnLnnwTBWLLbdjzKilg3Ns0vxZB8SCGATeGc3D8zIgEiqjtpP1VHo64vO9fjFvGcb2wueQETwI8a3w6kWdOoNdR3ZPzm0G50HOczY2AOILkA0fxlb1sboiLPcvNEYuRuCHe3kKe9TFOloFSQLvBrYrRdag0C6qpd3A9YW4XW4byQjsGOKhhPCgoA54nDwHoauLtS8hKD2XdSq9i6sPA"

var csrftokens map[string]string
var client *http.Client

func setupClient() {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client = &http.Client{Jar: jar}
}

func login() {
	if client == nil {
		setupClient()
	}
	loginurl := fmt.Sprintf("%slogin/ajax", basisURL)
	resp, err := client.PostForm(loginurl, url.Values{
		"func":       {"loginWithToken"},
		"id":         {userid},
		"token":      {token},
		"directtool": {"songsync"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Login successful")
}

type csrftoken struct {
	Token string `json:"data"`
}

func getCSRFToken() string {
	if csrftokens == nil {
		csrftokens = make(map[string]string)
	}
	if val, ok := csrftokens[domain]; ok {
		return val
	}
	url := fmt.Sprintf("https://%s/api/csrftoken", domain)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	csrftoken1 := csrftoken{}
	jsonErr := json.Unmarshal(data, &csrftoken1)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	log.Printf("Token: %s\n", csrftoken1.Token)
	csrftokens[domain] = csrftoken1.Token
	return csrftoken1.Token
}

func getRequest(url string, params map[string]string) http.Response {
	if client == nil {
		log.Fatal("please login first")
	}
	req, _ := http.NewRequest("GET", url, nil)
	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("CSRF-Token", getCSRFToken())
	log.Println(req.Header)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.Header)
	return *resp
}

func postRequest(client *http.Client, url string, params map[string]string) http.Response {
	if client == nil {
		log.Fatal("please login first")
	}
	req, _ := http.NewRequest("POST", url, nil)
	if params != nil {
		fmt.Println(params)
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("CSRF-Token", getCSRFToken())
	log.Println(req.Header)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.Header)
	return *resp
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// NewfileUploadRequest Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path, contentType string, uploadName ...string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileName := fi.Name()

	if len(uploadName) > 0 {
		fileName = uploadName[0]
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(paramName), escapeQuotes(fileName)))
	h.Set("Content-Type", contentType)
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("CSRF-Token", getCSRFToken())
	req.Header.Add("Content-Type", writer.FormDataContentType())
	log.Println(req.Header)
	// log.Println(body)
	return req, err

}
