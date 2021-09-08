package uadmin

import (
	"fmt"
	"github.com/sergeyglazyrindev/uadmin/core"
	"strings"
)

type CommandRegistry struct {
	Actions map[string]core.ICommand
}

func (r CommandRegistry) AddAction(name string, command core.ICommand) {
	r.Actions[name] = command
}

func (r CommandRegistry) IsRegisteredCommand(name string) bool {
	_, err := r.Actions[name]
	return !!err
}

func (r CommandRegistry) RunAction(command string, subaction string, args []string) error {
	action, _ := r.Actions[command]
	return action.Proceed(subaction, args)
}

func (r CommandRegistry) MakeHelpText() string {
	var helpParts []string
	var i = 1
	for action, handler := range r.Actions {
		helpParts = append(helpParts, fmt.Sprintf("%d. %s - %s", i, action, handler.GetHelpText()))
		i++
	}
	return strings.Join(helpParts, "\n")
}
