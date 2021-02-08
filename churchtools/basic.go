package churchtools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/publicsuffix"
)

const domain = "lkg-pfuhl.church.tools"
const loginurl = "https://lkg-pfuhl.church.tools/?q=login/ajax"
const churchServiceURL = "https://lkg-pfuhl.church.tools/?q=churchservice/ajax"

// These values are injected at build time via -ldflags
var userid = ""
var token = ""

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
