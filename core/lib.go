package core

import "github.com/SevereCloud/vksdk/v2/api/params"

func SendSimple(peerId int, msg string) error {
	b := params.NewMessagesSendBuilder()

	b.Message(msg)
	b.RandomID(0)
	b.PeerID(peerId)

	_, err := GetStorage().Vk.MessagesSend(b.Params)
	return err
}
