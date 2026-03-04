package churchtools_test

import (
	"testing"
	"time"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/stretchr/testify/assert"
)

func TestSongArrangement_GetModificationDate(t *testing.T) {
	t.Run("should parse modification date correctly", func(t *testing.T) {
		arrangement := churchtools.SongArrangement{
			ModifiedDate: "2024-03-15 14:30:00",
		}

		result, err := arrangement.GetModificationDate()

		assert.NoError(t, err)
		assert.Equal(t, 2024, result.Year())
		assert.Equal(t, time.March, result.Month())
		assert.Equal(t, 15, result.Day())
		assert.Equal(t, 14, result.Hour())
		assert.Equal(t, 30, result.Minute())
	})

	t.Run("should return error for invalid date format", func(t *testing.T) {
		arrangement := churchtools.SongArrangement{
			ModifiedDate: "invalid-date",
		}

		_, err := arrangement.GetModificationDate()

		assert.Error(t, err)
	})
}

func TestSongArrangement_GetName(t *testing.T) {
	t.Run("should return Bezeichnung when set (old API)", func(t *testing.T) {
		arrangement := churchtools.SongArrangement{
			Bezeichnung: "Old API Name",
			Name:        "",
		}

		result := arrangement.GetName()

		assert.Equal(t, "Old API Name", result)
	})

	t.Run("should return Name when Bezeichnung is empty (new API)", func(t *testing.T) {
		arrangement := churchtools.SongArrangement{
			Bezeichnung: "",
			Name:        "New API Name",
		}

		result := arrangement.GetName()

		assert.Equal(t, "New API Name", result)
	})

	t.Run("should prefer Bezeichnung when both are set", func(t *testing.T) {
		arrangement := churchtools.SongArrangement{
			Bezeichnung: "Old API Name",
			Name:        "New API Name",
		}

		result := arrangement.GetName()

		assert.Equal(t, "Old API Name", result)
	})
}
