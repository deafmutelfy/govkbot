package core

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
)

func Send(obj *events.MessageNewObject, msg string, b *params.MessagesSendBuilder) (api.MessagesSendUserIDsResponse, error) {
	b.Message(msg)
	b.RandomID(0)
	b.PeerIDs([]int{obj.Message.PeerID})

	return GetStorage().Vk.MessagesSendPeerIDs(b.Params)
}

func ReplySimple(obj *events.MessageNewObject, msg string, attachment ...interface{}) error {
	fromId := obj.Message.FromID

	return SendSimple(obj, "[id"+strconv.Itoa(fromId)+"|"+GetNickname(fromId)+"], "+msg, attachment...)
}

func SendSimple(obj *events.MessageNewObject, msg string, attachment ...interface{}) error {
	b := params.NewMessagesSendBuilder()

	b.Message(msg)
	b.DisableMentions(true)
	b.RandomID(0)
	b.PeerID(obj.Message.PeerID)
	b.DontParseLinks(true)

	if len(attachment) != 0 {
		b.Attachment(attachment[0])
	}

	_, err := GetStorage().Vk.MessagesSend(b.Params)
	return err
}

func GetAlias(userId int) string {
	s := GetStorage()
	key := fmt.Sprintf("nicknames.%d", userId)

	nickname, err := s.Db.Get(s.Ctx, key).Result()

	if err != nil {
		nickname = ""
	}

	_, err = s.Db.Get(s.Ctx, key+".initialized").Result()

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

func GetNicknameOrFullName(userId int) string {
	name := GetAlias(userId)

	if name == "" {
		b := params.NewUsersGetBuilder()

		b.UserIDs([]string{strconv.Itoa(userId)})

		res, err := GetStorage().Vk.UsersGet(b.Params)
		if err != nil {
			name = "<без имени>"
		} else {
			name = res[0].FirstName + " " + res[0].LastName
		}
	}

	return name
}

func ExtractArguments(obj *events.MessageNewObject) []string {
	return strings.Split(obj.Message.Text, " ")[1:]
}

func ExtractAttachments(obj *events.MessageNewObject, t ...string) []object.MessagesMessageAttachment {
	res := obj.Message.Attachments

	if obj.Message.ReplyMessage != nil {
		res = append(res, obj.Message.ReplyMessage.Attachments...)
	}

	for _, x := range obj.Message.FwdMessages {
		res = append(res, x.Attachments...)
	}

	if len(t) > 0 {
		found := []object.MessagesMessageAttachment{}

		types := strings.Split(t[0], ",")
		for _, x := range res {
			if IsInArray(types, x.Type) {
				found = append(found, x)
			}
		}

		return found
	}

	return res
}

func GetMention(obj *events.MessageNewObject) int {
	if obj.Message.ReplyMessage != nil {
		return obj.Message.ReplyMessage.FromID
	}

	if len(obj.Message.FwdMessages) > 0 {
		return obj.Message.FwdMessages[0].FromID
	}

	r := regexp.MustCompile(`\[id(\d*)\|.*]`)

	res := r.FindStringSubmatch(obj.Message.Text)

	if len(res) == 0 {
		return 0
	}

	id, _ := strconv.Atoi(res[1])

	return id
}

func PeerIdToChatId(obj *events.MessageNewObject) int {
	return obj.Message.PeerID - 2000000000
}

func Remove[T comparable](l []T, item T) []T {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func IsInArray[T comparable](l []T, item T) bool {
	for _, x := range l {
		if x == item {
			return true
		}
	}

	return false
}
