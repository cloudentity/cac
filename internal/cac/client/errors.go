package client

import (
	"reflect"

	"github.com/cloudentity/acp-client-go/clients/hub/client/workspace_configuration"
	"github.com/go-openapi/runtime"
	"golang.org/x/exp/slog"
)

func logErr(err error) {
	switch e := err.(type) {
	case *runtime.APIError:
		traceID := ""
		resp, ok := e.Response.(runtime.ClientResponse)
		if ok {
			traceID = resp.GetHeader("X-Trace-ID")
		}
		slog.Error("Request failed", "code", e.Code, "trace.id", traceID)
	case *workspace_configuration.PatchWorkspaceConfigRfc7396UnprocessableEntity:
	case *workspace_configuration.PatchWorkspaceConfigRfc6902BadRequest:
	case *workspace_configuration.ImportWorkspaceConfigBadRequest:
	case *workspace_configuration.ImportWorkspaceConfigUnprocessableEntity:
		slog.Error("Request failed", "code", e.Code, "message", e.Payload.Error)
	default:
		slog.Error("Request failed", "error", reflect.TypeOf(err), "message", err.Error())
	}
}