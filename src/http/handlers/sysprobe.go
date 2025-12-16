package handlers

import (
	"context"
	"net/http"

	"github.com/inx51/howlite-resources/logger"
)

type SysProbeHandler struct {
}

func (handler *SysProbeHandler) Method() string {
	return "HEAD"
}

func (handler *SysProbeHandler) Path() string {
	return "/$sys/probe"
}

func (handler *SysProbeHandler) Handle(
	ctx context.Context,
	req *http.Request,
	resp http.ResponseWriter) (int, error) {

	logger.Debug(ctx, "sysprobe called")

	resp.WriteHeader(http.StatusNoContent)
	return http.StatusNoContent, nil
}

func NewSysProbeHandler() Handler {
	return &SysProbeHandler{}
}
