package path

import (
	"os"
	"path/filepath"
)

func Relative(path string) string {
	if path[0:2] != "./" {
		return path
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(ex) + "\\" + path[2:]
}
