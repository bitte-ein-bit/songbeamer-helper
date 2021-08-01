package churchtools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/log"
)

// GetEvents returns the next events from today until daysInFuture days
func GetEvents(daysInFuture int) []Event {
	if client == nil {
		login()
	}
	params := map[string]string{
		// "from": time.Now().AddDate(0, 0, -14).Format("2006-01-02"),
		"from": time.Now().Format("2006-01-02"),
		"to":   time.Now().AddDate(0, 0, daysInFuture).Format("2006-01-02"),
	}
	log.Debugf("%s", params)
	resp := getRequest(fmt.Sprintf("https://%s/api/events", domain), params)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("%s", err)
	}
	r := getEventsResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	return r.Data
}
