package commemoration

import (
	"bytes"
	"io"
	"net/http"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const frame_file_path = "commands/commemoration/data.png"

var mw *imagick.MagickWand

func Register() core.Command {
	mw = imagick.NewMagickWand()
	mw.ReadImage(frame_file_path)

	return core.Command{
		Aliases:     []string{"поминки"},
		Description: "обернуть картинку похоронной рамкой",
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

	mw1 := imagick.NewMagickWand()
	mw1.ReadImageBlob(bt)
	mw1.ResizeImage(462, 584, imagick.FILTER_UNDEFINED, 1)
	mw1.TransformImageColorspace(imagick.COLORSPACE_GRAY)

	mw2 := mw.Clone()
	mw2.CompositeLayers(mw1, imagick.COMPOSITE_OP_DST_OVER, 496, 167)

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
