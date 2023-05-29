package cm

import (
	"strconv"
	"vkbot/core"
	"vkbot/subsystems/rolesystem"

	"github.com/SevereCloud/vksdk/v2/events"
)

func unban(obj *events.MessageNewObject) (err error) {
	if err = cmInit(obj); err != nil {
		core.ReplySimple(obj, err.Error())

		return
	}

	id := core.GetMention(obj)
	if id == 0 {
		core.ReplySimple(obj, core.ERR_NO_TARGET)

		return
	}

	if rolesystem.GetRole(obj) < rolesystem.ROLE_MODERATOR {
		core.ReplySimple(obj, core.ERR_ACCESS_DENIED)

		return
	}

	s := core.GetStorage()

	s.Db.Del(s.Ctx, "bans."+strconv.Itoa(obj.Message.PeerID)+"."+strconv.Itoa(id))

	core.ReplySimple(obj, "успешно")

	return
}
