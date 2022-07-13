package kick

import (
	"context"
	"errors"
	"vkbot/core"
	"vkbot/core/rolesystem"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"кик"},
		Description: "исключить участника беседы",
		Handler:     handle,
	}
}

func handle(ctx *context.Context, obj *events.MessageNewObject) {
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

	_, err := s.Vk.MessagesRemoveChatUser(b.Params)

	if errors.Is(err, api.ErrAccess) {
		core.ReplySimple(obj, core.ERR_NO_ACCESS_TO_CHAT+", или этого пользователя невозможно исключить")

		return
	} else if errors.Is(err, api.ErrMessagesChatUserNotInChat) {
		core.ReplySimple(obj, "ошибка: пользователь не найден")

		return
	} else if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "успешно")
}
