package help

import (
	"context"
	"strings"
	"sync"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

const repo_url = "https://github.com/deafmutelfy/govkbot"

var msg string = "список команд:\n"
var once sync.Once

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"хелп", "помощь"},
		Description: "получить справку о командах бота",
		Handler:     handle,
	}
}

func handle(ctx *context.Context, obj *events.MessageNewObject) {
	once.Do(func() {
		for _, x := range *core.GetStorage().CommandPool {
			if x.Hidden {
				continue
			}

			aliases := []string{}

			for _, v := range x.Aliases {
				aliases = append(aliases, "/"+v)
			}

			msg += strings.Join(aliases, ", ") + " - " + x.Description + "\n"
		}

		msg += "\nИсходный код: " + repo_url
	})

	core.ReplySimple(obj, msg)
}
