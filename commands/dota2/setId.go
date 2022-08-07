package dota2

import (
	"fmt"
	"strconv"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func setId(obj *events.MessageNewObject) {
	args := core.ExtractArguments(obj)
	if len(args) < 1 {
		core.ReplySimple(obj, "ошибка: необходимо указать свой ID (циферный идентификатор, его можно взять непосредственно в игре, либо на Dotabuff)")

		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		core.ReplySimple(obj, "ошибка: недопустимый ID")

		return
	}

	s := core.GetStorage()

	s.Db.Set(s.Ctx, fmt.Sprintf("dota2.%d.id", obj.Message.FromID), id, 0)

	core.ReplySimple(obj, "успешно")
}
