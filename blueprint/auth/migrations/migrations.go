package migrations

import (
	"github.com/sergeyglazyrindev/uadmin/core"
)

var BMigrationRegistry *core.MigrationRegistry

func init() {
	BMigrationRegistry = core.NewMigrationRegistry()

	BMigrationRegistry.AddMigration(initial1623582035{})
	// placeholder to insert next migration
}
