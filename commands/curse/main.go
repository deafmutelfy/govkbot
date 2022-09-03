package curse

import (
	"bytes"
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

func handle(obj *events.MessageNewObject) (err error) {
	atts := core.ExtractAttachments(obj, "photo")

	if len(atts) == 0 {
		core.ReplySimple(obj, core.ERR_NO_PICTURE)

		return
	}

	attachment := atts[0]

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

	vkPhoto, err := core.GetStorage().Vk.UploadMessagesPhoto(0, bytes.NewReader(mw.GetImageBlob()))

	mw.Destroy()

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "ваша картинка:", vkPhoto)

	return
}
