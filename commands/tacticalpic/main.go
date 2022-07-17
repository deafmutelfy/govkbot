package tacticalpic

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

var speech_bubbles = [...]string{
	"center",
	"right",
	"left",
}

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"боевая", "бой"},
		Description: "сделать картинку боевой",
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

	mw1 := imagick.NewMagickWand()
	mw1.ReadImageBlob(bt)

	args := core.ExtractArguments(obj)
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
			core.ReplySimple(obj, "типы боевых картинок: справа/слева/центр")
			return
		}
	}

	mw2 := imagick.NewMagickWand()
	mw2.ReadImage(fmt.Sprintf("commands/tacticalpic/speech-bubble%d.png", idx))
	width := mw1.GetImageWidth()
	height := uint(float32(mw2.GetImageHeight()) * (float32(width) / float32(mw2.GetImageWidth())))
	mw2.AdaptiveResizeImage(width, height)
	mw1.CompositeImage(mw2, imagick.COMPOSITE_OP_OVER, 0, 0)

	vkPhoto, err := core.GetStorage().Vk.UploadMessagesPhoto(obj.Message.PeerID, bytes.NewReader(mw1.GetImageBlob()))

	mw1.Destroy()
	mw2.Destroy()

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "ваша картинка:", vkPhoto)
}
