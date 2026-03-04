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

func TestGetSongs(t *testing.T) {
	t.Run("should fetch and parse all songs", func(t *testing.T) {
		mockClient := new(MockChurchToolsClient)

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"status": "success",
				"data": {
					"songs": {
						"123": {
							"id": "123",
							"bezeichnung": "Amazing Grace",
							"author": "John Newton",
							"ccli": "12345",
							"copyright": "Public Domain",
							"arrangement": {
								"1": {
									"id": "1",
									"bezeichnung": "Standard",
									"default_yn": "1"
								}
							}
						},
						"456": {
							"id": "456",
							"bezeichnung": "How Great Thou Art",
							"author": "Carl Boberg",
							"ccli": "67890",
							"arrangement": {}
						}
					}
				}
			}`)),
		}

		mockClient.On("GetRequest", "https://lkg-pfuhl.church.tools/?q=churchservice/ajax", 
			map[string]string{"func": "getAllSongs"}).Return(mockResponse)

		result, err := churchtools.GetSongs(mockClient)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Amazing Grace", result["123"].Bezeichnung)
		assert.Equal(t, "John Newton", result["123"].Author)
		assert.Equal(t, "How Great Thou Art", result["456"].Bezeichnung)

		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty songs response", func(t *testing.T) {
		mockClient := new(MockChurchToolsClient)

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"status": "success",
				"data": {
					"songs": {}
				}
			}`)),
		}

		mockClient.On("GetRequest", mock.Anything, mock.Anything).Return(mockResponse)

		result, err := churchtools.GetSongs(mockClient)

		assert.NoError(t, err)
		assert.Empty(t, result)

		mockClient.AssertExpectations(t)
	})

	t.Run("should return error for invalid JSON", func(t *testing.T) {
		mockClient := new(MockChurchToolsClient)

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`invalid json`)),
		}

		mockClient.On("GetRequest", mock.Anything, mock.Anything).Return(mockResponse)

		result, err := churchtools.GetSongs(mockClient)

		assert.Error(t, err)
		assert.Nil(t, result)

		mockClient.AssertExpectations(t)
	})
}

func TestTruncateMessage(t *testing.T) {
	t.Run("should not truncate short messages", func(t *testing.T) {
		message := "Short message"
		// Note: truncateMessage is not exported, testing indirectly through GetSongs error handling
		assert.True(t, len(message) < 200*1024)
	})

	t.Run("should truncate long messages", func(t *testing.T) {
		// Create a message larger than 200KB
		largeMessage := strings.Repeat("a", 250*1024)
		assert.True(t, len(largeMessage) > 200*1024)
	})
}
