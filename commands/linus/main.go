package linus

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

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

	atts := core.ExtractAttachments(obj)

	if len(atts) == 0 {
		core.ReplySimple(obj, core.ERR_NO_PICTURE)

		return
	}

	attachment := atts[0]
	if attachment.Type != "photo" {
		core.ReplySimple(obj, core.ERR_NO_PICTURE)

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

	vkPhoto, err := core.GetStorage().Vk.UploadMessagesPhoto(obj.Message.PeerID, bytes.NewReader(mw2.GetImageBlob()))

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "ваша картинка:", vkPhoto)
}
