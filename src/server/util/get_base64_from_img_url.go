package util

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/patrickmn/go-cache"
)

func GetBase64FromImgUrl(c *cache.Cache, imgUrl string) (string, error) {
	if data, found := c.Get("cachedBase64:" + imgUrl); found {
		log.Println("Found base64 in cache.")

		return data.(string), nil
	}

	resp, err := http.Get(imgUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Println("Image file size:", len(data))
	if len(data) > 2000000 {
		return "", fmt.Errorf("image is bigger than 2MB")
	}

	imgBase64Str := base64.StdEncoding.EncodeToString(data)

	c.Set("cachedBase64:"+imgUrl, imgBase64Str, cache.DefaultExpiration)

	return imgBase64Str, nil
}
