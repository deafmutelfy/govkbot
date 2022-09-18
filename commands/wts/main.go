package wts

import (
	"math/rand"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"втс"},
		Description: "пробелы в смайлики",
		Handler:     handle,
	}
}

var smiles = []string{
	"😖",
	"😣",
	"🤧",
	"😥",
	"🔪",
	"🙏🏻",
	"💞",
	"😭",
}

func handle(obj *events.MessageNewObject) (err error) {
	obj.Message.Text = strings.ToLower(obj.Message.Text)

	args := core.ExtractArguments(obj)
	if len(args) < 1 {
		core.ReplySimple(obj, "ошибка: недостаточно параметров")

		return
	}

	res := ""
	first := true
	for _, x := range args {
		if first {
			first = false
			res += x

			continue
		}

		res += smiles[rand.Intn(len(smiles))] + x
	}

	core.SendSimple(obj, res)

	return
}
