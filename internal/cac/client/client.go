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
	"github.com/go-openapi/runtime"
	"github.com/pkg/errors"
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

func (c *Client) Read(ctx context.Context, workspace string, opts ...api.SourceOpt) (models.Rfc7396PatchOperation, error) {
	var (
		options = &api.Options{}
		ok      *workspace_configuration.ExportWorkspaceConfigOK
		data    models.Rfc7396PatchOperation
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

	if data, err = utils.FromTreeServerToPatch(ok.Payload); err != nil {
		return nil, errors.Wrap(err, "failed to convert tree server to patch")
	}

	return data, nil
}

func (c *Client) Write(ctx context.Context, workspace string, data models.Rfc7396PatchOperation) error {
	var (
		mode = "update"
		err  error
	)

	if _, err = c.acp.Hub.WorkspaceConfiguration.
		PatchWorkspaceConfigRfc7396(workspace_configuration.
			NewPatchWorkspaceConfigRfc7396Params().
			WithContext(ctx).
			WithWid(workspace).
			WithMode(&mode).
			WithPatch(data), nil, func(operation *runtime.ClientOperation) {
			operation.PathPattern = "/workspaces/{wid}/promote/config-rfc7396"
		}); err != nil {
		return err
	}

	return nil
}

func (c *Client) String() string {
	return fmt.Sprintf("client: %v", c.acp.Config.IssuerURL)
}
