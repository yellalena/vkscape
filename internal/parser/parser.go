package parser

import (
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
	close(p.errs)
	return len(p.errs)
}

func convertDate(timestamp int) string {
	// Convert Unix timestamp to a readable date format
	// Example: 1633072800 -> "20060102 (YYYYMMDD)"
	i, err := strconv.ParseInt(strconv.Itoa(timestamp), 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm.Format(DateFormat)
}

func downloadImage(url, outputDir, filename string) error {
	response, e := http.Get(url) //nolint:gosec // URL from trusted VK API
	if e != nil {
		return e
	}
	defer response.Body.Close() //nolint:errcheck

	return utils.SaveObject(outputDir, filename, response.Body)
}
