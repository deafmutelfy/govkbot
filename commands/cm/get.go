package cm

import (
	"fmt"
	"vkbot/core"
	"vkbot/subsystems/rolesystem"

	"github.com/SevereCloud/vksdk/v2/events"
)

func getrole(obj *events.MessageNewObject) (err error) {
	if err = cmInit(obj); err != nil {
		core.ReplySimple(obj, err.Error())

		return
	}
	
	s := core.GetStorage()

	initialized, err := s.Db.Get(s.Ctx, fmt.Sprintf("roles.%d.initialized", obj.Message.PeerID)).Result()
	if initialized != "true" {
		core.ReplySimple(obj, core.ERR_NO_ROLESYSTEM)

		return
	}

	role := ""

	switch rolesystem.GetRole(obj) {
	case rolesystem.ROLE_BOT_OWNER:
		role = core.GetAlias(obj.Message.FromID)
	case rolesystem.ROLE_OWNER:
		role = "владелец"
	case rolesystem.ROLE_ADMINISTRATOR:
		role = "администратор"
	case rolesystem.ROLE_MODERATOR:
		role = "модератор"
	case rolesystem.ROLE_MEMBER:
		role = "участник"
	}

	core.ReplySimple(obj, "ваша роль: "+role)

	return
}
