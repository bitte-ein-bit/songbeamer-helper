package churchtools_test

import (
	"fmt"
	"testing"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/stretchr/testify/assert"
)

func TestAPISong_GetDefaultArrangement(t *testing.T) {
	t.Run("should return the default arrangement", func(t *testing.T) {
		song := churchtools.APISong{
			ID:          123,
			Bezeichnung: "Test Song",
			Arrangements: []churchtools.APISongArrangement{
				{
					ID:      1,
					Name:    "Non-default",
					Default: false,
				},
				{
					ID:      2,
					Name:    "Default Arrangement",
					Default: true,
				},
				{
					ID:      3,
					Name:    "Another Non-default",
					Default: false,
				},
			},
		}

		result := song.GetDefaultArrangement()

		assert.Equal(t, 2, result.ID)
		assert.Equal(t, "Default Arrangement", result.Name)
		assert.True(t, result.Default)
	})

	t.Run("should return empty arrangement if no default exists", func(t *testing.T) {
		song := churchtools.APISong{
			ID:          123,
			Bezeichnung: "Test Song",
			Arrangements: []churchtools.APISongArrangement{
				{
					ID:      1,
					Name:    "Non-default",
					Default: false,
				},
			},
		}

		result := song.GetDefaultArrangement()

		assert.Equal(t, 0, result.ID)
		assert.Equal(t, "", result.Name)
	})

	t.Run("should return empty arrangement if no arrangements", func(t *testing.T) {
		song := churchtools.APISong{
			ID:           123,
			Bezeichnung:  "Test Song",
			Arrangements: []churchtools.APISongArrangement{},
		}

		result := song.GetDefaultArrangement()

		assert.Equal(t, 0, result.ID)
	})
}

func TestAPISong_ToSong(t *testing.T) {
	t.Run("should convert APISong to Song correctly", func(t *testing.T) {
		apiSong := churchtools.APISong{
			ID:          123,
			Bezeichnung: "Amazing Grace",
			Category: churchtools.APISongCategory{
				ID:   5,
				Name: "Worship",
			},
			Practice:  true,
			Author:    "John Newton",
			CCLI:      "12345",
			Copyright: "Public Domain",
			Note:      "Test notes",
			Arrangements: []churchtools.APISongArrangement{
				{
					ID:               1,
					Name:             "Standard",
					Default:          true,
					KeyOfArrangement: "G",
					BPM:              "120",
					Beat:             "4/4",
					Duration:         185, // 3 minutes 5 seconds
					Note:             "Arrangement notes",
				},
				{
					ID:      2,
					Name:    "Acoustic",
					Default: false,
				},
			},
		}

		result := apiSong.ToSong()

		assert.Equal(t, 123, result.ID)
		assert.Equal(t, "Amazing Grace", result.Bezeichnung)
		assert.Equal(t, 5, result.SongcategoryID)
		assert.Equal(t, "John Newton", result.Author)
		assert.Equal(t, "12345", result.CCLI)
		assert.Equal(t, "Public Domain", result.Copyright)
		assert.Equal(t, "Test notes", result.Note)
		assert.Len(t, result.Arrangements, 2)

		arrangement1, exists := result.Arrangements["1"]
		assert.True(t, exists)
		assert.Equal(t, 1, arrangement1.ID)
		assert.Equal(t, "Standard", arrangement1.Bezeichnung)
		assert.Equal(t, 1, arrangement1.Default)
		assert.Equal(t, "G", arrangement1.Tonality)
		assert.Equal(t, "120", arrangement1.BPM)
		assert.Equal(t, 3, arrangement1.Minutes)
		assert.Equal(t, 5, arrangement1.Seconds)

		arrangement2, exists := result.Arrangements["2"]
		assert.True(t, exists)
		assert.Equal(t, 2, arrangement2.ID)
		assert.Equal(t, 0, arrangement2.Default)
	})

	t.Run("should handle empty arrangements", func(t *testing.T) {
		apiSong := churchtools.APISong{
			ID:           123,
			Bezeichnung:  "Test Song",
			Arrangements: []churchtools.APISongArrangement{},
		}

		result := apiSong.ToSong()

		assert.Equal(t, 123, result.ID)
		assert.Empty(t, result.Arrangements)
	})
}

func TestAPISongCategory(t *testing.T) {
	t.Run("should create APISongCategory with all fields", func(t *testing.T) {
		category := churchtools.APISongCategory{
			ID:             1,
			Name:           "Worship",
			NameTranslated: "Anbetung",
			SortKey:        10,
			CampusID:       5,
		}

		assert.Equal(t, 1, category.ID)
		assert.Equal(t, "Worship", category.Name)
		assert.Equal(t, "Anbetung", category.NameTranslated)
		assert.Equal(t, 10, category.SortKey)
		assert.Equal(t, 5, category.CampusID)
	})
}

func TestAPISongArrangement_ToArrangement(t *testing.T) {
	t.Run("should convert APISongArrangement to SongArrangement with default true", func(t *testing.T) {
		apiArrangement := churchtools.APISongArrangement{
			ID:               123,
			Name:             "Standard",
			Default:          true,
			KeyOfArrangement: "G",
			BPM:              "120",
			Beat:             "4/4",
			Duration:         245, // 4 minutes 5 seconds
			Note:             "Test notes",
		}

		result := apiArrangement.ToArrangement()

		assert.Equal(t, 123, result.ID)
		assert.Equal(t, "Standard", result.Bezeichnung)
		assert.Equal(t, 1, result.Default)
		assert.Equal(t, "G", result.Tonality)
		assert.Equal(t, "120", result.BPM)
		assert.Equal(t, "4/4", result.Beat)
		assert.Equal(t, 4, result.Minutes)
		assert.Equal(t, 5, result.Seconds)
		assert.Equal(t, "Test notes", result.Note)
	})

	t.Run("should convert APISongArrangement to SongArrangement with default false", func(t *testing.T) {
		apiArrangement := churchtools.APISongArrangement{
			ID:       456,
			Name:     "Acoustic",
			Default:  false,
			Duration: 125, // 2 minutes 5 seconds
		}

		result := apiArrangement.ToArrangement()

		assert.Equal(t, 456, result.ID)
		assert.Equal(t, "Acoustic", result.Bezeichnung)
		assert.Equal(t, 0, result.Default)
		assert.Equal(t, 2, result.Minutes)
		assert.Equal(t, 5, result.Seconds)
	})

	t.Run("should handle duration edge cases", func(t *testing.T) {
		testCases := []struct {
			duration        int
			expectedMinutes int
			expectedSeconds int
		}{
			{0, 0, 0},
			{59, 0, 59},
			{60, 1, 0},
			{61, 1, 1},
			{3661, 61, 1},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("duration_%d", tc.duration), func(t *testing.T) {
				apiArrangement := churchtools.APISongArrangement{
					Duration: tc.duration,
				}

				result := apiArrangement.ToArrangement()

				assert.Equal(t, tc.expectedMinutes, result.Minutes)
				assert.Equal(t, tc.expectedSeconds, result.Seconds)
			})
		}
	})
}
