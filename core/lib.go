package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
)

func ReplySimple(obj *events.MessageNewObject, msg string, attachment ...interface{}) error {
	b := params.NewMessagesSendBuilder()
	fromId := obj.Message.FromID

	b.Message("[id" + strconv.Itoa(fromId) + "|" + GetNickname(fromId) + "], " + msg)
	b.DisableMentions(true)
	b.RandomID(0)
	b.PeerID(obj.Message.PeerID)

	if len(attachment) != 0 {
		b.Attachment(attachment[0])
	}

	_, err := GetStorage().Vk.MessagesSend(b.Params)
	return err
}

func GetNicknameWithoutSetup(userId int) string {
	s := GetStorage()
	key := fmt.Sprintf("nicknames.%d", userId)

	nickname, err := s.Db.Get(s.Ctx, key).Result()

	if err != nil {
		nickname = ""
	}

	return nickname
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

func ExtractAttachments(obj *events.MessageNewObject) []object.MessagesMessageAttachment {
	res := obj.Message.Attachments

	if obj.Message.ReplyMessage != nil {
		res = append(res, obj.Message.ReplyMessage.Attachments...)
	}

	for _, x := range obj.Message.FwdMessages {
		res = append(res, x.Attachments...)
	}

	return res
}
