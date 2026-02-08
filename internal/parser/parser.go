package parser

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/yellalena/vkscape/internal/utils"
)

const (
	PostTypePost = "post"
	DateFormat   = "20060102"
)

type VKParser struct {
	logger *slog.Logger
	errs   chan error
}

func InitParser(logger *slog.Logger) VKParser {
	return VKParser{
		logger: logger,
	}
}

func (p *VKParser) CloseErrorsAndCount() int {
	if p.errs == nil {
		return 0
	}
	close(p.errs)
	count := len(p.errs)
	p.errs = nil
	return count
}

func convertDate(timestamp int) (string, error) {
	// Convert Unix timestamp to a readable date format
	// Example: 1633072800 -> "20060102 (YYYYMMDD)"
	i, err := strconv.ParseInt(strconv.Itoa(timestamp), 10, 64)
	if err != nil {
		return "", err
	}
	tm := time.Unix(i, 0)
	return tm.Format(DateFormat), nil
}

func downloadImage(ctx context.Context, url, outputDir, filename string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) //nolint:gosec // URL from trusted VK API
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close() //nolint:errcheck

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d for %s", response.StatusCode, url)
	}

	return utils.SaveObject(outputDir, filename, response.Body)
}
