package ping

import (
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

func handle(obj *events.MessageNewObject) (err error) {
	core.ReplySimple(obj, "понг")

	return
}
