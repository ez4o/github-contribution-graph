package util

import "fmt"

func GetImgTypeFromBase64(base64Head byte) (string, error) {
	switch base64Head {
	case '/':
		return "image/jpeg", nil
	case 'i':
		return "image/png", nil
	case 'R':
		return "image/gif", nil
	case 'U':
		return "image/webp", nil
	default:
		return "", fmt.Errorf("unsupported image type")
	}
}
