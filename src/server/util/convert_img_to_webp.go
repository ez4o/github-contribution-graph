package util

import (
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ConvertImgToWebp(buf []byte) ([]byte, error) {
	uuid := uuid.New().String()

	err := ioutil.WriteFile("./"+uuid, buf, 0644)
	if err != nil {
		return nil, err
	}

	err = ffmpeg.
		Input("./temp/"+uuid).
		Output("./temp/"+uuid+".webp", ffmpeg.KwArgs{
			"c:v": "libwebp",
			"crf": "51",
		}).
		Run()

	if err != nil {
		return nil, err
	}

	buf, err = ioutil.ReadFile("./" + uuid + ".webp")
	if err != nil {
		return nil, err
	}

	os.Remove("./temp/" + uuid)
	os.Remove("./temp/" + uuid + ".webp")

	return buf, nil
}
