package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/cloudentity/acp-client-go"
	"github.com/cloudentity/acp-client-go/clients/hub/client/workspace_configuration"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/go-openapi/runtime"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slog"
	"net/http"
)

type Client struct {
	acp *acpclient.Client
}

func WithSecrets(secrets bool) api.SourceOpt {
	return func(o *api.Options) {
		o.Secrets = secrets
	}
}

var _ api.Source = &Client{}

func InitClient(config *Configuration) (c *Client, err error) {
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

func (c *Client) Read(ctx context.Context, workspace string, opts ...api.SourceOpt) (*models.TreeServer, error) {
	var (
		ok      *workspace_configuration.ExportWorkspaceConfigOK
		options = &api.Options{}
		err     error
	)

	for _, opt := range opts {
		opt(options)
	}

	if ok, err = c.acp.Hub.WorkspaceConfiguration.
		ExportWorkspaceConfig(workspace_configuration.
			NewExportWorkspaceConfigParams().
			WithContext(ctx).
			WithWithCredentials(&options.Secrets).
			WithWid(workspace), nil); err != nil {
		return nil, err
	}

	return ok.Payload, nil
}

func (c *Client) Write(ctx context.Context, workspace string, input *models.TreeServer) error {
	var (
		mode    = "update"
		out     = models.Rfc7396PatchOperation{}
		decoder *mapstructure.Decoder
		err     error
	)

	if decoder, err = utils.Decoder(&out); err != nil {
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

func (c *Client) String() string {
	return fmt.Sprintf("client: %v", c.acp.Config.IssuerURL)
}
