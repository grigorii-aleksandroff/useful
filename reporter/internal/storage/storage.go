package storage

import "git.bububla.com/kilogramix/asia/reporter.git/internal/defs"

type Download struct {
	Type    defs.DownloadType
	Content []byte
	Link    string
}

type Storage interface {
	Code() defs.StorageCode
	DownloadType() defs.DownloadType
	Save(name string, content []byte) (string, error)
	Download(locator string) (Download, error)
}
