package client

import acpclient "github.com/cloudentity/acp-client-go"

type Configuration struct {
	// nolint
	acpclient.Config `json:",inline,squash"`

	Insecure bool `json:"insecure"`
}

var DefaultConfig = func() *Configuration {
	return &Configuration{
		Insecure: false,
		Config: acpclient.Config{
			Scopes: []string{"manage_configuration"},
		},
	}
}
