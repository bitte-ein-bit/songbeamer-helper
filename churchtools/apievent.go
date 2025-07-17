package churchtools

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/log"
)

type getEventsResponse struct {
	Data []Event `json:"data"`
}

// Event describes an event as returned by the REST api
type Event struct {
	ID          int       `json:"id"`
	EndDate     time.Time `json:"endDate"`
	StartDate   time.Time `json:"startDate"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// GetSongs returns all Songs the user has access to that are part of this event
func (e *Event) GetSongs(client ChurchToolsClient) []APISong {
	log.Debugf("GetSongs for event %d", e.ID)
	if e.ID == 0 {
		log.Fatal("Cannot load songs for uninitialzed event")
	}

	resp := client.GetRequest(fmt.Sprintf("https://%s/api/events/%d/agenda/songs", domain, e.ID), nil)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	r := getSongsResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	return r.Data
}

// GetAgenda retuns the Agenda for the Event
func (e *Event) GetAgenda(client ChurchToolsClient) APIAgenda {
	if e.ID == 0 {
		log.Fatal("Cannot load songs for uninitialzed event")
	}

	// TODO: no Agenda setup yet...
	resp := client.GetRequest(fmt.Sprintf("https://%s/api/events/%d/agenda", domain, e.ID), nil)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	r := apiAgendaResponse{}
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(data), jsonErr.Error())
	}
	return r.Agenda
}
