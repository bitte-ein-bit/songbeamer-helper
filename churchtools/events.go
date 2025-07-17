package churchtools

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/log"
)

// GetEvents returns the next events from today until daysInFuture days
func GetEvents(client ChurchToolsClient, daysInFuture, daysInPast int) []Event {
	params := map[string]string{
		"from": time.Now().AddDate(0, 0, -daysInPast).Format("2006-01-02"),
		"to":   time.Now().AddDate(0, 0, daysInFuture).Format("2006-01-02"),
	}
	log.Debugf("%s", params)
	url := fmt.Sprintf("https://%s/api/events", domain)
	resp := client.GetRequest(url, params)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

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
