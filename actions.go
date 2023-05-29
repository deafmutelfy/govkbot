package main

import (
	"fmt"
	"strconv"
	"strings"
	"vkbot/commands/cm"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/dlclark/regexp2"
)

func handleChatInviteUser(obj *events.MessageNewObject) {
	s := core.GetStorage()

	if obj.Message.Action.MemberID == -s.Cfg.GroupId {
		return
	}

	ok, _ := s.Db.Get(s.Ctx, fmt.Sprintf("bans.%d.%d", obj.Message.PeerID, obj.Message.Action.MemberID)).Result()
	if ok == "true" && obj.Message.Action.MemberID > 0 {
		b := params.NewMessagesRemoveChatUserBuilder()

		b.MemberID(obj.Message.Action.MemberID)
		b.ChatID(core.PeerIdToChatId(obj))

		s := core.GetStorage()

		_, err := s.Vk.MessagesRemoveChatUser(b.Params)
		if err != nil {
			core.SendSimple(obj, "Возникла ошибка при исключении [id"+strconv.Itoa(obj.Message.Action.MemberID)+"|заблокированного пользователя] из беседы. Для уточнения причины воспользуйтесь командой \"/чм кик\"")

			return
		}

		core.SendSimple(obj, "[id"+strconv.Itoa(obj.Message.Action.MemberID)+"|Заблокированный пользователь] исключён")

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

	action, err = regexp2.MustCompile(`(?i)\bя\b`, 0).Replace(action, "[id"+
		strconv.Itoa(obj.Message.FromID)+
		"|"+
		n1+
		"]", 0, -1)
	if err != nil {
		return
	}
	action, err = regexp2.MustCompile(`(?i)\bцель\b`, 0).Replace(action, "[id"+
		strconv.Itoa(id)+
		"|"+
		n2+
		"]", 0, -1)
	if err != nil {
		return
	}

	core.SendSimple(obj, "* "+action)
}
