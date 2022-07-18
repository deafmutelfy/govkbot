package curse

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"перекос"},
		Description: "перекосить картинку",
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

	bt, err := io.ReadAll(response.Body)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	mw := imagick.NewMagickWand()
	mw.ReadImageBlob(bt)
	mw.LiquidRescaleImage(mw.GetImageWidth()/2, mw.GetImageHeight()/2, 1, 0)
	mw.ResizeImage(mw.GetImageWidth()*2, mw.GetImageHeight()*2, imagick.FILTER_UNDEFINED, 1)

	mw.Destroy()

	vkPhoto, err := core.GetStorage().Vk.UploadMessagesPhoto(obj.Message.PeerID, bytes.NewReader(mw.GetImageBlob()))

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "ваша картинка:", vkPhoto)
}
