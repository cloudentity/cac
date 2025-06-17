package storage_test

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/diff"
	"github.com/cloudentity/cac/internal/cac/logging"
	"github.com/cloudentity/cac/internal/cac/storage"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"
)

func TestTenantStorage(t *testing.T) {
	tcs := []struct {
		desc    string
		data    *models.TreeTenant
		files   []string
		filters []string
		assert  func(t *testing.T, path string, bts []byte)
	}{
		{
			desc: "workspace and mfa_methods",
			data: &models.TreeTenant{
				Servers: models.TreeServers{
					"demo": models.TreeServer{
						Name:           "demo workspace",
						AccessTokenTTL: strfmt.Duration(time.Minute * 10),
						Idps: models.TreeIDPs{
							"oidc": models.TreeIDP{
								Name:     "oidc",
								Disabled: true,
							},
						},
					},
				},
				MfaMethods: models.TreeMFAMethods{
					"sms": models.TreeMFAMethod{
						Enabled:   true,
						Mechanism: "sms",
					},
				},
			},
			files: []string{
				"mfa_methods/sms.yaml",
				"workspaces/demo/server.yaml",
				"workspaces/demo/idps/oidc.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				switch path {
				case "mfa_methods/sms.yaml":
					require.YAMLEq(t, `enabled: true
id: sms
mechanism: sms`, string(bts))
				case "workspaces/demo/server.yaml":
					require.YAMLEq(t, `access_token_ttl: 10m0s
authentication_mechanisms: []
authorization_code_ttl: 0s
scope_claim_formats: []
backchannel_token_delivery_modes_supported: []
backchannel_user_code_parameter_supported: false
cookie_max_age: 0s
do_not_create_default_claims: false
enable_idp_discovery: false
enable_legacy_clients_with_no_software_statement: false
enable_quick_access: false
enable_trust_anchor: false
enforce_id_token_encryption: false
enforce_pkce: false
enforce_pkce_for_public_clients: false
grant_types: []
id: demo
id_token_ttl: 0s
initialize: false
name: demo workspace
pushed_authorization_request_ttl: 0s
refresh_token_ttl: 0s
require_pushed_authorization_requests: false
rotated_secrets: []
subject_identifier_types: []
template: false
tenant_id: ""
token_endpoint_auth_methods: []
token_endpoint_auth_signing_alg_values: []
token_endpoint_authn_methods: []
version: 0`, string(bts))
				case "workspaces/demo/idps/oidc.yaml":
					require.YAMLEq(t, `disabled: true
display_order: 0
hidden: false
id: oidc
name: oidc
static_amr: []
version: 0`, string(bts))
				}
			},
		},
		{
			desc:    "filtered workspace and mfa_methods",
			filters: []string{"mfa_methods"},
			data: &models.TreeTenant{
				Servers: models.TreeServers{
					"demo": models.TreeServer{
						Name:           "demo workspace",
						AccessTokenTTL: strfmt.Duration(time.Minute * 10),
						Idps: models.TreeIDPs{
							"oidc": models.TreeIDP{
								Name:     "oidc",
								Disabled: true,
							},
						},
					},
				},
				MfaMethods: models.TreeMFAMethods{
					"sms": models.TreeMFAMethod{
						Enabled:   true,
						Mechanism: "sms",
					},
				},
			},
			files: []string{
				"mfa_methods/sms.yaml",
				"workspaces/demo/server.yaml",
				"workspaces/demo/idps/oidc.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				switch path {
				case "mfa_methods/sms.yaml":
					require.YAMLEq(t, `enabled: true
id: sms
mechanism: sms`, string(bts))
				}
			},
		},
		{
			desc: "themes and templates",
			data: &models.TreeTenant{
				Themes: models.TreeThemes{
					"theme1": models.TreeTheme{
						Name: "theme1",
						Templates: models.TreeTemplates{
							"pages/error/index.tmpl": models.TreeTemplate{
								Content: "template1 content",
							},
							"shared/footer.tmpl": models.TreeTemplate{
								Content: "footer content",
							},
						},
					},
				},
			},
			files: []string{
				"themes/theme1/theme.yaml",
				"themes/theme1/templates/pages_error_index.tmpl",
				"themes/theme1/templates/pages_error_index.tmpl.yaml",
				"themes/theme1/templates/shared_footer.tmpl",
				"themes/theme1/templates/shared_footer.tmpl.yaml",
			},
			assert: func(t *testing.T, path string, bts []byte) {
				switch path {
				case "themes/theme1/theme.yaml":
					require.YAMLEq(t, `name: theme1`, string(bts))
				case "themes/theme1/templates/pages_error_index.tmpl":
					require.Equal(t, "template1 content", string(bts))
				case "themes/theme1/templates/shared_footer.tmpl":
					require.Equal(t, "footer content", string(bts))
				case "themes/theme1/templates/pages/error_index.tmpl.yaml":
					require.Equal(t, `content: {{ include "error_index.tmpl" | nindent 2 }}
created_at: "0001-01-01T00:00:00.000Z"
id: pages/error/index.tmpl
updated_at: "0001-01-01T00:00:00.000Z"`, string(bts))
				case "themes/theme1/templates/shared_footer.tmpl.yaml":
					require.Equal(t, `id: shared/footer.tmpl
content: {{ include "shared_footer.tmpl" | nindent 2 }}
created_at: 0001-01-01T00:00:00.000Z
updated_at: 0001-01-01T00:00:00.000Z
`, string(bts))
				}
			},
		},
	}

	for _, tc := range tcs {
        if tc.desc != "workspace and mfa_methods" {
            continue
        }
		t.Run(tc.desc, func(t *testing.T) {
			err := logging.InitLogging(&logging.Configuration{
				Level: "debug",
			})

			require.NoError(t, err)

			st, err := storage.InitMultiStorage(&storage.MultiStorageConfiguration{
				DirPath: []string{t.TempDir()},
			}, storage.InitTenantStorage)

			require.NoError(t, err)

			patchData, err := utils.FromModelToPatch(tc.data)
			require.NoError(t, err)

			err = st.Write(context.Background(), &api.TenantPatch{
				Data: patchData,
				Ext:  &api.TenantExtensions{},
			}, api.WithWorkspace("demo"))
			require.NoError(t, err)

			var files []string

			for _, dir := range st.Config.DirPath {
				err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if !info.IsDir() {
						if path, err = filepath.Rel(dir, path); err != nil {
							return err
						}

						files = append(files, path)
					}
					return nil
				})
			}

			require.NoError(t, err)
			require.ElementsMatch(t, tc.files, files)

			// checking if files written to fs have expected content
			for _, f := range tc.files {
				// using first dirpath as multi storage stores everything there
				bts, err := os.ReadFile(filepath.Join(st.Config.DirPath[0], f))
				require.NoError(t, err)

				if tc.assert != nil {
					tc.assert(t, f, bts)
				}
			}

			var readServer api.Patch
			readServer, err = st.Read(context.Background(),
				api.WithWorkspace("demo"),
				api.WithFilters(tc.filters))

			require.NoError(t, err)

			// verifying if the data read from fs is the same as the provided test data
			patchData, err = utils.FilterPatch(patchData, tc.filters)
			require.NoError(t, err)

			d, err := diff.Tree(&api.TenantPatch{
				Data: patchData,
				Ext:  &api.TenantExtensions{},
			}, readServer)
			require.NoError(t, err)
			require.Empty(t, d)
		})
	}
}
