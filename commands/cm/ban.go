package cm

import (
	"errors"
	"strconv"
	"vkbot/core"
	"vkbot/subsystems/rolesystem"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
)

func ban(obj *events.MessageNewObject) (err error) {
	if err = cmInit(obj); err != nil {
		core.ReplySimple(obj, err.Error())

		return
	}

	id := core.GetMention(obj)
	if id == 0 {
		core.ReplySimple(obj, core.ERR_NO_TARGET)

		return
	}

	role1 := rolesystem.GetRole(obj)
	role2 := rolesystem.GetRole(&events.MessageNewObject{Message: object.MessagesMessage{FromID: id}})

	if role1 <= role2 && (obj.Message.FromID != id) {
		core.ReplySimple(obj, core.ERR_ACCESS_DENIED)

		return
	}

	b := params.NewMessagesRemoveChatUserBuilder()

	b.MemberID(id)
	b.ChatID(core.PeerIdToChatId(obj))

	s := core.GetStorage()

	_, err = s.Vk.MessagesRemoveChatUser(b.Params)

	msg := ""
	if errors.Is(err, api.ErrAccess) {
		msg = "возникла ошибка с исключением пользователя из беседы, но он внесён в список заблокированных"
	} else if errors.Is(err, api.ErrMessagesChatUserNotInChat) {
		msg = "пользователя нет в беседе, но он внесён в список заблокированных"
	} else if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	s.Db.Set(s.Ctx, "bans."+strconv.Itoa(obj.Message.PeerID)+"."+strconv.Itoa(id), "true", 0)

	if msg == "" {
		msg = "успешно"
	} else {
		msg = "успешно: " + msg
	}

	core.ReplySimple(obj, msg)

	return
}
