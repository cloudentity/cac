package storage

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"path/filepath"
)

func storeScripts(scripts models.TreeScripts, path string) error {
	for id, script := range scripts {
		var (
			sc   = NewWithID(id, script)
			name = normalize(script.Name)
			jsn  = name + ".js"
			raw  Writer[[]byte]
			err  error
		)

		if raw, err = RawWriter(path); err != nil {
			return err
		}

		if err = raw(jsn, []byte(sc.Other.Body)); err != nil {
			return err
		}

		sc.Other.Body = createMultilineIncludeTemplate(jsn, 2)

		if err = writeFile(sc, filepath.Join(path, name)); err != nil {
			return err
		}
	}

	return nil
}
