package online

import (
	"context"
	"errors"
	"strconv"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"онлайн"},
		Description: "получить список участников беседы, находящихся в сети",
		Handler:     handle,
	}
}

func handle(ctx *context.Context, obj *events.MessageNewObject) {
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

	txt := "список участников, находящихся в сети:\n"

	for _, x := range m.Profiles {
		txt += "[id" + strconv.Itoa(x.ID) + "|" + x.FirstName + " " + x.LastName + "]\n"
	}

	core.ReplySimple(obj, txt)
}
