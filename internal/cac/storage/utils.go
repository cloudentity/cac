package storage

import "github.com/cloudentity/acp-client-go/clients/hub/models"

func readFileToMap(server models.Rfc7396PatchOperation, key string, path string) error {
	var err error

	if server[key], err = readFile(path); err != nil {
		return err
	}

	if v, ok := server[key].(map[string]any); ok && len(v) == 0 {
		delete(server, key)
	}

	return nil
}

func readFilesToMap(server models.Rfc7396PatchOperation, key string, path string) error {
	var err error

	if server[key], err = readFiles(path); err != nil {
		return err
	}

	if v, ok := server[key].(map[string]any); ok && len(v) == 0 {
		delete(server, key)
	}

	return nil
}
