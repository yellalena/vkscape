package vkscape

import (
	"fmt"
	"os"
)

const (
	OutputDir      = "vkscape_output"
	OutputGroupDir = "group_%s"
)

func CreateGroupDirectory(groupID string) string {
	groupDir := fmt.Sprintf("%s/%s", OutputDir, fmt.Sprintf(OutputGroupDir, groupID))
	_ = os.MkdirAll(groupDir, 0755)
	return groupDir
}
