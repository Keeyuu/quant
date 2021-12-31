package imgutil

import (
	"strings"
)

func FillImgUrl(hostPrefix string, imgPath string) string {
	if strings.HasPrefix(imgPath, "http") {
		return imgPath
	} else if !strings.HasPrefix(imgPath, "/") {
		return hostPrefix + "/" + imgPath
	} else {
		return hostPrefix + imgPath
	}
}
