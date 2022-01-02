package util

import (
	"net/url"
	"server/model"
	"strconv"
)

func GetRequestParams(v url.Values) model.RequestParams {
	username := v.Get("username")
	imgUrl := v.Get("img_url")
	lastNDays, err := strconv.Atoi(v.Get("last_n_days"))
	if lastNDays <= 0 || err != nil {
		lastNDays = 7
	}

	params := model.RequestParams{
		Username:  username,
		ImgUrl:    imgUrl,
		LastNDays: lastNDays,
	}

	return params
}
