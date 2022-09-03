package tacticalpic

// легаси код ай донт рили вонт ту си зэт

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"боевая", "бой"},
		Description: "сделать картинку боевой",
		Handler:     handle,
	}
}

func list(obj *events.MessageNewObject) {
	core.ReplySimple(obj, "возможные расположения диалогового облака:\nсправа\nцентр\nслева")
}

const bubble_height_default = 260

func handle(obj *events.MessageNewObject) (err error) {
	args := core.ExtractArguments(obj)
	if len(args) > 0 {
		if args[0] == "лист" {
			list(obj)

			return
		}
	}

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

	var idx int
	if len(args) == 0 {
		idx = rand.Intn(2)
	} else {
		switch args[0] {
		case "справа":
			idx = 0
		case "слева":
			idx = 1
		case "центр":
			idx = 2
		default:
			list(obj)

			return
		}
	}

	mw2 := imagick.NewMagickWand()
	mw2.ReadImage(fmt.Sprintf("commands/tacticalpic/speech-bubble%d.png", idx))
	width := mw1.GetImageWidth()
	ratio := float32(width) / float32(mw2.GetImageWidth())
	height := uint(float32(mw2.GetImageHeight()) * ratio)
	mw2.AdaptiveResizeImage(width, height)

	bubble_height := bubble_height_default * ratio
	mw1.ExtentImage(width, uint(float32(mw1.GetImageHeight())+bubble_height), 0, int(bubble_height)*-1)
	mw1.CompositeImage(mw2, imagick.COMPOSITE_OP_OVER, 0, 0)

	vkPhoto, err := core.GetStorage().Vk.UploadMessagesPhoto(0, bytes.NewReader(mw1.GetImageBlob()))

	mw1.Destroy()
	mw2.Destroy()

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "ваша картинка:", vkPhoto)

	return
}
