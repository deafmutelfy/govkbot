package rptool

import (
	"fmt"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func create(obj *events.MessageNewObject) (err error) {
	args := core.ExtractArguments(obj)

	trueargs := strings.Split(strings.Join(args, " "), ": ")
	if len(trueargs) < 2 {
		core.ReplySimple(obj, "ошибка: параметры не указаны или указаны неверно")

		return
	}

	trueargs[0] = strings.Trim(trueargs[0], " ")
	trueargs[1] = strings.Trim(strings.Join(trueargs[1:], " "), " ")

	s := core.GetStorage()

	s.Db.Set(s.Ctx, fmt.Sprintf("customrp.%d.%s", obj.Message.FromID, trueargs[0]), trueargs[1], 0)

	core.ReplySimple(obj, "успешно")

	return
}
