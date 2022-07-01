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

func handle(ctx *context.Context, obj *events.MessageNewObject) {
	core.SendSimple(obj.Message.PeerID, "Понг")
}
