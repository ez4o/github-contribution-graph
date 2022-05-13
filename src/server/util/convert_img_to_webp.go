package util

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

func ConvertImgToWebp(buf []byte) ([]byte, error) {
	uuid := uuid.New().String()

	i := "./temp/" + uuid
	o := "./temp/" + uuid + ".webp"

	err := ioutil.WriteFile(i, buf, 0644)
	if err != nil {
		return nil, err
	}

	err = exec.Command("ffmpeg", "-i", i, "-c:v", "libwebp", "-crf", "51", o).Run()

	if err != nil {
		return nil, err
	}

	buf, err = ioutil.ReadFile(o)
	if err != nil {
		return nil, err
	}

	os.Remove(i)
	os.Remove(o)

	return buf, nil
}
