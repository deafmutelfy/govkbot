package wts

import (
	"math/rand"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/dlclark/regexp2"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"втс"},
		Description: "пробелы в смайлики",
		Handler:     handle,
	}
}

var r = regexp2.MustCompile(`[^\p{L}\p{N} ]+`, 0)

var smiles = []string{
	"😖",
	"😣",
	"🤧",
	"😥",
	"🔪",
	"🙏🏻",
	"💞",
	"😭",
	"💜",
	"🚀",
	"👀",
	"💥",
	"💔",
}

func handle(obj *events.MessageNewObject) (err error) {
	txt, err := r.Replace(obj.Message.Text, "", 0, -1)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)
	}

	obj.Message.Text = strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(txt)), " "))

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
