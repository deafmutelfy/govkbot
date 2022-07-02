package nick

import (
	"context"
	"fmt"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"ник"},
		Description: "изменить свой никнейм",
		Handler:     handle,
	}
}

func handle(ctx *context.Context, obj *events.MessageNewObject) {
	nickname := strings.Join(core.ExtractArguments(obj), " ")

	if nickname == "" {
		core.ReplySimple(obj, "ошибка: не указан желаемый никнейм")

		return
	}
	nickname = strings.ReplaceAll(nickname, "\n", "")
	if len(nickname) > 32 {
		core.ReplySimple(obj, "ошибка: максимальная длина никнейма - 32 символа")

		return
	}

	s := core.GetStorage()

	s.Db.Set(s.Ctx, fmt.Sprintf("nicknames.%d", obj.Message.FromID), nickname, 0)

	core.ReplySimple(obj, "успешно")
}
