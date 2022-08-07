package help

import (
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

const doc_url = "https://vkbot.deafmute.xyz"

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"хелп", "помощь"},
		Description: "получить справку о командах бота",
		Handler:     handle,
	}
}

func handle(obj *events.MessageNewObject) {
	core.ReplySimple(obj, "документация по боту находится здесь: "+doc_url)
}
