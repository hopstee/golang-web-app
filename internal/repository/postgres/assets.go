package postgres

import "embed"

//go:embed migrations/*.sql
var assets embed.FS
