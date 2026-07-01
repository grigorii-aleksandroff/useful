package storage

import (
	"os"
	"path/filepath"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
)

type File struct {
	dir          string
	downloadType defs.DownloadType
}

func NewFile(cfg FileConfig) File {
	downloadType := cfg.DownloadType
	if downloadType == "" {
		downloadType = defs.DownloadTypeFile
	}

	return File{dir: cfg.Dir, downloadType: downloadType}
}

func (File) Code() defs.StorageCode {
	return defs.StorageFile
}

func (f File) DownloadType() defs.DownloadType {
	return f.downloadType
}

func (f File) Save(name string, content []byte) (string, error) {
	dir := f.dir
	if dir == "" {
		dir = "/tmp/reporter"
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	full := filepath.Join(dir, name)
	if err := os.WriteFile(full, content, 0o644); err != nil {
		return "", err
	}

	return full, nil
}

func (f File) Download(locator string) (Download, error) {
	switch f.downloadType {
	case defs.DownloadTypeLink:
		return Download{Type: defs.DownloadTypeLink, Link: locator}, nil
	default:
		content, err := os.ReadFile(locator)
		if err != nil {
			return Download{}, err
		}

		return Download{Type: defs.DownloadTypeFile, Content: content}, nil
	}
}

var _ Storage = File{}
