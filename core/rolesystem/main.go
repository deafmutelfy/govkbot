package rolesystem

import (
	"fmt"
	"strconv"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func GetRole(obj *events.MessageNewObject) int {
	s := core.GetStorage()

	if strconv.Itoa(obj.Message.FromID) == s.Cfg.BotOwnerId {
		return ROLE_BOT_OWNER
	}

	role, _ := s.Db.Get(s.Ctx, fmt.Sprintf("roles.%d.%d", obj.Message.PeerID, obj.Message.FromID)).Result()
	rint, _ := strconv.Atoi(role)

	return rint
}
