package fileutil

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Vfs is virtual file system
var Vfs = afero.NewOsFs()

// SetMemMapFs changes Vfs to NewMemMapFs
func SetMemMapFs() {
	Vfs = afero.NewMemMapFs()
}

// FolderExists ensures that folder exists
func FolderExists(dir string) error {
	if dir == "" {
		return errors.Errorf("invalid parameter: dir")
	}

	stat, err := Vfs.Stat(dir)
	if err != nil {
		return errors.WithStack(err)
	}

	if !stat.IsDir() {
		return errors.Errorf("not a folder: %q", dir)
	}

	return nil
}

// FileExists ensures that file exists
func FileExists(file string) error {
	if file == "" {
		return errors.Errorf("invalid parameter: file")
	}

	stat, err := Vfs.Stat(file)
	if err != nil {
		return errors.WithStack(err)
	}

	if stat.IsDir() {
		return errors.Errorf("not a file: %q", file)
	}

	return nil
}
