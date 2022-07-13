package setrole

import (
	"context"
	"fmt"
	"vkbot/core"
	"vkbot/core/rolesystem"

	"github.com/SevereCloud/vksdk/v2/events"
)

func handle(ctx *context.Context, obj *events.MessageNewObject, targetRole int) {
	id := core.GetMention(obj)

	if id == 0 {
		core.ReplySimple(obj, core.ERR_NO_TARGET)

		return
	}
	if id < 0 {
		core.ReplySimple(obj, "вы не можете менять роли сообществам")

		return
	}

	senderRole := rolesystem.GetRole(obj)
	if senderRole < rolesystem.ROLE_OWNER {
		core.ReplySimple(obj, core.ERR_ACCESS_DENIED)

		return
	}
	if obj.Message.FromID == id {
		core.ReplySimple(obj, "ошибка: вы не можете сменить свою роль")

		return
	}

	s := core.GetStorage()
	key := fmt.Sprintf("roles.%d.%d", obj.Message.PeerID, id)

	if targetRole != rolesystem.ROLE_MEMBER {
		s.Db.Set(s.Ctx, key, targetRole, 0)
	} else {
		s.Db.Del(s.Ctx, key)
	}

	core.ReplySimple(obj, "успешно")
}
