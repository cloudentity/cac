package client

import (
	"context"
	"crypto/tls"
	"github.com/cloudentity/acp-client-go"
	"github.com/cloudentity/acp-client-go/clients/hub/client/workspace_configuration"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/go-openapi/runtime"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slog"
	"net/http"
)

type Client struct {
	acp *acpclient.Client
}

func InitClient(config Configuration) (c *Client, err error) {
	var (
		acp acpclient.Client
	)

	if config.Insecure {
		config.HttpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	if acp, err = acpclient.New(config.Config); err != nil {
		return nil, err
	}

	slog.With("client", acp).Debug("Initiated client")

	return &Client{
		acp: &acp,
	}, nil
}

func (c *Client) PullWorkspaceConfiguration(ctx context.Context, workspace string, secrets bool) (*models.TreeServer, error) {
	var (
		ok  *workspace_configuration.ExportWorkspaceConfigOK
		err error
	)

	if ok, err = c.acp.Hub.WorkspaceConfiguration.
		ExportWorkspaceConfig(workspace_configuration.
			NewExportWorkspaceConfigParams().
			WithContext(ctx).
			WithWithCredentials(&secrets).
			WithWid(workspace), nil); err != nil {
		return nil, err
	}

	return ok.Payload, nil
}

func (c *Client) PushWorkspaceConfiguration(ctx context.Context, workspace string, input *models.TreeServer) error {
	var (
		mode    = "update"
		out     = models.Rfc7396PatchOperation{}
		decoder *mapstructure.Decoder
		err     error
	)

	if decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &out,
		TagName: "json",
	}); err != nil {
		return err
	}

	if err = decoder.Decode(input); err != nil {
		return err
	}

	if out, err = maputil.Compact(out); err != nil {
		return err
	}

	if _, err = c.acp.Hub.WorkspaceConfiguration.
		PatchWorkspaceConfigRfc7396(workspace_configuration.
			NewPatchWorkspaceConfigRfc7396Params().
			WithContext(ctx).
			WithWid(workspace).
			WithMode(&mode).
			WithPatch(out), nil, func(operation *runtime.ClientOperation) {
			operation.PathPattern = "/workspaces/{wid}/promote/config-rfc7396"
		}); err != nil {
		return err
	}

	return nil
}
