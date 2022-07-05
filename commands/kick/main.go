package kick

import (
	"context"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"кик"},
		Description: "исключить участника беседы",
		Handler:     handle,
	}
}

func handle(ctx *context.Context, obj *events.MessageNewObject) {
	id := core.GetMention(obj)

	if id == 0 {
		core.ReplySimple(obj, "ошибка: не указана цель. Перешлите сообщение цели, или укажите её с помощью упоминания")

		return
	}

}
