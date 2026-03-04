package churchtools_test

import (
	"testing"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/stretchr/testify/assert"
)

func TestSongFile_ToAPIFile(t *testing.T) {
	t.Run("should convert SongFile to APIFile correctly", func(t *testing.T) {
		songFile := churchtools.SongFile{
			ID:           123,
			DomainType:   "song_arrangement",
			DomainID:     456,
			Bezeichnung:  "Test Song.sng",
			Filename:     "somehash.sng",
			ModifiedDate: "2024-03-15 14:30:00",
		}

		result := songFile.ToAPIFile()

		assert.Equal(t, "song_arrangement", result.DomainType)
		assert.Equal(t, 456, result.DomainID)
		assert.Equal(t, "Test Song.sng", result.Name)
		assert.Equal(t, "somehash.sng", result.Filename)
		assert.Contains(t, result.FileURL, "id=123")
		assert.Contains(t, result.FileURL, "filename=somehash.sng")
		assert.Contains(t, result.FileURL, "public/filedownload")
	})
}

func TestSongFile_GetModificationDate(t *testing.T) {
	t.Run("should parse modification date with timezone correctly", func(t *testing.T) {
		songFile := churchtools.SongFile{
			ModifiedDate: "2024-03-15 14:30:00",
		}

		result, err := songFile.GetModificationDate()

		assert.NoError(t, err)
		assert.Equal(t, 2024, result.Year())
		assert.Equal(t, time.March, result.Month())
		assert.Equal(t, 15, result.Day())
		// Time should be adjusted for +0100 timezone
		assert.Equal(t, 14, result.Hour())
		assert.Equal(t, 30, result.Minute())
	})

	t.Run("should return error for invalid date format", func(t *testing.T) {
		songFile := churchtools.SongFile{
			ModifiedDate: "invalid-date",
		}

		_, err := songFile.GetModificationDate()

		assert.Error(t, err)
	})

	t.Run("should round to nearest second", func(t *testing.T) {
		songFile := churchtools.SongFile{
			ModifiedDate: "2024-03-15 14:30:00",
		}

		result, err := songFile.GetModificationDate()

		assert.NoError(t, err)
		// Verify it's rounded to the second (no nanoseconds)
		assert.Equal(t, 0, result.Nanosecond())
	})
}
