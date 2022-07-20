package main

import (
	"fmt"
	"vkbot/commands/cm"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func handleChatInviteUser(obj *events.MessageNewObject) {
	s := core.GetStorage()

	if obj.Message.Action.MemberID == -s.Cfg.GroupId {
		return
	}

	msg, err := s.Db.Get(s.Ctx, fmt.Sprintf("greetings.%d", obj.Message.PeerID)).Result()
	if msg == "" || err != nil {
		msg = cm.GREETING_DEFAULT
	}

	obj.Message.FromID = obj.Message.Action.MemberID
	core.ReplySimple(obj, msg)
}
