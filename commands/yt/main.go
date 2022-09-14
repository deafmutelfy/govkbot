package yt

import (
	"bytes"
	"os/exec"
	"strconv"
	"vkbot/core"
	"vkbot/subsystems/audiosystem"

	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/kkdai/youtube/v2"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"добавить"},
		Description: "перекачать видео с YouTube в аудиозаписи ВКонтакте",
		Handler:     handle,
		Queue: &core.CommandQueueParams{
			Name: "audio",
		},
	}
}

func handle(obj *events.MessageNewObject) (err error) {
	args := core.ExtractArguments(obj)
	if len(args) < 1 {
		core.ReplySimple(obj, "ошибка: необходимо указать ссылку на YouTube-ролик")

		return
	}

	cl := youtube.Client{}

	v, err := cl.GetVideo(args[0])
	if err != nil {
		core.ReplySimple(obj, "ошибка: неверная ссылка")

		return
	}

	f := v.Formats.FindByItag(140)

	res, _, err := cl.GetStream(v, f)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	buffconvert := bytes.Buffer{}
	cmdbass := exec.Command("ffmpeg", "-i", "-", "-b:a", "192K", "-vn", "-f", "mp3", "-")

	cmdbass.Stdin = res
	cmdbass.Stdout = &buffconvert

	if cmdbass.Run() != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	d, err := audiosystem.UploadAudio(&buffconvert, "deafmute bot", v.Title)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	r := (*d).(map[string]interface{})

	core.ReplySimple(obj, "ваша аудиозапись:",
		"audio"+
			strconv.FormatInt(int64(r["owner_id"].(float64)), 10)+
			"_"+
			strconv.FormatInt(int64(r["id"].(float64)), 10),
	)

	return
}
