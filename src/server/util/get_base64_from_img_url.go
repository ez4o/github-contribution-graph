package util

import (
	"encoding/base64"
	"io/ioutil"
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

	imgBase64Str := base64.StdEncoding.EncodeToString(data)
	return imgBase64Str, nil
}
