package client

import (
	"reflect"

	"github.com/go-openapi/runtime"
	"golang.org/x/exp/slog"
)

func logErr(err error) {
	if e, ok := err.(*runtime.APIError); ok {
		traceID := ""
		resp, ok := e.Response.(runtime.ClientResponse)
		if ok {
			traceID = resp.GetHeader("X-Trace-ID")
		}
		slog.Error("Request failed", "code", e.Code, "trace.id", traceID)
	} else if e, ok := err.(errr); ok{
		slog.Error("Request failed", "code", e.Code(), "message", e.Error())
	} else {
		slog.Error("Request failed", "error", reflect.TypeOf(err), "message", err.Error())
	}
}

type errr interface {
	Error() string
	Code() int
}