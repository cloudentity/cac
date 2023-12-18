package storage_test

import (
	"encoding/json"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/storage"
	"github.com/go-openapi/strfmt"
	"github.com/imdario/mergo"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	tcs := []struct {
		desc   string
		data   *models.TreeServer
		files  []string
		assert func(t *testing.T, path string, bts []byte)
	}{
		{
			desc: "server",
			data: &models.TreeServer{
				Name:           "demo workspace",
				AccessTokenTTL: strfmt.Duration(10 * time.Minute),
			},
			files: []string{
				"workspaces/demo/server.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `id: demo
name: demo workspace
access_token_ttl: 10m0s`, string(bts))
			},
		},
		{
			desc: "clients",
			data: &models.TreeServer{
				Clients: models.TreeClients{
					"demo-demo": models.TreeClient{
						ClientName: "Demo Portal",
					},
				},
			},
			files: []string{
				"workspaces/demo/clients/Demo_Portal.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `id: demo-demo
client_name: Demo Portal`, string(bts))
			},
		},
		{
			desc: "idps",
			data: &models.TreeServer{
				Idps: models.TreeIDPs{
					"some-idp": models.TreeIDP{
						Name: "Some IDP",
					},
				},
			},
			files: []string{
				"workspaces/demo/idps/Some_IDP.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `id: some-idp
name: Some IDP`, string(bts))
			},
		},
		{
			desc: "claims",
			data: &models.TreeServer{
				Claims: models.TreeClaims{
					"access_token": models.TreeClaimType{
						"customer_id": models.TreeClaim{
							Mapping:    "customer_id",
							Scopes:     []string{"customer"},
							SourcePath: "customer_id",
							SourceType: "authnCtx",
						},
					},
				},
			},
			files: []string{
				"workspaces/demo/claims.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `access_token:
  customer_id:
    mapping: customer_id
    scopes: [customer]
    source_path: customer_id
    source_type: authnCtx`, string(bts))
			},
		},
		{
			desc: "custom_apps",
			data: &models.TreeServer{
				CustomApps: models.TreeCustomApps{
					"some-app": models.TreeCustomApp{
						Name: "Some App",
						URL:  "https://some-app.com",
					},
				},
			},
			files: []string{
				"workspaces/demo/custom_apps/Some_App.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `id: some-app
name: Some App
url: https://some-app.com`, string(bts))
			},
		},
		{
			desc: "gateways",
			data: &models.TreeServer{
				Gateways: models.TreeGateways{
					"some-gateway": models.TreeGateway{
						Name: "Some Gateway",
					},
				},
			},
			files: []string{
				"workspaces/demo/gateways/Some_Gateway.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `id: some-gateway
name: "Some Gateway"`, string(bts))
			},
		},
		{
			desc: "policy_execution_points",
			data: &models.TreeServer{
				PolicyExecutionPoints: models.TreePolicyExecutionPoints{
					"server_user_token": "some_policy_id",
				},
			},
			files: []string{
				"workspaces/demo/policy_execution_points.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `server_user_token: some_policy_id`, string(bts))
			},
		},
		{
			desc: "pools",
			data: &models.TreeServer{
				Pools: models.TreePools{
					"some-pool": models.TreePool{
						Name: "Some Pool",
					},
				},
			},
			files: []string{
				"workspaces/demo/pools/Some_Pool.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `id: some-pool
name: "Some Pool"`, string(bts))
			},
		},
		{
			desc: "scopes",
			data: &models.TreeServer{
				ScopesWithoutService: models.TreeScopes{
					"some_scope": models.TreeScope{
						Description: "Some Scope",
						PolicyExecutionPoints: models.TreePolicyExecutionPoints{
							"scope_user_grant": "some_policy_id",
						},
						Transient: false,
					},
				},
			},
			files: []string{
				"workspaces/demo/scopes.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `some_scope:
  description: Some Scope
  policy_execution_points: 
    scope_user_grant: some_policy_id`, string(bts))
			},
		},
		{
			desc: "script_execution_points",
			data: &models.TreeServer{
				ScriptExecutionPoints: models.TreeScriptExecutionPoints{
					"token_minting": {
						"demo": models.TreeScriptExecutionPoint{
							ScriptID: "some_script_id",
						},
					},
				},
			},
			files: []string{
				"workspaces/demo/script_execution_points.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `token_minting: 
  demo:
    script_id: some_script_id`, string(bts))
			},
		},
		{
			desc: "consent",
			data: &models.TreeServer{
				ServerConsent: &models.TreeServerConsent{
					Custom: &models.CustomServerConsent{
						ServerConsentURL: "https://example.com/consent",
					},
					Type: "custom",
				},
			},
			files: []string{
				"workspaces/demo/consent.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `custom:
  server_consent_url: https://example.com/consent
type: custom`, string(bts))
			},
		},
		{
			desc: "server bindings",
			data: &models.TreeServer{
				ServersBindings: models.TreeServersBindings{
					"other_server": true,
				},
			},
			files: []string{
				"workspaces/demo/server_bindings.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `bindings:
  - other_server`, string(bts))
			},
		},
		{
			desc: "services",
			data: &models.TreeServer{
				Services: models.TreeServices{
					"some_service": models.TreeService{
						Name: "Some Service",
						Scopes: models.TreeScopes{
							"some_scope": models.TreeScope{
								Description: "Some Scope",
								PolicyExecutionPoints: models.TreePolicyExecutionPoints{
									"scope_user_grant": "some_policy_id",
								},
							},
						},
					},
				},
			},
			files: []string{
				"workspaces/demo/services/Some_Service.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `id: some_service
name: Some Service
scopes: 
  some_scope:
    description: Some Scope
    policy_execution_points:
      scope_user_grant: some_policy_id`, string(bts))
			},
		},
		{
			desc: "theme binding",
			data: &models.TreeServer{
				ThemeBinding: &models.TreeThemeBinding{
					ThemeID: "some_theme",
				},
			},
			files: []string{
				"workspaces/demo/theme_binding.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `theme_id: some_theme`, string(bts))
			},
		},
		{
			desc: "webhooks",
			data: &models.TreeServer{
				Webhooks: models.TreeWebhooks{
					"hook_id": models.TreeWebhook{
						Active: true,
						URL:    "https://example.com",
					},
				},
			},
			files: []string{
				"workspaces/demo/webhooks/hook_id.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `id: hook_id
active: true
url: "https://example.com"`, string(bts))
			},
		},
		{
			desc: "rego policies",
			data: &models.TreeServer{
				Policies: models.TreePolicies{
					"some_policy": models.TreePolicy{
						Definition: `
package acp.authz

default allow = false
`,
						Language:   "rego",
						PolicyName: "Some Rego Policy",
						Type:       "api",
					},
				},
			},
			files: []string{
				"workspaces/demo/policies/Some_Rego_Policy.yaml",
				"workspaces/demo/policies/Some_Rego_Policy.rego",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				if strings.Contains(path, ".yaml") {
					require.Equal(t, `definition: |-
  {{ include "Some_Rego_Policy.rego" | indent 2 }}
id: some_policy
language: rego
policy_name: Some Rego Policy
type: api
`, string(bts))
				} else {
					require.Equal(t, `
package acp.authz

default allow = false
`, string(bts))
				}
			},
		},
		{
			desc: "ce policies",
			data: &models.TreeServer{
				Policies: models.TreePolicies{
					"some_policy": models.TreePolicy{
						Language:   "cloudentity",
						PolicyName: "Some CE Policy",
						Type:       "api",
						Validators: []*models.ValidatorConfig{
							{
								Conf: map[string]any{
									"fields": []map[string]any{
										{
											"comparator": "contains",
											"field":      "login.verified_recovery_methods",
											"value": []string{
												"mfa",
											},
										},
									},
								},
								Recovery: []*models.RecoveryConfig{
									{
										Type: "mfa",
									},
								},
							},
						},
					},
				},
			},
			files: []string{
				"workspaces/demo/policies/Some_CE_Policy.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				require.YAMLEq(t, `id: some_policy
language: cloudentity
type: api
policy_name: "Some CE Policy"
validators:
  - recovery:
      - type: mfa
    conf:
      fields:
        - comparator: contains
          field: login.verified_recovery_methods
          value:
          - mfa`, string(bts))
			},
		},
		{
			desc: "js extensions",
			data: &models.TreeServer{
				Scripts: models.TreeScripts{
					"some_script": models.TreeScript{
						Body: `
module.exports = async function(context) {
      return {
        access_token: {
          x: "123"
        }
      };
  }
`,
						Name: "Some Script",
					},
				},
			},
			files: []string{
				"workspaces/demo/scripts/Some_Script.yaml",
				"workspaces/demo/scripts/Some_Script.js",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				if strings.Contains(path, ".yaml") {
					require.Equal(t, `body: |-
  {{ include "Some_Script.js" | indent 2 }}
id: some_script
name: Some Script
`, string(bts))
				} else {
					require.Equal(t, `
module.exports = async function(context) {
      return {
        access_token: {
          x: "123"
        }
      };
  }
`, string(bts))
				}

			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			st := storage.InitStorage(storage.Configuration{
				DirPath: t.TempDir(),
			})

			err := st.Store("demo", tc.data)
			require.NoError(t, err)

			var files []string
			err = filepath.Walk(st.Config.DirPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					if path, err = filepath.Rel(st.Config.DirPath, path); err != nil {
						return err
					}

					files = append(files, path)
				}
				return nil
			})

			require.NoError(t, err)
			require.ElementsMatch(t, slices.Compact(append(tc.files, "workspaces/demo/server.yaml")), files)

			for _, f := range tc.files {
				bts, err := os.ReadFile(filepath.Join(st.Config.DirPath, f))
				require.NoError(t, err)

				if tc.assert != nil {
					tc.assert(t, f, bts)
				}
			}

			// fill the server data to make the comparison easier
			err = mergo.Merge(tc.data, models.TreeServer{
				Clients:              models.TreeClients{},
				Idps:                 models.TreeIDPs{},
				CustomApps:           models.TreeCustomApps{},
				Gateways:             models.TreeGateways{},
				Pools:                models.TreePools{},
				Webhooks:             models.TreeWebhooks{},
				Scripts:              models.TreeScripts{},
				Policies:             models.TreePolicies{},
				Services:             models.TreeServices{},
				ScopesWithoutService: models.TreeScopes{},
				ServersBindings:      models.TreeServersBindings{},
			})
			require.NoError(t, err)

			var readServer models.TreeServer
			readServer, err = st.Read("demo")

			require.NoError(t, err)

			serData, err := json.Marshal(tc.data)
			require.NoError(t, err)

			serReadServer, err := json.Marshal(readServer)
			require.NoError(t, err)

			require.Equal(t, string(serData), string(serReadServer))
		})
	}
}
