package parser

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	vkObject "github.com/SevereCloud/vksdk/v2/object"
	"github.com/stretchr/testify/assert"
)

func TestProcessPost(t *testing.T) {
	wd, _ := os.Getwd()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("img"))
	}))
	t.Cleanup(server.Close)

	tests := []struct {
		name      string
		post      vkObject.WallWallpost
		setup     func(t *testing.T) string
		validate  bool
		wantError bool
		skipMsg   string
	}{
		{
			name: "non-post-type",
			post: vkObject.WallWallpost{
				ID:       1,
				Date:     1700000000,
				PostType: "copy",
				Text:     "test",
			},
			setup:     tempDirSetup(t, wd),
			wantError: false,
		},
		{
			name: "empty-post",
			post: vkObject.WallWallpost{
				ID:       2,
				Date:     1700000000,
				PostType: PostTypePost,
			},
			setup:     tempDirSetup(t, wd),
			wantError: false,
		},
		{
			name: "invalid-date",
			post: vkObject.WallWallpost{
				ID:       3,
				Date:     1700000000,
				PostType: PostTypePost,
				Text:     "test",
			},
			setup:   tempDirSetup(t, wd),
			skipMsg: "convertDate cannot fail with int input; requires refactor to test invalid date",
		},
		{
			name: "valid-date-saves-text",
			post: vkObject.WallWallpost{
				ID:       4,
				Date:     1700000000,
				PostType: PostTypePost,
				Text:     "hello",
			},
			setup:     tempDirSetup(t, wd),
			validate:  true,
			wantError: false,
		},
		{
			name: "save-file-fails",
			post: vkObject.WallWallpost{
				ID:       5,
				Date:     1700000000,
				PostType: PostTypePost,
				Text:     "test",
			},
			setup:   tempDirSetup(t, wd),
			skipMsg: "SaveFile failure requires injection or file permission control between CreateSubDirectory and SaveFile",
		},
		{
			name: "photo-no-sizes",
			post: vkObject.WallWallpost{
				ID:       6,
				Date:     1700000000,
				PostType: PostTypePost,
				Text:     "test",
				Attachments: []vkObject.WallWallpostAttachment{
					{Type: "photo", Photo: vkObject.PhotosPhoto{ID: 1}},
				},
			},
			setup:     tempDirSetup(t, wd),
			wantError: false,
		},
		{
			name: "photo-valid-sizes",
			post: vkObject.WallWallpost{
				ID:       7,
				Date:     1700000000,
				PostType: PostTypePost,
				Text:     "test",
				Attachments: []vkObject.WallWallpostAttachment{
					{
						Type: "photo",
						Photo: vkObject.PhotosPhoto{
							ID: 1,
							Sizes: []vkObject.PhotosPhotoSizes{{
								BaseImage: vkObject.BaseImage{URL: server.URL + "/img"},
							}},
						},
					},
				},
			},
			setup:     tempDirSetup(t, wd),
			validate:  true,
			wantError: false,
		},
		{
			name: "non-photo-attachment",
			post: vkObject.WallWallpost{
				ID:       8,
				Date:     1700000000,
				PostType: PostTypePost,
				Text:     "test",
				Attachments: []vkObject.WallWallpostAttachment{
					{Type: "video"},
				},
			},
			setup:     tempDirSetup(t, wd),
			wantError: false,
		},
	}

	p := VKParser{logger: slog.Default()}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skipMsg != "" {
				t.Skip(tc.skipMsg)
			}
			outDir := tc.setup(t)
			err := p.processPost(context.Background(), outDir, tc.post)
			if tc.wantError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tc.validate {
				validatePost(t, outDir, tc.post)
			}
		})
	}
}

func tempDirSetup(t *testing.T, wd string) func(t *testing.T) string {
	return func(t *testing.T) string {
		tmp := t.TempDir()
		t.Cleanup(func() { _ = os.Chdir(wd) })
		_ = os.Chdir(tmp)
		return tmp
	}
}

func validatePost(t *testing.T, outDir string, post vkObject.WallWallpost) {
	t.Helper()

	dateStr, err := convertDate(post.Date)
	assert.NoError(t, err)
	postName := fmt.Sprintf(PostFileNameTemplate, post.ID, dateStr)
	path := filepath.Join(outDir, postName)
	_, statErr := os.Stat(path)
	assert.NoError(t, statErr)
}
