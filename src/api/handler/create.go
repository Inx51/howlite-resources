package handler

import (
	"net/http"

	"github.com/inx51/howlite/resources/api/resource"
)

func CreateResource(resp *http.ResponseWriter, req *http.Request) {
	identifier := resource.GetIdentifier(&req.URL.Path)
	res := resource.New(&identifier, &req.Body, &req.Header)
	resource.Presist(res)
	// fileLocation := fmt.Sprintf("C:/test/%s.bin", identifier)
	// file, err := os.Create(fileLocation)
	// if err != nil {
	// 	panic(err)
	// }
	// buff := make([]byte, 1024)
	// io.CopyBuffer(file, req.Body, buff)

	// defer func() {
	// 	if err := file.Close(); err != nil {
	// 		panic(err)
	// 	}
	// 	if err := req.Body.Close(); err != nil {
	// 		panic(err)
	// 	}
	// }()
}
