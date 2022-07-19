package cm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"vkbot/core"
	"vkbot/core/rolesystem"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

func extractId(key string) int {
	l := strings.Split(key, ".")
	idint, _ := strconv.Atoi(l[len(l)-1])

	return idint
}

func listrole(_ *context.Context, obj *events.MessageNewObject) {
	if err := cmInit(obj); err != nil {
		core.ReplySimple(obj, err.Error())

		return
	}

	s := core.GetStorage()

	initialized, _ := s.Db.Get(s.Ctx, fmt.Sprintf("roles.%d.initialized", obj.Message.PeerID)).Result()
	if initialized != "true" {
		core.ReplySimple(obj, core.ERR_NO_ROLESYSTEM)

		return
	}

	keys, err := s.Db.Keys(s.Ctx, fmt.Sprintf("roles.%d.*", obj.Message.PeerID)).Result()
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}
	keys = core.Remove(keys, fmt.Sprintf("roles.%d.initialized", obj.Message.PeerID))

	values, err := s.Db.MGet(s.Ctx, keys...).Result()
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

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

	rlist := struct {
		owner  int
		admins []int
		moders []int
	}{}

	for i, x := range values {
		rint, _ := strconv.Atoi(x.(string))

		switch rint {
		case rolesystem.ROLE_OWNER:
			rlist.owner = extractId(keys[i])
		case rolesystem.ROLE_ADMINISTRATOR:
			rlist.admins = append(rlist.admins, extractId(keys[i]))
		case rolesystem.ROLE_MODERATOR:
			rlist.moders = append(rlist.moders, extractId(keys[i]))
		}
	}

	lookupName := func(id int) string {
		name := core.GetAlias(id)

		if name == "" {
			for _, x := range m.Profiles {
				if x.ID == id {
					name = x.FirstName + " " + x.LastName
				}
			}
		}

		return name
	}

	msg := "\nОснователь беседы: [id" + strconv.Itoa(rlist.owner) + "|" + lookupName(rlist.owner) + "]"

	if len(rlist.admins) != 0 {
		msg += "\n\nАдминистраторы:"

		for _, x := range rlist.admins {
			msg += "\n- [id" + strconv.Itoa(x) + "|" + lookupName(x) + "]"
		}
	}
	if len(rlist.moders) != 0 {
		msg += "\n\nМодераторы:"

		for _, x := range rlist.moders {
			msg += "\n- [id" + strconv.Itoa(x) + "|" + lookupName(x) + "]"
		}
	}

	core.ReplySimple(obj, msg)
}
