package util

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetBase64FromImgUrl(imgUrl string) (string, error) {
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
		return "", fmt.Errorf("Image is bigger than 2MB.")
	}

	imgBase64Str := base64.StdEncoding.EncodeToString(data)
	return imgBase64Str, nil
}
