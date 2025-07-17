package churchtools_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChurchToolsClient struct {
	mock.Mock
}

func (m *MockChurchToolsClient) GetRequest(url string, params map[string]string) *http.Response {
	args := m.Called(url, params)
	return args.Get(0).(*http.Response)
}

func (m *MockChurchToolsClient) Login() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockChurchToolsClient) DeleteRequest(url string, params map[string]string) (*http.Response, error) {
	args := m.Called(url, params)
	return args.Get(0).(*http.Response), nil
}
func (m *MockChurchToolsClient) PostRequest(url string, params map[string]string) *http.Response {
	args := m.Called(url, params)
	return args.Get(0).(*http.Response)
}


func TestGetEvents(t *testing.T) {

	t.Run("should fetch events for the specified date range", func(t *testing.T) {
		// Create mock client
		mockClient := new(MockChurchToolsClient)

		// Mock response
		// Get yesterday's date
		yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

		// create a time for tomorrow, 16:00 UTC
		tomorrow := time.Now().AddDate(0, 0, 1)
		tomorrow2 := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 16, 0, 0, 0, time.UTC)

		tomorrowStr := tomorrow2.Format("2006-01-02")

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"data": [
					{
						"id": 123,
						"name": "Sunday Service",
						"startDate": "` + yesterdayStr + `T15:00:00Z",
						"endDate": "` + yesterdayStr + `T16:00:00Z"
					},
					{
						"id": 124,
						"name": "Prayer Meeting",
						"startDate": "` + tomorrowStr + `T16:00:00Z",
						"endDate": "` + tomorrowStr + `T17:00:00Z"
					}
				]
			}`)),
		}

		// Setup expectations
		expectedURL := "https://lkg-pfuhl.church.tools/api/events"
		today := time.Now()
		expectedFromDate := today.AddDate(0, 0, -2).Format("2006-01-02")
		expectedToDate := today.AddDate(0, 0, 7).Format("2006-01-02")
		expectedParams := map[string]string{
			"from": expectedFromDate,
			"to":   expectedToDate,
		}
		mockClient.On("GetRequest", expectedURL, expectedParams).Return(mockResponse)

		// Call function under test
		events := churchtools.GetEvents(mockClient, 7, 2)

		// Assertions
		assert.Len(t, events, 2)
		assert.Equal(t, 123, events[0].ID)
		assert.Equal(t, "Sunday Service", events[0].Name)
		assert.Equal(t, yesterdayStr, events[0].StartDate.Format("2006-01-02"))
		assert.Equal(t, 124, events[1].ID)
		assert.Equal(t, "Prayer Meeting", events[1].Name)
		assert.Equal(t, tomorrow2, events[1].StartDate)

		// Verify expectations
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty events response", func(t *testing.T) {
		// Create mock client
		mockClient := new(MockChurchToolsClient)

		// Mock response
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{"data": []}`)),
		}

		// Setup expectations
		mockClient.On("GetRequest", mock.Anything, mock.Anything).Return(mockResponse)

		// Call function under test
		events := churchtools.GetEvents(mockClient, 7, 0)

		// Assertions
		assert.Empty(t, events)

		// Verify expectations
		mockClient.AssertExpectations(t)
	})
}
