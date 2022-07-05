package top

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"топ"},
		Description: "составить топ участников беседы",
		Handler:     handle,
	}
}

func handle(ctx *context.Context, obj *events.MessageNewObject) {
	name := strings.Join(core.ExtractArguments(obj), " ")

	if name == "" {
		core.ReplySimple(obj, "ошибка: не указано имя топа")

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

	rand.Shuffle(len(m.Profiles), func(i, j int) { m.Profiles[i], m.Profiles[j] = m.Profiles[j], m.Profiles[i] })

	n := rand.Intn(len(m.Profiles)) + 1

	if n == 0 {
		n = 1
	}

	txt := fmt.Sprintf("топ %d %s:\n", n, name)

	for i := 0; i < n; i++ {
		x := m.Profiles[i]

		name := core.GetAlias(x.ID)
		if name == "" {
			name = x.FirstName + " " + x.LastName
		}

		txt += "- [id" + strconv.Itoa(x.ID) + "|" + name + "]\n"
	}

	core.ReplySimple(obj, txt)
}
