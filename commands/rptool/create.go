package rptool

import (
	"fmt"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func create(obj *events.MessageNewObject) {
	args := core.ExtractArguments(obj)

	trueargs := strings.Split(strings.ToLower(strings.Join(args, " ")), ": ")
	if len(trueargs) < 2 {
		core.ReplySimple(obj, "ошибка: параметры не указаны или указаны неверно")

		return
	}

	trueargs[0] = strings.Trim(trueargs[0], " ")
	trueargs[1] = strings.Trim(trueargs[1], " ")

	s := core.GetStorage()

	s.Db.Set(s.Ctx, fmt.Sprintf("customrp.%d.%s", obj.Message.FromID, trueargs[0]), trueargs[1], 0)

	core.ReplySimple(obj, "успешно")
}
