package linus

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const error_image_not_attached = "ошибка: нужно прикрепить картинку"
const linus_file_path = "commands/linus/linus.png"

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"линус"},
		Description: "перенести картинку на конференцию с Линусом Торвальдсом",
		Handler:     handle,
	}
}

func handle(ctx *context.Context, obj *events.MessageNewObject) {
	imagick.Initialize()
	defer imagick.Terminate()

	if len(obj.Message.Attachments) == 0 {
		core.ReplySimple(obj, error_image_not_attached)

		return
	}

	attachment := obj.Message.Attachments[0]
	if attachment.Type != "photo" {
		core.ReplySimple(obj, error_image_not_attached)

		return
	}

	response, err := http.Get(attachment.Photo.MaxSize().URL)

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	mask := []float64{
		0, 0, 6, 23,
		0, 275, 0, 275,
		435, 275, 435, 250,
		435, 0, 435, -74,
	}

	bt, err := io.ReadAll(response.Body)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	mw1 := imagick.NewMagickWand()
	mw1.ReadImageBlob(bt)
	mw1.ResizeImage(435, 275, imagick.FILTER_UNDEFINED, 1)
	mw1.SetImageVirtualPixelMethod(imagick.VIRTUAL_PIXEL_TRANSPARENT)
	mw1.DistortImage(imagick.DISTORTION_PERSPECTIVE, mask, false)

	mw2 := imagick.NewMagickWand()
	mw2.ReadImage(linus_file_path)
	mw2.CompositeLayers(mw1, imagick.COMPOSITE_OP_DST_OVER, 205, 0)

	s := core.GetStorage()

	vkPhoto, err := s.Vk.UploadMessagesPhoto(obj.Message.PeerID, bytes.NewReader(mw2.GetImageBlob()))

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	b := params.NewMessagesSendBuilder()
	fromId := obj.Message.FromID

	b.Message("[id" + strconv.Itoa(fromId) + "|" + core.GetNickname(fromId) + "], ваша картинка:")
	b.DisableMentions(true)
	b.RandomID(0)
	b.PeerID(obj.Message.PeerID)
	b.Attachment(vkPhoto)

	s.Vk.MessagesSend(b.Params)
}
