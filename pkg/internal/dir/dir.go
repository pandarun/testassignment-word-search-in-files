package dir

import (
	"io/fs"
	"path/filepath"
)

// FilesFS returns a list of files in the directory
func FilesFS(fsys fs.FS, dir string) ([]string, error) {
	if dir == "" {
		dir = "."
	}
	var fileNames []string
	err := fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileNames = append(fileNames, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileNames, nil
}

type FileNameWithoutExtension struct {
	value string
}

func (f FileNameWithoutExtension) String() string {
	return f.value
}

func NewFileNameWithoutExtension(fileName string) FileNameWithoutExtension {
	return FileNameWithoutExtension{
		value: fileName[:len(fileName)-len(filepath.Ext(fileName))],
	}
}
