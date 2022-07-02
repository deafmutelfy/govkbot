package nick

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"
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
	if utf8.RuneCountInString(nickname) > 128 {
		core.ReplySimple(obj, "ошибка: максимальная длина никнейма - 128 символов")

		return
	}

	s := core.GetStorage()

	s.Db.Set(s.Ctx, fmt.Sprintf("nicknames.%d", obj.Message.FromID), nickname, 0)

	core.ReplySimple(obj, "успешно")
}
