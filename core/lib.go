package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

func ReplySimple(obj *events.MessageNewObject, msg string) error {
	b := params.NewMessagesSendBuilder()
	fromId := obj.Message.FromID

	b.Message("[id" + strconv.Itoa(fromId) + "|" + GetNickname(fromId) + "], " + msg)
	b.DisableMentions(true)
	b.RandomID(0)
	b.PeerID(obj.Message.PeerID)

	_, err := GetStorage().Vk.MessagesSend(b.Params)
	return err
}

func GetNickname(userId int) string {
	s := GetStorage()
	key := fmt.Sprintf("nicknames.%d", userId)

	nickname, err := s.Db.Get(s.Ctx, key).Result()

	if nickname == "" || err != nil {
		b := params.NewUsersGetBuilder()

		b.UserIDs([]string{strconv.Itoa(userId)})

		u, err := s.Vk.UsersGet(b.Params)

		if err != nil {
			nickname = "<без имени>"
		} else {
			nickname = u[0].FirstName

			s.Db.Set(s.Ctx, key, nickname, 0)
		}
	}

	return nickname
}

func ExtractArguments(obj *events.MessageNewObject) []string {
	return strings.Split(obj.Message.Text, " ")[1:]
}
