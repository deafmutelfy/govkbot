package core

import (
	"github.com/SevereCloud/vksdk/v2/events"
)

type CommandHandler func(obj *events.MessageNewObject)

type CommandQueueParams struct {
	Name    string
	Handler CommandHandler
}

type Command struct {
	Aliases     []string
	Description string
	Handler     CommandHandler
	Metacommand bool
	Subcommands *[]Command
	Hidden      bool
	NoPrefix    bool
	Queue       *CommandQueueParams
}
