package storage

import "git.bububla.com/kilogramix/asia/reporter.git/internal/defs"

type Config struct {
	File FileConfig
}

type FileConfig struct {
	Dir          string
	DownloadType defs.DownloadType
}
