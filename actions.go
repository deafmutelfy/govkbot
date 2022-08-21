package main

import (
	"fmt"
	"strconv"
	"strings"
	"vkbot/commands/cm"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/dlclark/regexp2"
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

func handleUserRPAction(obj *events.MessageNewObject) {
	s := core.GetStorage()

	msg := strings.ToLower(obj.Message.Text)

	id := core.GetMention(obj)
	if id <= 0 {
		return
	}

	enabled, _ := s.Db.Get(s.Ctx, fmt.Sprintf("rp.%d.enabled", obj.Message.PeerID)).Result()
	if enabled == "false" {
		return
	}

	action, err := s.Db.Get(s.Ctx, fmt.Sprintf("customrp.%d.%s", obj.Message.FromID, msg)).Result()
	if err != nil || action == "" {
		return
	}

	n1 := core.GetNicknameOrFullName(obj.Message.FromID)
	n2 := core.GetNicknameOrFullName(id)

	action, err = regexp2.MustCompile(`(?i)\bя\b`, 0).Replace(action, "[id" +
				strconv.Itoa(obj.Message.FromID) +
				"|" +
				n1 +
				"]", 0, -1)
	if err != nil {
		return
	}
	action, err = regexp2.MustCompile(`(?i)\bцель\b`, 0).Replace(action, "[id" +
				strconv.Itoa(id) +
				"|" +
				n2 +
				"]", 0, -1)
	if err != nil {
		return
	}

	core.SendSimple(obj, action)
}
