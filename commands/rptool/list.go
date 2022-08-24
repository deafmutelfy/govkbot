package rptool

import (
	"fmt"
	"strconv"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func list(obj *events.MessageNewObject) (err error) {
	s := core.GetStorage()

	keys, err := s.Db.Keys(s.Ctx, fmt.Sprintf("customrp.%d.*", obj.Message.FromID)).Result()
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	msg := "ваши РП-команды:\n"
	leastOne := false

	for _, x := range keys {
		leastOne = true

		msg += "- " + strings.Replace(x, "customrp."+strconv.Itoa(obj.Message.FromID)+".", "", 1) + "\n"
	}

	if !leastOne {
		msg = "ошибка: вы ещё не создали ни одной РП-команды"
	}

	core.ReplySimple(obj, msg)

	return
}
