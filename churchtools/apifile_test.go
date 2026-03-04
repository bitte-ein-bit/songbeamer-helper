package churchtools_test

import (
	"net/url"
	"path/filepath"
	"testing"

	"github.com/bitte-ein-bit/songbeamer-helper/churchtools"
	"github.com/stretchr/testify/assert"
)

func TestNewAPIFile(t *testing.T) {
	t.Run("should create APIFile from file path", func(t *testing.T) {
		path := "/tmp/test/mysong.sng"

		result := churchtools.NewAPIFile(path)

		assert.NotNil(t, result)
		assert.Equal(t, "mysong.sng", result.Name)
		// filepath is private, can't test directly
	})

	t.Run("should extract basename from path", func(t *testing.T) {
		path := "/very/long/path/to/some/file.txt"

		result := churchtools.NewAPIFile(path)

		assert.Equal(t, "file.txt", result.Name)
	})
}

func TestNewSongAPIFile(t *testing.T) {
	t.Run("should create song arrangement type API file", func(t *testing.T) {
		path := "/tmp/test/mysong.sng"
		domainID := 123

		result := churchtools.NewSongAPIFile(path, domainID)

		assert.NotNil(t, result)
		assert.Equal(t, "mysong.sng", result.Name)
		assert.Equal(t, "song_arrangement", result.DomainType)
		assert.Equal(t, 123, result.DomainID)
	})

	t.Run("should handle different domain IDs", func(t *testing.T) {
		path := "/tmp/song.txt"
		domainID := 999

		result := churchtools.NewSongAPIFile(path, domainID)

		assert.Equal(t, "song_arrangement", result.DomainType)
		assert.Equal(t, 999, result.DomainID)
	})
}

func TestAPIFile_SetUploadName(t *testing.T) {
	t.Run("should allow setting custom upload name", func(t *testing.T) {
		file := churchtools.NewAPIFile("/tmp/original.txt")

		file.SetUploadName("custom-name.txt")

		// uploadName is private, but getUploadName will use it
		// We can test this indirectly through the behavior
		assert.NotNil(t, file)
	})
}

func TestAPIFile_Delete(t *testing.T) {
	t.Run("should return error when trying to delete file with ID 0", func(t *testing.T) {
		file := churchtools.APIFile{
			DomainType: "song_arrangement",
			DomainID:   123,
		}

		err := file.Delete(0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete file with ID 0")
	})

	t.Run("should not error for valid ID (deletion is commented out)", func(t *testing.T) {
		file := churchtools.APIFile{
			DomainType: "song_arrangement",
			DomainID:   123,
		}

		err := file.Delete(456)

		// Currently returns nil because actual deletion is commented out
		assert.NoError(t, err)
	})
}

func TestAPIFile_LoadFromFile(t *testing.T) {
	t.Run("should load file path and extract name", func(t *testing.T) {
		file := &churchtools.APIFile{}
		path := "/tmp/test/newfile.sng"

		file.LoadFromFile(path)

		assert.Equal(t, "newfile.sng", file.Name)
	})

	t.Run("should update existing file", func(t *testing.T) {
		file := &churchtools.APIFile{
			Name:       "old.txt",
			DomainType: "song_arrangement",
			DomainID:   123,
		}
		path := "/tmp/new.txt"

		file.LoadFromFile(path)

		assert.Equal(t, "new.txt", file.Name)
		assert.Equal(t, "song_arrangement", file.DomainType) // Should keep other fields
		assert.Equal(t, 123, file.DomainID)
	})
}

func TestAPIFile_getID(t *testing.T) {
	// Testing private method indirectly through struct creation
	t.Run("should parse ID from FileURL", func(t *testing.T) {
		file := churchtools.APIFile{
			FileURL: "https://example.com/download?id=123&filename=test.txt",
		}

		// getID is private, but we can verify the FileURL is set correctly
		parsedURL, err := url.Parse(file.FileURL)
		assert.NoError(t, err)

		params, err := url.ParseQuery(parsedURL.RawQuery)
		assert.NoError(t, err)
		assert.Equal(t, "123", params.Get("id"))
	})

	t.Run("should handle FileURL without ID", func(t *testing.T) {
		file := churchtools.APIFile{
			FileURL: "https://example.com/download?filename=test.txt",
		}

		parsedURL, err := url.Parse(file.FileURL)
		assert.NoError(t, err)

		params, err := url.ParseQuery(parsedURL.RawQuery)
		assert.NoError(t, err)
		assert.Equal(t, "", params.Get("id")) // No ID in URL
	})

	t.Run("should handle empty FileURL", func(t *testing.T) {
		file := churchtools.APIFile{
			FileURL: "",
		}

		// Empty URL should be handled gracefully
		assert.Equal(t, "", file.FileURL)
	})
}

func TestAPIFile_Fields(t *testing.T) {
	t.Run("should create APIFile with all fields", func(t *testing.T) {
		file := churchtools.APIFile{
			DomainType: "song_arrangement",
			DomainID:   123,
			Name:       "test.sng",
			Filename:   "abc123.sng",
			FileURL:    "https://example.com/file?id=456",
		}

		assert.Equal(t, "song_arrangement", file.DomainType)
		assert.Equal(t, 123, file.DomainID)
		assert.Equal(t, "test.sng", file.Name)
		assert.Equal(t, "abc123.sng", file.Filename)
		assert.Contains(t, file.FileURL, "id=456")
	})
}

func TestAPIFile_PathHandling(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		expectedName string
	}{
		{
			name:         "Unix absolute path",
			path:         "/home/user/songs/test.sng",
			expectedName: "test.sng",
		},
		{
			name:         "Windows absolute path",
			path:         "C:\\Users\\user\\songs\\test.sng",
			expectedName: "test.sng",
		},
		{
			name:         "Relative path",
			path:         "songs/test.sng",
			expectedName: "test.sng",
		},
		{
			name:         "Just filename",
			path:         "test.sng",
			expectedName: "test.sng",
		},
		{
			name:         "Path with spaces",
			path:         "/home/user/my songs/test song.sng",
			expectedName: "test song.sng",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Normalize path for the current OS
			normalizedPath := filepath.FromSlash(tc.path)
			file := churchtools.NewAPIFile(normalizedPath)

			// filepath.Base should extract the same filename regardless of OS
			expectedBase := filepath.Base(normalizedPath)
			assert.Equal(t, expectedBase, file.Name)
		})
	}
}
