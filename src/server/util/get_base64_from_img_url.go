package util

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server/model"

	"github.com/gabriel-vasile/mimetype"
	"github.com/patrickmn/go-cache"
)

func GetImageFromUrl(c *cache.Cache, imgUrl string) (model.Image, error) {
	if imgBase64Str, found := c.Get("cachedBase64:" + imgUrl); found {
		log.Println("Found base64 in cache.")

		return imgBase64Str.(model.Image), nil
	}

	resp, err := http.Get(imgUrl)
	if err != nil {
		return model.Image{}, err
	}
	defer resp.Body.Close()

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.Image{}, err
	}

	log.Println("Image file size:", len(buffer))
	if len(buffer) > 2000000 {
		return model.Image{}, fmt.Errorf("image is bigger than 2MB")
	}

	var image model.Image
	var imageBytes []byte

	image.MimeType = mimetype.Detect(buffer).String()

	if image.MimeType == "image/gif" || image.MimeType == "video/webm" {
		imageBytes, err = ConvertGifToWebm(buffer)

		if err != nil {
			return model.Image{}, err
		}

		image.MimeType = "video/webm"
	} else if image.MimeType == "image/jpeg" || image.MimeType == "image/png" || image.MimeType == "image/webp" {
		imageBytes, err = ConvertImgToWebp(buffer)

		if err != nil {
			return model.Image{}, err
		}

		image.MimeType = "image/webp"
	}

	image.Base64String = base64.StdEncoding.EncodeToString(imageBytes)

	c.Set("cachedBase64:"+imgUrl, image, cache.DefaultExpiration)

	return image, nil
}
