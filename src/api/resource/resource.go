package resource

import (
	"os"
)

type Resource struct {
	Body       *os.File
	Identifier string
}
