package core

import (
	"context"

	"github.com/SevereCloud/vksdk/v2/events"
)

type CommandHandler func(ctx *context.Context, obj *events.MessageNewObject)

type Command struct {
	Aliases     []string
	Description string
	Handler     CommandHandler
}
