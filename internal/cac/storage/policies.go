package storage

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"path/filepath"
)

func StorePolicies(policies models.TreePolicies, path string) error {
	for id, policy := range policies {
		var (
			sc   = NewWithID(id, policy)
			name = normalize(policy.PolicyName)
			err  error
		)

		if policy.Language == "rego" {
			var (
				fname = name + ".rego"
				raw   Writer[[]byte]
			)

			if raw, err = RawWriter(path); err != nil {
				return err
			}

			if err = raw(fname, []byte(sc.Other.Definition)); err != nil {
				return err
			}

			sc.Other.Definition = createMultilineIncludeTemplate(fname, 2)
		}

		if err = writeFile(sc, filepath.Join(path, name)); err != nil {
			return err
		}
	}

	return nil
}
