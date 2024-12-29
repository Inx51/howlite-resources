package resource

import "github.com/inx51/howlite/resources/api/utils"

func GetIdentifier(path string) string {
	return utils.Base64HashString(path)
}
