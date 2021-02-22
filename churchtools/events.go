package churchtools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

// GetEvents returns the next events from today until daysInFuture days
func GetEvents(daysInFuture int) []Event {
	if client == nil {
		login()
	}
	params := map[string]string{
		"from": time.Now().Format("2006-01-02"),
		"to":   time.Now().AddDate(0, 0, daysInFuture).Format("2006-01-02"),
	}
	fmt.Println(params)
	resp := getRequest(fmt.Sprintf("https://%s/api/events", domain), params)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	r := getEventsResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	return r.Data
}
