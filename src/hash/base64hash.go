package hash

import (
	"crypto/md5"
	"encoding/base64"
)

func Base64HashString(value string) string {
	var encBytes = md5.Sum([]byte(value))
	return base64.StdEncoding.EncodeToString(encBytes[:])
}
