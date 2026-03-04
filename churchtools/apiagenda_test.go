package churchtools_test

import (
	"strings"
	"testing"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/stretchr/testify/assert"
)

func TestAPIAgendaSong_ToFilename(t *testing.T) {
	t.Run("should convert song to filename correctly", func(t *testing.T) {
		song := churchtools.APIAgendaSong{
			SongID:        123,
			ArrangementID: 456,
			Title:         "Amazing Grace",
			Arrangement:   "Standard",
		}

		result := song.ToFilename()

		assert.Equal(t, "Amazing Grace - Standard.sng", result)
	})

	t.Run("should replace slashes in title with underscores", func(t *testing.T) {
		song := churchtools.APIAgendaSong{
			Title:       "Jesus/Savior",
			Arrangement: "Acoustic",
		}

		result := song.ToFilename()

		assert.Equal(t, "Jesus_Savior - Acoustic.sng", result)
		assert.NotContains(t, result, "/")
	})

	t.Run("should handle multiple slashes", func(t *testing.T) {
		song := churchtools.APIAgendaSong{
			Title:       "Alpha/Omega/Beginning/End",
			Arrangement: "Standard",
		}

		result := song.ToFilename()

		assert.Equal(t, "Alpha_Omega_Beginning_End - Standard.sng", result)
		assert.NotContains(t, result, "/")
	})
}

func TestAPIAgendaItem_ToSongbeamerItem(t *testing.T) {
	t.Run("should convert song type item correctly", func(t *testing.T) {
		item := churchtools.APIAgendaItem{
			ID:       1,
			Position: 1,
			Type:     "song",
			Title:    "Opening Song",
			Song: churchtools.APIAgendaSong{
				SongID:      123,
				Title:       "Amazing Grace",
				Arrangement: "Standard",
			},
		}

		result := item.ToSongbeamerItem()

		assert.Contains(t, result, "item")
		assert.Contains(t, result, "Caption = 'Amazing Grace - Standard'")
		assert.Contains(t, result, "Color = clBlue")
		assert.Contains(t, result, "FileName = 'Amazing Grace - Standard.sng'")
		assert.Contains(t, result, "end")
	})

	t.Run("should convert header type item correctly", func(t *testing.T) {
		item := churchtools.APIAgendaItem{
			ID:       2,
			Position: 2,
			Type:     "header",
			Title:    "Worship",
		}

		result := item.ToSongbeamerItem()

		assert.Contains(t, result, "Caption = 'Worship'")
		assert.Contains(t, result, "Color = clBlack")
		assert.NotContains(t, result, "FileName")
	})

	t.Run("should convert normal type item correctly", func(t *testing.T) {
		item := churchtools.APIAgendaItem{
			ID:       3,
			Position: 3,
			Type:     "normal",
			Title:    "Sermon",
		}

		result := item.ToSongbeamerItem()

		assert.Contains(t, result, "Caption = 'Sermon'")
		assert.Contains(t, result, "Color = 33023")
	})

	t.Run("should handle Count-Down / Jingle as intro", func(t *testing.T) {
		item := churchtools.APIAgendaItem{
			ID:       4,
			Position: 1,
			Type:     "normal",
			Title:    "Count-Down / Jingle",
		}

		result := item.ToSongbeamerItem()

		assert.Contains(t, result, "Color = clBlack")
		assert.Contains(t, result, "FileName = 'C:\\Program Files (x86)\\SongBeamer\\Intro.mp3'")
	})

	t.Run("should handle Kinderprogramm with MP3", func(t *testing.T) {
		item := churchtools.APIAgendaItem{
			ID:       5,
			Position: 1,
			Type:     "normal",
			Title:    "Kinderprogramm Ankündigung",
		}

		result := item.ToSongbeamerItem()

		assert.Contains(t, result, "Color = clBlack")
		assert.Contains(t, result, "FileName = 'C:\\Program Files (x86)\\SongBeamer\\Kinder.mp3'")
	})

	t.Run("should escape single quotes in caption", func(t *testing.T) {
		item := churchtools.APIAgendaItem{
			ID:       6,
			Position: 1,
			Type:     "normal",
			Title:    "God's Love",
		}

		result := item.ToSongbeamerItem()

		assert.Contains(t, result, "Caption = 'God'#39's Love'")
		assert.NotContains(t, result, "God's Love")
	})

	t.Run("should escape single quotes in filename", func(t *testing.T) {
		item := churchtools.APIAgendaItem{
			ID:       7,
			Position: 1,
			Type:     "song",
			Title:    "Song Title",
			Song: churchtools.APIAgendaSong{
				Title:       "God's Love",
				Arrangement: "Standard",
			},
		}

		result := item.ToSongbeamerItem()

		assert.Contains(t, result, "FileName = 'God'#39's Love - Standard.sng'")
	})

	t.Run("should handle unknown type with default blue color", func(t *testing.T) {
		item := churchtools.APIAgendaItem{
			ID:       8,
			Position: 1,
			Type:     "unknown_type",
			Title:    "Unknown Item",
		}

		result := item.ToSongbeamerItem()

		assert.Contains(t, result, "Color = clBlue")
	})
}

func TestAPIAgenda(t *testing.T) {
	t.Run("should create APIAgenda with all fields", func(t *testing.T) {
		agenda := churchtools.APIAgenda{
			ID:      123,
			Name:    "Sunday Service",
			IsFinal: true,
			Items: []churchtools.APIAgendaItem{
				{
					ID:       1,
					Position: 1,
					Type:     "song",
					Title:    "Opening Song",
				},
			},
		}

		assert.Equal(t, 123, agenda.ID)
		assert.Equal(t, "Sunday Service", agenda.Name)
		assert.True(t, agenda.IsFinal)
		assert.Len(t, agenda.Items, 1)
	})
}

func TestAPIAgendaSong_Fields(t *testing.T) {
	t.Run("should create APIAgendaSong with all fields", func(t *testing.T) {
		song := churchtools.APIAgendaSong{
			SongID:        123,
			ArrangementID: 456,
			Title:         "Amazing Grace",
			Arrangement:   "Standard",
			Category:      "Worship",
			Key:           "G",
			BPM:           120,
			IsDefault:     true,
		}

		assert.Equal(t, 123, song.SongID)
		assert.Equal(t, 456, song.ArrangementID)
		assert.Equal(t, "Amazing Grace", song.Title)
		assert.Equal(t, "Standard", song.Arrangement)
		assert.Equal(t, "Worship", song.Category)
		assert.Equal(t, "G", song.Key)
		assert.Equal(t, 120, song.BPM)
		assert.True(t, song.IsDefault)
	})
}

func TestAPIAgendaItem_Fields(t *testing.T) {
	t.Run("should create APIAgendaItem with all fields", func(t *testing.T) {
		item := churchtools.APIAgendaItem{
			ID:            1,
			Position:      5,
			Type:          "song",
			Title:         "Opening Song",
			IsBeforeEvent: true,
			ArrangementID: 123,
			Song: churchtools.APIAgendaSong{
				SongID: 456,
				Title:  "Test Song",
			},
		}

		assert.Equal(t, 1, item.ID)
		assert.Equal(t, 5, item.Position)
		assert.Equal(t, "song", item.Type)
		assert.Equal(t, "Opening Song", item.Title)
		assert.True(t, item.IsBeforeEvent)
		assert.Equal(t, 123, item.ArrangementID)
		assert.Equal(t, 456, item.Song.SongID)
	})
}

func TestAPIAgendaItem_SpecialTitles(t *testing.T) {
	testCases := []struct {
		name         string
		title        string
		originalType string
		expectedType string
	}{
		{
			name:         "Count-Down / Jingle should become intro",
			title:        "Count-Down / Jingle",
			originalType: "normal",
			expectedType: "intro",
		},
		{
			name:         "Kinderprogramm should become kinder",
			title:        "Kinderprogramm Ankündigung",
			originalType: "header",
			expectedType: "kinder",
		},
		{
			name:         "Regular title should keep type",
			title:        "Regular Item",
			originalType: "normal",
			expectedType: "normal",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			item := churchtools.APIAgendaItem{
				Title: tc.title,
				Type:  tc.originalType,
			}

			result := item.ToSongbeamerItem()

			// Verify the correct type behavior through color and filename checks
			if tc.expectedType == "intro" {
				assert.Contains(t, result, "Intro.mp3")
			} else if tc.expectedType == "kinder" {
				assert.Contains(t, result, "Kinder.mp3")
			}
		})
	}
}

func TestAPIAgendaItem_FileNameEscaping(t *testing.T) {
	t.Run("should handle filename with slashes and quotes", func(t *testing.T) {
		item := churchtools.APIAgendaItem{
			Type:  "song",
			Title: "Test",
			Song: churchtools.APIAgendaSong{
				Title:       "Jesus/God's Love",
				Arrangement: "Standard's Version",
			},
		}

		result := item.ToSongbeamerItem()

		// Slashes should be replaced with underscores
		assert.Contains(t, result, "Jesus_God")
		// Single quotes should be escaped
		assert.True(t, strings.Contains(result, "'#39'"))
	})
}
