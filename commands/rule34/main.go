package rule34

import (
	"context"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/valyala/fastjson"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"р34", "r34"},
		Description: "найти материал в Rule 34",
		Handler:     handle,
	}
}

func handle(ctx *context.Context, obj *events.MessageNewObject) {
	tags := strings.Join(strings.Split(obj.Message.Text, " ")[1:], "+")

	if tags == "" {
		core.ReplySimple(obj, "ошибка: теги не указаны")
	}

	q := url.Values{}
	q.Set("tags", tags)
	q.Set("limit", "200")

	u := &url.URL{
		Scheme:   "https",
		Host:     "r34-json-api.herokuapp.com",
		Path:     "posts",
		RawQuery: q.Encode(),
	}

	response, err := http.Get(u.String())
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	bt, err := io.ReadAll(response.Body)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	v, err := fastjson.ParseBytes(bt)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	l := len(v.GetArray())

	if l == 0 {
		core.ReplySimple(obj, "ошибка: ничего не найдено")

		return
	}

	post := v.GetArray()[rand.Intn(l)]

	pic, err := http.Get(string(post.GetStringBytes("file_url")))
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
