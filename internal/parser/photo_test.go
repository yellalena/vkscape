package parser

import (
	"context"
	"log/slog"
	"testing"

	vkObject "github.com/SevereCloud/vksdk/v2/object"
	"github.com/stretchr/testify/assert"
)

func TestProcessPhotoEmptySizes(t *testing.T) {
	p := VKParser{logger: slog.Default()}
	err := p.processPhoto(context.Background(), "out", "1", vkObject.PhotosPhoto{})
	assert.Error(t, err)
}
