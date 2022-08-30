package window

import (
	"bytes"
	"io"
	"net/http"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const window_file_path = "commands/window/win.png"

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"окно"},
		Description: "перенести картинку на шаблон с окном",
		Handler:     handle,
	}
}

func handle(obj *events.MessageNewObject) (err error) {
	imagick.Initialize()
	defer imagick.Terminate()

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

	mw1 := imagick.NewMagickWand()
	mw1.ReadImageBlob(bt)
	mw1.ResizeImage(701, 397, imagick.FILTER_UNDEFINED, 1)

	mw2 := imagick.NewMagickWand()
	mw2.ReadImage(window_file_path)
	mw2.CompositeLayers(mw1, imagick.COMPOSITE_OP_DST_OVER, 789, 153)

	vkPhoto, err := core.GetStorage().Vk.UploadMessagesPhoto(0, bytes.NewReader(mw2.GetImageBlob()))

	mw1.Destroy()
	mw2.Destroy()

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "ваша картинка:", vkPhoto)

	return
}
