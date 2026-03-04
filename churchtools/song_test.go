package churchtools_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/stretchr/testify/assert"
)

func TestSong_GetDefaultArrangement(t *testing.T) {
	t.Run("should return the default arrangement", func(t *testing.T) {
		song := churchtools.Song{
			ID:          123,
			Bezeichnung: "Test Song",
			Arrangements: map[string]churchtools.SongArrangement{
				"1": {
					ID:          1,
					Bezeichnung: "Non-default",
					Default:     0,
				},
				"2": {
					ID:          2,
					Bezeichnung: "Default Arrangement",
					Default:     1,
				},
			},
		}

		result := song.GetDefaultArrangement()
		assert.Equal(t, 2, result.ID)
		assert.Equal(t, "Default Arrangement", result.Bezeichnung)
	})

	t.Run("should return empty arrangement if no default", func(t *testing.T) {
		song := churchtools.Song{
			ID:          123,
			Bezeichnung: "Test Song",
			Arrangements: map[string]churchtools.SongArrangement{
				"1": {
					ID:          1,
					Bezeichnung: "Non-default",
					Default:     0,
				},
			},
		}

		result := song.GetDefaultArrangement()
		assert.Equal(t, 0, result.ID)
	})
}

func TestSong_AddArrangement(t *testing.T) {
	t.Run("should add arrangement and return ID", func(t *testing.T) {
		mockClient := new(MockChurchToolsClient)
		song := churchtools.Song{
			ID:           123,
			Bezeichnung:  "Test Song",
			Arrangements: map[string]churchtools.SongArrangement{},
		}

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"data": {
					"id": 456,
					"name": "New Arrangement"
				}
			}`)),
		}

		expectedURL := "https://lkg-pfuhl.church.tools/api/songs/123/arrangements"
		expectedParams := map[string]string{"name": "New Arrangement"}
		mockClient.On("PostRequest", expectedURL, expectedParams).Return(mockResponse)

		arrangementID, err := song.AddArrangement(mockClient, "New Arrangement")

		assert.NoError(t, err)
		assert.Equal(t, 456, arrangementID)
		_, exists := song.Arrangements["456"]
		assert.True(t, exists)

		mockClient.AssertExpectations(t)
	})
}

func TestSong_GetModificationDate(t *testing.T) {
	t.Run("should parse modification date correctly", func(t *testing.T) {
		song := churchtools.Song{
			ModifiedDate: "2024-03-15 14:30:00",
		}

		result, err := song.GetModificationDate()

		assert.NoError(t, err)
		assert.Equal(t, 2024, result.Year())
		assert.Equal(t, time.March, result.Month())
		assert.Equal(t, 15, result.Day())
		assert.Equal(t, 14, result.Hour())
		assert.Equal(t, 30, result.Minute())
	})

	t.Run("should return error for invalid date format", func(t *testing.T) {
		song := churchtools.Song{
			ModifiedDate: "invalid-date",
		}

		_, err := song.GetModificationDate()

		assert.Error(t, err)
	})
}
