package service

import (
	"slices"
	"strings"

	"github.com/inx51/howlite/resources/config"
	"github.com/inx51/howlite/resources/hash"
	"github.com/inx51/howlite/resources/resource"
)

func filterHeaders(headers *map[string][]string) map[string][]string {
	forbiddenHeaders := []string{"host", "accept-encoding", "connection", "accepts", "user-agent"}
	var result = make(map[string][]string)
	for k, v := range *headers {
		if slices.Contains(forbiddenHeaders, strings.ToLower(k)) {
			continue
		}

		result[k] = v
	}

	return result
}

func getPath(identifier *resource.ResourceIdentifier) string {
	hashedFileName := hash.Base64HashString(*identifier.Value)
	return config.Instance.Storage.Path + "\\" + hashedFileName + ".bin"
}
