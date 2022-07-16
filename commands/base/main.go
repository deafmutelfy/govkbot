package base

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/valyala/fastjson"
)

const bazman_uri = "https://bazman.ctw.re/"

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"база"},
		Description: "сгенерировать мудрую, умную мысль",
		Handler:     handle,
	}
}

func list(obj *events.MessageNewObject) {
	r, err := http.Get(bazman_uri + "bases")
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	p, err := fastjson.ParseBytes(b)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	msg := "список источников:"

	for k, x := range p.GetArray() {
		msg += "\n" + strconv.Itoa(k) + ": " + string(x.GetStringBytes("name"))
	}

	core.ReplySimple(obj, msg)
}

func handle(ctx *context.Context, obj *events.MessageNewObject) {
	args := core.ExtractArguments(obj)
	if len(args) < 1 {
		core.ReplySimple(obj, "ошибка: не указан источник. Список доступных источников можно получить командой \"/база лист\"")

		return
	}
	if args[0] == "лист" {
		list(obj)

		return
	}

	_, err := strconv.Atoi(args[0])
	if err != nil {
		core.ReplySimple(obj, "ошибка: недопустимый идентификатор источника")

		return
	}

	q := url.Values{}
	q.Set("num", args[0])

	u := &url.URL{
		Scheme:   "https",
		Host:     "bazman.ctw.re",
		Path:     "gen",
		RawQuery: q.Encode(),
	}

	r, err := http.Get(u.String())
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}
	if r.StatusCode == 500 {
		core.ReplySimple(obj, "ошибка: такого источника не существует")

		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	bl := params.NewMessagesSendBuilder()

	bl.Message(string(b))
	bl.DisableMentions(true)
	bl.RandomID(0)
	bl.PeerID(obj.Message.PeerID)
	bl.DontParseLinks(true)

	core.GetStorage().Vk.MessagesSend(bl.Params)
}
