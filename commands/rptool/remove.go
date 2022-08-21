package rptool

import (
	"fmt"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func remove(obj *events.MessageNewObject) {
	args := strings.Join(core.ExtractArguments(obj), " ")
	if args == "" {
		core.ReplySimple(obj, "ошибка: не указано название команды")

		return
	}

	s := core.GetStorage()

	n, err := s.Db.Del(s.Ctx, fmt.Sprintf("customrp.%d.%s", obj.Message.FromID, args)).Result()
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}
	if n == 0 {
		core.ReplySimple(obj, "ошибка: вы ещё не создали РП-команду с таким названием")

		return
	}

	core.ReplySimple(obj, "успешно")
}
