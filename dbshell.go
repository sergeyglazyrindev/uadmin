package uadmin

import (
	"github.com/sergeyglazyrindev/uadmin/core"
)

type DbShellCommand struct {
}

func (c DbShellCommand) Proceed(subaction string, args []string) error {
	adapter := core.NewDbAdapter(nil, core.CurrentDatabaseSettings.Default.Type)
	return adapter.StartDBShell(core.CurrentDatabaseSettings.Default)
}

func (c DbShellCommand) GetHelpText() string {
	return "Start shell for your database"
}
