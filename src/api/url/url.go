package url

import (
	"fmt"
	"net/url"
)

func GetAbsolute(url *url.URL) string {
	return fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, url.Path)
}
