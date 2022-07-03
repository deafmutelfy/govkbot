package rolesystem

import (
	"fmt"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func GetRole(obj *events.MessageNewObject) string {
	s := core.GetStorage()

	role, _ := s.Db.Get(s.Ctx, fmt.Sprintf("roles.%d.%d", obj.Message.PeerID, obj.Message.FromID)).Result()

	return role
}
