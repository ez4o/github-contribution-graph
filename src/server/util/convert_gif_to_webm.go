package util

import (
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ConvertGifToWebm(buf []byte) ([]byte, error) {
	uuid := uuid.New().String()

	err := ioutil.WriteFile("./"+uuid+".gif", buf, 0644)
	if err != nil {
		return nil, err
	}

	err = ffmpeg.
		Input("./temp/"+uuid+".gif").
		Output("./temp/"+uuid+".webm", ffmpeg.KwArgs{
			"y":        "",
			"r":        "16",
			"c:v":      "libvpx",
			"quality":  "good",
			"cpu-used": "0",
			"b:v":      "500K",
			"crf":      "51",
			"pix_fmt":  "yuv420p",
			"movflags": "faststart",
		}).
		Run()

	if err != nil {
		return nil, err
	}

	buf, err = ioutil.ReadFile("./" + uuid + ".webm")
	if err != nil {
		return nil, err
	}

	os.Remove("./temp/" + uuid + ".gif")
	os.Remove("./temp/" + uuid + ".webm")

	return buf, nil
}
