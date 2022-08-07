package cm

import (
	"fmt"
	"strings"
	"vkbot/core"
	"vkbot/subsystems/rolesystem"

	"github.com/SevereCloud/vksdk/v2/events"
)

const GREETING_DEFAULT = "добро пожаловать!"

func greeting(obj *events.MessageNewObject) {
	if err := cmInit(obj); err != nil {
		core.ReplySimple(obj, err.Error())

		return
	}

	s := core.GetStorage()
	key := fmt.Sprintf("greetings.%d", obj.Message.PeerID)

	args := core.ExtractArguments(obj)
	if len(args) == 0 {
		msg, err := s.Db.Get(s.Ctx, key).Result()
		if msg == "" || err != nil {
			msg = GREETING_DEFAULT
		}

		core.ReplySimple(obj, "текущее приветствие: \""+msg+"\"")

		return
	}

	if rolesystem.GetRole(obj) < rolesystem.ROLE_ADMINISTRATOR {
		core.ReplySimple(obj, core.ERR_ACCESS_DENIED)

		return
	}

	s.Db.Set(s.Ctx, key, strings.Join(args, " "), 0)

	core.ReplySimple(obj, "успешно")
}
