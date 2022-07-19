package cm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"vkbot/core"
	"vkbot/core/rolesystem"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

func initrole(_ *context.Context, obj *events.MessageNewObject) {
	if err := cmInit(obj); err != nil {
		core.ReplySimple(obj, err.Error())

		return
	}

	b := params.NewMessagesGetConversationMembersBuilder()

	b.PeerID(obj.Message.PeerID)

	m, err := core.GetStorage().Vk.MessagesGetConversationMembers(b.Params)

	if errors.Is(err, api.ErrMessagesChatUserNoAccess) {
		core.ReplySimple(obj, core.ERR_NO_ACCESS_TO_CHAT)

		return
	} else if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	s := core.GetStorage()

	initializedKey := fmt.Sprintf("roles.%d.initialized", obj.Message.PeerID)
	initialized, _ := s.Db.Get(s.Ctx, initializedKey).Result()

	if initialized == "true" {
		core.ReplySimple(obj, "ошибка: система ролей для этой беседы уже инициализирована")

		return
	}

	txt := "успешно:\n"

	for _, x := range m.Items {
		if (!x.IsOwner && !x.IsAdmin) || x.MemberID < 0 {
			continue
		}

		r := fmt.Sprintf("roles.%d.%d", obj.Message.PeerID, x.MemberID)

		role := ""

		if x.IsOwner {
			s.Db.Set(s.Ctx, r, rolesystem.ROLE_OWNER, 0)

			txt += "Основателю беседы "
			role = "Основатель"
		} else if x.IsAdmin {
			s.Db.Set(s.Ctx, r, rolesystem.ROLE_MODERATOR, 0)

			txt += "Администратору беседы "
			role = "Модератор"
		}

		name := core.GetAlias(x.MemberID)
		if name == "" {
			for _, v := range m.Profiles {
				if v.ID == x.MemberID {
					name = v.FirstName + " " + v.LastName
				}
			}
		}

		txt += "[id" + strconv.Itoa(x.MemberID) + "|" + name + "] выдана роль \"" + role + "\"\n"
	}

	s.Db.Set(s.Ctx, initializedKey, "true", 0)

	core.ReplySimple(obj, txt)
}
