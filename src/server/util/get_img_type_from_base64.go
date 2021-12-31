package util

import "fmt"

func GetImgTypeFromBase64(base64String string) (string, error) {
	switch base64String[0] {
	case '/':
		return "image/jpeg", nil
	case 'i':
		return "image/png", nil
	case 'R':
		return "image/gif", nil
	case 'U':
		return "image/webp", nil
	default:
		return "", fmt.Errorf("Unsupported image type.")
	}
}
