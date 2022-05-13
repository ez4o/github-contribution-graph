package util

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

func CompressGif(minBuf []byte) []byte {
	uuid := uuid.New().String()

	i := "./temp/" + uuid + ".gif"
	o1 := "./temp/" + uuid + "-opt1.gif"
	o2 := "./temp/" + uuid + "-opt2.gif"

	err := ioutil.WriteFile(i, minBuf, 0644)
	if err != nil {
		return minBuf
	}

	err = exec.Command("ffmpeg", "-y", "-i", i, "-vf", "fps=12,scale=320:-1", o1).Run()
	if err != nil {
		return minBuf
	}

	buf1, err := ioutil.ReadFile(o1)
	if err != nil {
		return minBuf
	}

	if len(buf1) < len(minBuf) {
		minBuf = buf1
	}

	err = exec.Command("gifsicle", "-i", o1, "-O3", "--colors", "32", "-o", o2).Run()
	if err != nil {
		return minBuf
	}

	buf2, err := ioutil.ReadFile(o2)
	if err != nil {
		return minBuf
	}

	if len(buf2) < len(minBuf) {
		minBuf = buf2
	}

	os.Remove(i)
	os.Remove(o1)
	os.Remove(o2)

	return minBuf
}
