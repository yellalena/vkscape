package parser

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/yellalena/vkscape/internal/vkapi"
)

const PostTypePost = "post"

type VKParser struct {
	Client vkapi.VKClient
}

func InitParser(client vkapi.VKClient) VKParser {
	return VKParser{
		Client: client,
	}
}

func convertDate(timestamp int) string {
	// Convert Unix timestamp to a readable date format
	// Example: 1633072800 -> "20060102 (YYYYMMDD)"
	i, err := strconv.ParseInt(strconv.Itoa(timestamp), 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm.Format("20060102")
}

func downloadImage(url, filepath string) error {
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	return err
}
