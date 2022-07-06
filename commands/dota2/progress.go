package dota2

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/jasonodonnell/go-opendota"
)

func progress(ctx *context.Context, obj *events.MessageNewObject) {
	args := core.ExtractArguments(obj)
	if len(args) < 2 {
		core.ReplySimple(obj, "ошибка: необходимо указать количество дней")

		return
	}

	days, err := strconv.Atoi(args[1])
	if err != nil || days < 1 {
		core.ReplySimple(obj, "ошибка: недопустимое значение количества дней")

		return
	}
	if days > 7 {
		core.ReplySimple(obj, "ошибка: значение количества дней не должно превышать 7")

		return
	}

	s := core.GetStorage()

	id, err := s.Db.Get(s.Ctx, fmt.Sprintf("dota2.%d.id", obj.Message.FromID)).Result()
	if id == "" || err != nil {
		core.ReplySimple(obj, ERR_NO_ID)

		return
	}
	numId, _ := strconv.ParseInt(id, 10, 64)

	client := opendota.NewClient(&http.Client{})

	m, _, err := client.PlayerService.Matches(numId, &opendota.PlayerParam{Date: days})
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	var w, l int

	for _, x := range m {
		if x.RadiantWin {
			w++
		} else {
			l++
		}
	}

	mmrdiff := 0 + (30 * w) - (30 * l)

	msg := "\nМатчей сыграно: " + strconv.Itoa(w+l)
	msg += "\nПобед: " + strconv.Itoa(w)
	msg += "\nПоражений: " + strconv.Itoa(l)

	msg += "\n\nВаш рейтинг "
	if mmrdiff != 0 {
		if mmrdiff > 0 {
			msg += "повысился"
		} else {
			msg += "понизился"
			mmrdiff *= -1
		}

		msg += " на " + strconv.Itoa(mmrdiff) + " MMR"
	} else {
		msg += "не изменился "
	}

	core.ReplySimple(obj, msg)
}
