package linus

import (
	"bytes"
	"io"
	"net/http"
	"vkbot/core"
	"vkbot/subsystems/queuesystem"

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
	atts := core.ExtractAttachments(obj, "photo,doc")

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

	if attachment.Type == "doc" {
		queuesystem.Add(obj, func(obj *events.MessageNewObject) (err error) {
			aw := mw1.CoalesceImages()

			mw2 := mw.Clone()

			for i := 0; i < int(aw.GetNumberImages()); i++ {
				aw.SetIteratorIndex(i)
				img := aw.GetImage()

				img.ResizeImage(435, 275, imagick.FILTER_UNDEFINED, 1)
				img.SetImageVirtualPixelMethod(imagick.VIRTUAL_PIXEL_TRANSPARENT)
				img.DistortImage(imagick.DISTORTION_PERSPECTIVE, mask, false)
				img.SetImagePage(0, 0, 205, 0)

				img.CompositeLayers(cutset, imagick.COMPOSITE_OP_SRC_OVER, 0, 0)

				mw2.AddImage(img)
				img.Destroy()
			}

			mw2.SetFormat("gif")
			vkPhoto, err := core.GetStorage().Vk.UploadMessagesDoc(obj.Message.PeerID, "doc", "deafmute-bot.gif", "", bytes.NewReader(mw2.GetImagesBlob()))

			mw1.Destroy()
			mw2.Destroy()
			aw.Destroy()

			if err != nil {
				core.ReplySimple(obj, core.ERR_UNKNOWN)

				return
			}

			core.ReplySimple(obj, "ваша картинка:", vkPhoto.Doc)

			return
		})

		return
	}

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
