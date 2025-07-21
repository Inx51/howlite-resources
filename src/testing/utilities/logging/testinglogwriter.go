package logging

type TestingLogWriter struct {
}

func (writer *TestingLogWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
