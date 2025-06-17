package storage

import (
	"os"
	"path/filepath"
)

// var files []string

//             for _, dir := range st.Config.DirPath {
//                 err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
//                     if err != nil {
//                         return err
//                     }

//                     if !info.IsDir() {
//                         if path, err = filepath.Rel(dir, path); err != nil {
//                             return err
//                         }

//                         files = append(files, path)
//                     }
//                     return nil
//                 })
//             }

// write a function that returns a list of all files in the input directories
func ListFilesInDirectories(dirs ...string) ([]string, error) {
	var files []string

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				relPath, err := filepath.Rel(dir, path)
				if err != nil {
					return err
				}
				files = append(files, relPath)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}