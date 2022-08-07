package dota2

import (
	"fmt"
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

	msg += "\n\nПобед: " + strconv.Itoa(wl.Win)
	msg += "\nПоражений: " + strconv.Itoa(wl.Lose)
	msg += "\nВинрейт: " + fmt.Sprintf("%.2f", float64(100)/float64(wl.Lose+wl.Win)*float64(wl.Win)) + "%"

	h, _, err := client.PlayerService.Heroes(accId, nil)
	if err != nil {
		msg = core.ERR_UNKNOWN

		return msg
	}
	hId, _ := strconv.Atoi(h[0].HeroID)

	msg += "\n\nСигнатурный герой: " + heroes[hId-1]
	msg += "\nВинрейт на сигнатурном герое: " + fmt.Sprintf("%.2f", float64(100)/float64(h[0].Games)*float64(h[0].Win)) + "%"

	rm, _, err := client.PlayerService.RecentMatches(accId)
	if err != nil {
		msg = core.ERR_UNKNOWN

		return msg
	}

	msg += "\n\nПоследний матч:\n"
	if rm[0].RadiantWin {
		msg += "Победа"
	} else {
		msg += "Поражение"
	}
	msg += " - " + heroes[rm[0].HeroID-1]
	msg += " - " + fmt.Sprintf("%d/%d/%d", rm[0].Kills, rm[0].Deaths, rm[0].Assists)

	return msg
}

func statistics(obj *events.MessageNewObject) {
	s := core.GetStorage()

	id, err := s.Db.Get(s.Ctx, fmt.Sprintf("dota2.%d.id", obj.Message.FromID)).Result()
	if id == "" || err != nil {
		core.ReplySimple(obj, ERR_NO_ID)

		return
	}

	numId, _ := strconv.ParseInt(id, 10, 64)

	core.ReplySimple(obj, statisticsGenMsg(numId))
}
