package rptool

import (
	"fmt"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func clear(obj *events.MessageNewObject) (err error) {
	s := core.GetStorage()

	keys, err := s.Db.Keys(s.Ctx, fmt.Sprintf("customrp.%d.*", obj.Message.FromID)).Result()
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}
	if len(keys) == 0 {
		core.ReplySimple(obj, "ошибка: вы ещё не создали ни одной РП-команды")

		return
	}

	for _, x := range keys {
		s.Db.Del(s.Ctx, x)
	}

	core.ReplySimple(obj, "успешно")

	return
}
