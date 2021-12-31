package util

import (
	"net/url"
	"server/model"
)

func GetRequestParams(v url.Values) model.RequestParams {
	username := v.Get("username")
	imgUrl := v.Get("img_url")

	params := model.RequestParams{
		Username: username,
		ImgUrl:   imgUrl,
	}

	return params
}
