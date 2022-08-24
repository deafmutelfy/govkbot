package who

import (
	"errors"
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
		Aliases:     []string{"кто"},
		Description: "выбрать случайного участника беседы",
		Handler:     handle,
	}
}

func handle(obj *events.MessageNewObject) (err error) {
	desc := strings.Join(core.ExtractArguments(obj), " ")

	if desc == "" {
		core.ReplySimple(obj, "ошибка: не указано описание")

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

	u := m.Profiles[rand.Intn(len(m.Profiles))]

	name := core.GetAlias(u.ID)
	if name == "" {
		name = u.FirstName + " " + u.LastName
	}

	core.ReplySimple(obj, "кто "+desc+"? Возможно, это [id"+strconv.Itoa(u.ID)+"|"+name+"]")

	return
}
