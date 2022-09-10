package linus

import (
	"bytes"
	"io"
	"net/http"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const linus_file_path = "commands/linus/linus.png"

var mw *imagick.MagickWand
var cutset *imagick.MagickWand

func Register() core.Command {
	mw = imagick.NewMagickWand()
	mw.ReadImage(linus_file_path)

	cutset = mw.Clone()
	cutset.CropImage(435, 275, 205, 0)

	return core.Command{
		Aliases:     []string{"линус"},
		Description: "перенести картинку на конференцию с Линусом Торвальдсом",
		Handler:     handle,
	}
}

func handle(obj *events.MessageNewObject) (err error) {
	// atts := core.ExtractAttachments(obj, "photo,doc")
	atts := core.ExtractAttachments(obj, "photo")

	if len(atts) == 0 {
		core.ReplySimple(obj, core.ERR_NO_PICTURE)

		return
	}

	attachment := atts[0]

	link := ""

	switch attachment.Type {
	case "photo":
		link = attachment.Photo.MaxSize().URL
	case "doc":
		link = attachment.Doc.URL

		if attachment.Doc.Size > 30*1024*1024 {
			core.ReplySimple(obj, core.ERR_LARGE_GIF)
		}
	}

	response, err := http.Get(link)

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

	mw2 := mw.Clone()
	mw2.CompositeLayers(mw1, imagick.COMPOSITE_OP_DST_OVER, 205, 0)

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
