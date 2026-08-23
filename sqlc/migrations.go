package sqlc

import "embed"

//go:embed all:migrations
var MigrationsFS embed.FS
