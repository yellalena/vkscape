package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yellalena/vkscape/internal/models"
)

func TestCreateAlbumDirectory(t *testing.T) {
	cases := []struct {
		name      string
		setup     func(t *testing.T)
		wantError bool
	}{
		{
			name: "success",
			setup: func(t *testing.T) {
				tmp := t.TempDir()
				wd, _ := os.Getwd()
				t.Cleanup(func() { _ = os.Chdir(wd) })
				_ = os.Chdir(tmp)
			},
			wantError: false,
		},
		{
			name: "invalid path",
			setup: func(t *testing.T) {
				tmp := t.TempDir()
				wd, _ := os.Getwd()
				t.Cleanup(func() { _ = os.Chdir(wd) })
				_ = os.Chdir(tmp)
				if err := os.WriteFile(OutputDir, []byte("not a dir"), 0600); err != nil {
					assert.NoError(t, err)
				}
			},
			wantError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(t)
			album := models.PhotoAlbum{ID: 123, Title: "Title", Description: "Desc"}
			dir, err := CreateAlbumDirectory(album)
			if tc.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			infoPath := filepath.Join(dir, "album_info.txt")
			_, statErr := os.Stat(infoPath)
			assert.NoError(t, statErr)
		})
	}
}
