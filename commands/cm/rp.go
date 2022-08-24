package cm

import (
	"fmt"
	"vkbot/core"
	"vkbot/subsystems/rolesystem"

	"github.com/SevereCloud/vksdk/v2/events"
)

func rp(obj *events.MessageNewObject) (err error) {
	if err = cmInit(obj); err != nil {
		core.ReplySimple(obj, err.Error())

		return
	}

	if rolesystem.GetRole(obj) < rolesystem.ROLE_OWNER {
		core.ReplySimple(obj, core.ERR_ACCESS_DENIED)

		return
	}

	s := core.GetStorage()
	key := fmt.Sprintf("rp.%d.enabled", obj.Message.PeerID)

	status, _ := s.Db.Get(s.Ctx, key).Result()

	msgstatus := ""

	if status != "false" {
		status = "false"
		msgstatus = "отключен"
	} else {
		status = "true"
		msgstatus = "включен"
	}

	s.Db.Set(s.Ctx, key, status, 0)

	core.ReplySimple(obj, "успешно. RP команды теперь "+msgstatus+"ы")

	return
}
