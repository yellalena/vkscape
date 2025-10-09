package parser

import (
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
}

func InitParser() VKParser {
	return VKParser{}
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
	response, e := http.Get(url)
	if e != nil {
		return e
	}
	defer response.Body.Close()

	return utils.SaveObject(outputDir, filename, response.Body)
}
