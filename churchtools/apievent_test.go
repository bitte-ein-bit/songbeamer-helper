package churchtools_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEvent_GetSongs(t *testing.T) {
	t.Run("should fetch songs for an event", func(t *testing.T) {
		mockClient := new(MockChurchToolsClient)
		event := churchtools.Event{
			ID:   123,
			Name: "Sunday Service",
		}

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"data": [
					{
						"id": 456,
						"name": "Amazing Grace",
						"author": "John Newton",
						"ccli": "12345",
						"category": {
							"id": 1,
							"name": "Worship"
						},
						"arrangements": [
							{
								"id": 1,
								"name": "Standard",
								"isDefault": true
							}
						]
					},
					{
						"id": 789,
						"name": "How Great Thou Art",
						"author": "Carl Boberg",
						"arrangements": []
					}
				]
			}`)),
		}

		expectedURL := "https://lkg-pfuhl.church.tools/api/events/123/agenda/songs"
		mockClient.On("GetRequest", expectedURL, mock.Anything).Return(mockResponse)

		songs := event.GetSongs(mockClient)

		assert.Len(t, songs, 2)
		assert.Equal(t, 456, songs[0].ID)
		assert.Equal(t, "Amazing Grace", songs[0].Bezeichnung)
		assert.Equal(t, "John Newton", songs[0].Author)
		assert.Len(t, songs[0].Arrangements, 1)
		assert.Equal(t, 789, songs[1].ID)
		assert.Equal(t, "How Great Thou Art", songs[1].Bezeichnung)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty songs list", func(t *testing.T) {
		mockClient := new(MockChurchToolsClient)
		event := churchtools.Event{
			ID:   123,
			Name: "Empty Event",
		}

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{"data": []}`)),
		}

		mockClient.On("GetRequest", "https://lkg-pfuhl.church.tools/api/events/123/agenda/songs",
			mock.Anything).Return(mockResponse)

		songs := event.GetSongs(mockClient)

		assert.Empty(t, songs)

		mockClient.AssertExpectations(t)
	})
}

func TestEvent_GetAgenda(t *testing.T) {
	t.Run("should fetch agenda for an event", func(t *testing.T) {
		mockClient := new(MockChurchToolsClient)
		event := churchtools.Event{
			ID:   123,
			Name: "Sunday Service",
		}

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"data": {
					"id": 456,
					"name": "Service Agenda",
					"isFinal": true,
					"items": [
						{
							"id": 1,
							"position": 1,
							"type": "song",
							"title": "Opening Song",
							"isBeforeEvent": false,
							"song": {
								"songId": 789,
								"arrangementId": 1,
								"title": "Amazing Grace",
								"arrangement": "Standard",
								"category": "Worship",
								"key": "G",
								"bpm": 120,
								"isDefault": true
							}
						},
						{
							"id": 2,
							"position": 2,
							"type": "normal",
							"title": "Sermon",
							"isBeforeEvent": false
						}
					]
				}
			}`)),
		}

		expectedURL := "https://lkg-pfuhl.church.tools/api/events/123/agenda"
		mockClient.On("GetRequest", expectedURL, mock.Anything).Return(mockResponse)

		agenda := event.GetAgenda(mockClient)

		assert.Equal(t, 456, agenda.ID)
		assert.Equal(t, "Service Agenda", agenda.Name)
		assert.True(t, agenda.IsFinal)
		assert.Len(t, agenda.Items, 2)
		assert.Equal(t, "song", agenda.Items[0].Type)
		assert.Equal(t, "Amazing Grace", agenda.Items[0].Song.Title)
		assert.Equal(t, "normal", agenda.Items[1].Type)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty agenda", func(t *testing.T) {
		mockClient := new(MockChurchToolsClient)
		event := churchtools.Event{
			ID:   123,
			Name: "Empty Event",
		}

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"data": {
					"id": 0,
					"name": "",
					"isFinal": false,
					"items": []
				}
			}`)),
		}

		mockClient.On("GetRequest", "https://lkg-pfuhl.church.tools/api/events/123/agenda",
			mock.Anything).Return(mockResponse)

		agenda := event.GetAgenda(mockClient)

		assert.Equal(t, 0, agenda.ID)
		assert.Empty(t, agenda.Items)
		assert.False(t, agenda.IsFinal)

		mockClient.AssertExpectations(t)
	})
}
