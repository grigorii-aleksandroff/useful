package migrations

import "embed"

func FS() embed.FS {
	return migrations
}

//go:embed *.sql
var migrations embed.FS
