package rule34

import (
	"context"
	"math/rand"
	"net/http"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	gonnachan "github.com/insomnyawolf/Gonnachan"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"р34", "r34"},
		Description: "найти материал в Rule 34",
		Handler:     handle,
	}
}

func handle(ctx *context.Context, obj *events.MessageNewObject) {
	tags := core.ExtractArguments(obj)

	if len(tags) == 0 {
		core.ReplySimple(obj, "ошибка: теги не указаны")
	}

	req := gonnachan.PostRequest{Tags: tags, MaxResults: 64, TargetAPI: gonnachan.ServerRule34}

	r, err := req.GetResults()

	l := len(r)

	if l == 0 || err != nil {
		core.ReplySimple(obj, "ошибка: ничего не найдено")

		return
	}

	post := r[rand.Intn(l)]

	pic, err := http.Get(post.FileURL)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	vkPhoto, err := core.GetStorage().Vk.UploadMessagesPhoto(obj.Message.PeerID, pic.Body)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "ваша картинка:", vkPhoto)
}
