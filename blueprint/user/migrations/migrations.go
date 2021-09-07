package migrations

import (
	"github.com/sergeyglazyrindev/uadmin/core"
)

var BMigrationRegistry *core.MigrationRegistry

func init() {
	BMigrationRegistry = core.NewMigrationRegistry()

	BMigrationRegistry.AddMigration(initial1621680132{})

	BMigrationRegistry.AddMigration(addinguse1623259185{})

	// placeholder to insert next migration
}
