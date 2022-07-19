package ping

import (
	"context"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"пинг", "ping"},
		Description: "проверить работоспособность бота",
		Handler:     handle,
	}
}

func handle(_ *context.Context, obj *events.MessageNewObject) {
	core.ReplySimple(obj, "понг")
}
