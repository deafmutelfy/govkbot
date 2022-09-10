package wrong

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const data_file_path = "commands/wrong/data.png"

var mw *imagick.MagickWand

func Register() core.Command {
	mw = imagick.NewMagickWand()
	mw.ReadImage(data_file_path)

	return core.Command{
		Aliases:     []string{"неправ"},
		Description: "в чём я не прав?",
		Handler:     handle,
	}
}

func handle(obj *events.MessageNewObject) (err error) {
	atts := core.ExtractAttachments(obj, "photo")

	url := ""
	if len(atts) == 0 {
		id := core.GetMention(obj)
		if id == 0 {
			core.ReplySimple(obj, core.ERR_NO_PICTURE)

			return
		}

		b := params.NewUsersGetBuilder()

		b.UserIDs([]string{strconv.Itoa(id)})
		b.Fields([]string{"photo_400_orig"})

		r, err := core.GetStorage().Vk.UsersGet(b.Params)
		if err != nil {
			core.ReplySimple(obj, core.ERR_UNKNOWN)

			return err
		}

		url = r[0].Photo400Orig
	} else {
		url = atts[0].Photo.MaxSize().URL
	}

	response, err := http.Get(url)
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
	mw1.ResizeImage(651, 623, imagick.FILTER_UNDEFINED, 1)

	mw2 := mw.Clone()
	mw2.CompositeLayers(mw1, imagick.COMPOSITE_OP_DST_OVER, 188, 136)

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
