package storage

import (
	"fmt"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
)

func Register(cfg Config) (map[defs.StorageCode]Storage, map[defs.StorageCode]bool, error) {
	all := []Storage{
		NewFile(cfg.File),
	}

	result := make(map[defs.StorageCode]Storage, len(all))
	codes := make(map[defs.StorageCode]bool, len(all))
	for _, current := range all {
		if _, ok := result[current.Code()]; ok {
			return nil, nil, fmt.Errorf("storage %q is registered twice", current.Code())
		}
		result[current.Code()] = current
		codes[current.Code()] = true
	}
	return result, codes, nil
}
