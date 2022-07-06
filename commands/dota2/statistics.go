package dota2

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/jasonodonnell/go-opendota"
)

func statisticsGetRank(p *opendota.Player) string {
	if p.RankTier == 0 {
		return "не откалиброван"
	}

	rank := []string{"Рекрут", "Страж", "Рыцарь", "Герой", "Легенда", "Властелин", "Божество", "Титан"}[int(p.RankTier/10)-1]
	star := strconv.Itoa(p.RankTier % 10)

	if star != "0" {
		rank += " " + star
	}

	return rank
}

func statisticsGenMsg(accId int64) string {
	msg := ""
	client := opendota.NewClient(&http.Client{})

	p, _, err := client.PlayerService.Player(accId)
	if err != nil {
		msg = core.ERR_UNKNOWN

		return msg
	}
	if p.Profile.AccountID == 0 {
		msg = "ошибка: информации по вашему ID не найдено. Установите идентификатор заново, или же, если ID верный, включите общедоступную историю матчей в настройках игры"

		return msg
	}

	msg += "\nНикнейм: " + p.Profile.Personaname
	msg += "\nРанг: " + statisticsGetRank(&p)

	wl, _, err := client.PlayerService.WinLoss(accId, nil)
	if err != nil {
		msg = core.ERR_UNKNOWN

		return msg
	}

	log.Println(float64(100) / float64(wl.Lose+wl.Win))

	msg += "\n\nПобед: " + strconv.Itoa(wl.Win)
	msg += "\nПоражений: " + strconv.Itoa(wl.Lose)
	msg += "\nВинрейт: " + fmt.Sprintf("%.2f", float64(100)/float64(wl.Lose+wl.Win)*float64(wl.Win)) + "%"

	return msg
}

func statistics(ctx *context.Context, obj *events.MessageNewObject) {
	s := core.GetStorage()

	id, err := s.Db.Get(s.Ctx, fmt.Sprintf("dota2.%d.id", obj.Message.FromID)).Result()

	if id == "" || err != nil {
		core.ReplySimple(obj, "ошибка: ваш ID не установлен. Для установки ID воспользуйтесь командой \"/д2 айди <ваш ID>\"")

		return
	}

	numId, _ := strconv.ParseInt(id, 10, 64)

	core.ReplySimple(obj, statisticsGenMsg(numId))
}
