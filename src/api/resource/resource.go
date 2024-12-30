package resource

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"

	reserrors "github.com/inx51/howlite/resources/api/resource/errors"
	"github.com/inx51/howlite/resources/api/utils"
)

type Resource struct {
	Body       *os.File
	Identifier *string
}

func New(identifier *string, body *io.ReadCloser, headers *http.Header) Resource {

}

func GetIdentifier(path *string) string {
	return utils.Base64HashString(*path)
}

func Get(identifier *string) (Resource, error) {

	file, err := os.Open("")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Resource{}, &reserrors.NotFoundError{Identifier: identifier}
		}

		return Resource{}, err
	}

	return Resource{
		Identifier: identifier,
		Body:       file,
	}, nil
}
