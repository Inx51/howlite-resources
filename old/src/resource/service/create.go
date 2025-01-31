package service

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/inx51/howlite/resources/resource"
	"github.com/inx51/howlite/resources/storage"
	"github.com/vmihailenco/msgpack/v5"
)

func Create(res *resource.Resource, storage *storage.Storage) error {
	if Exists(res.Identifier, storage) {
		return resource.AlreadyExistsError{Identifier: res.Identifier}
	}

	//Output
	path := getPath(res.Identifier)
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//Headers
	res.Headers = filterHeaders(&res.Headers)
	headersLength := len(res.Headers)
	if headersLength == 0 {
		file.Write(make([]byte, 8))
	} else {
		msgPackedHeaders, err := msgpack.Marshal(res.Headers)
		if err != nil {
			panic(err)
		}
		headersMsgPackLength := len(msgPackedHeaders)
		binary.Write(file, binary.LittleEndian, uint64(headersMsgPackLength))
		file.Write(msgPackedHeaders)
	}

	//Body
	buff := make([]byte, 1024)
	readCloser := io.NopCloser(*res.Body)
	_, err = io.CopyBuffer(file, readCloser, buff)
	if err != nil {
		panic(err)
	}

	return nil
}
