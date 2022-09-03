package soyjack

import (
	"bytes"
	"io"
	"net/http"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const soyboy_file_path = "commands/soyjack/soyboy.png"
const gurba_file_path = "commands/soyjack/gurba.png"

const (
	mode_soyboy = "сойбой"
	mode_gurba  = "нс"
)

type mode_data struct {
	Mask []float64

	Wand *imagick.MagickWand

	Width  uint
	Height uint

	PosX int
	PosY int
}

var mwsoy, mwgurba *imagick.MagickWand

func Register() core.Command {
	mwsoy = imagick.NewMagickWand()
	mwsoy.ReadImage(soyboy_file_path)

	mwgurba = imagick.NewMagickWand()
	mwgurba.ReadImage(gurba_file_path)

	return core.Command{
		Aliases:     []string{"сой", "сойджек"},
		Description: "обмазать картинку соей",
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

	data := mode_data{}

	if mw1.GetImageWidth() < mw1.GetImageHeight() {
		data = mode_data{
			Mask: []float64{
				0, 0, 1, 20,
				1, 443, 33, 443,
				275, 443, 273, 423,
				275, 1, 238, 2,
			},
			Wand:   mwsoy,
			Width:  275,
			Height: 443,
			PosX:   25,
			PosY:   66,
		}
	} else {
		data = mode_data{
			Mask: []float64{
				0, 0, 3, 1,
				0, 414, 1, 412,
				595, 414, 595, 362,
				595, 0, 566, 6,
			},
			Wand:   mwgurba,
			Width:  595,
			Height: 414,
			PosX:   8,
			PosY:   1249,
		}
	}

	mw1.ResizeImage(data.Width, data.Height, imagick.FILTER_UNDEFINED, 1)
	mw1.SetImageVirtualPixelMethod(imagick.VIRTUAL_PIXEL_TRANSPARENT)
	mw1.DistortImage(imagick.DISTORTION_PERSPECTIVE, data.Mask, false)

	mw2 := data.Wand.Clone()
	mw2.CompositeLayers(mw1, imagick.COMPOSITE_OP_DST_OVER, data.PosX, data.PosY)

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
