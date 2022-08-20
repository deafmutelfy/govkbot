package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func handleUserRPAction(obj *events.MessageNewObject) {
	s := core.GetStorage()

	msg := strings.ToLower(obj.Message.Text)

	id := core.GetMention(obj)
	if id <= 0 {
		return
	}

	action, err := s.Db.Get(s.Ctx, fmt.Sprintf("customrp.%d.%s", obj.Message.FromID, msg)).Result()
	if err != nil || action == "" {
		return
	}

	me := regexp.MustCompile(`(?i)(?:\A|)я(?:||\z)`)
	target := regexp.MustCompile(`(?i)(?:\A|)цель(?:||\z)`)

	action = me.ReplaceAllString(action, "[id"+
		strconv.Itoa(obj.Message.FromID)+
		"|"+
		core.GetNicknameOrFullName(obj.Message.FromID)+
		"] ")
	action = target.ReplaceAllString(action, "[id"+
		strconv.Itoa(id)+
		"|"+
		core.GetNicknameOrFullName(id)+
		"]")

	core.SendSimple(obj, action)
}
