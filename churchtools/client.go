package churchtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"os"
	"strings"

	"github.com/bitte-ein-bit/songbeamer-helper/log"
	"golang.org/x/net/publicsuffix"
)

// CTClient abstracts a little the HTTP handling for ChurchTools
type CTClient struct {
	Client       *http.Client
	quoteEscaper *strings.Replacer
	csrftokens   map[string]string
}

func (c *CTClient) setupClient() {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	c.Client = &http.Client{Jar: jar}
	c.quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
	c.csrftokens = make(map[string]string)
}

// Login logs into ChurchTools
func (c *CTClient) Login() error {
	if c.Client == nil {
		c.setupClient()
	}
	loginurl := fmt.Sprintf("%slogin/ajax", basisURL)
	resp, err := c.Client.PostForm(loginurl, url.Values{
		"func":       {"loginWithToken"},
		"id":         {userid},
		"token":      {token},
		"directtool": {"songsync"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Login successful")
	return nil
}

func (c *CTClient) getCSRFToken() string {
	log.Debugf("Get CSRF Token")
	if c.csrftokens == nil {
		c.csrftokens = make(map[string]string)
	}
	if val, ok := c.csrftokens[domain]; ok {
		return val
	}
	url := fmt.Sprintf("https://%s/api/csrftoken", domain)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	resp, _ := c.Client.Do(req)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	csrftoken1 := csrftoken{}
	jsonErr := json.Unmarshal(data, &csrftoken1)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	log.Debugf("Token: %s\n", csrftoken1.Token)
	c.csrftokens[domain] = csrftoken1.Token
	return csrftoken1.Token
}

// GetRequest executes a HTTP Get request
func (c *CTClient) GetRequest(url string, params map[string]string) *http.Response {
	req, _ := http.NewRequest("GET", url, nil)
	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("CSRF-Token", c.getCSRFToken())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "Login "+token)
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

// PostRequest executes a HTTP Post request
func (c *CTClient) PostRequest(url string, params map[string]string) *http.Response {
	req, _ := http.NewRequest("POST", url, nil)
	if params != nil {
		// log.Infof(params)
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("CSRF-Token", c.getCSRFToken())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "Login "+token)
	// log.Println(req.Header)
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(resp.Header)
	return resp
}

// DeleteRequest executes a HTTP Delete request
func (c *CTClient) DeleteRequest(url string, params map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest("DELETE", url, nil)
	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("CSRF-Token", c.getCSRFToken())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "Login "+token)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("delete request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to delete file (status %d): %s", resp.StatusCode, truncateMessage(string(bodyBytes)))
	}
	return resp, nil
}

func (c *CTClient) escapeQuotes(s string) string {
	return c.quoteEscaper.Replace(s)
}

// NewfileUploadRequest Creates a new file upload http request with optional extra params
func (c *CTClient) NewfileUploadRequest(uri string, params map[string]string, paramName, path, contentType string, uploadName ...string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileContents, err := io.ReadAll(file)
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
			c.escapeQuotes(paramName), c.escapeQuotes(fileName)))
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
	req.Header.Set("CSRF-Token", c.getCSRFToken())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "Login "+token)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	log.Println(req.Header)
	// log.Println(body)
	return req, err

}

// ChurchToolsClient defines the interface for interacting with ChurchTools
type ChurchToolsClient interface {
	GetRequest(url string, params map[string]string) *http.Response
	PostRequest(url string, params map[string]string) *http.Response
	DeleteRequest(url string, params map[string]string) (*http.Response, error)
	// Login() error
}

// Ensure CTClient implements ChurchToolsClient
var _ ChurchToolsClient = &CTClient{}

// SetGlobalClient sets the global CTClient instance for use across the package
func SetGlobalClient(c *CTClient) {
	globalCTClient = c
}

// NewClient returns a new initialized ChurchToolsClient
func NewClient() ChurchToolsClient {
	c := &CTClient{}
	c.setupClient()
	c.Login()
	SetGlobalClient(c)
	return c
}
