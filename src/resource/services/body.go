package services

import "io"

func WriteBody(resourceStream *io.WriteCloser, body *io.ReadCloser) {
	buff := make([]byte, 1024)
	readCloser := io.NopCloser(*body)
	_, err := io.CopyBuffer(*resourceStream, readCloser, buff)
	if err != nil {
		panic(err)
	}

	(*resourceStream).Close()
}
