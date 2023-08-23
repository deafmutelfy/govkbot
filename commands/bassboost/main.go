package bassboost

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"vkbot/core"
	"vkbot/subsystems/audiosystem"

	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"бассбуст", "пушечка"},
		Description: "инструмент для бассбуста",
		Handler:     handle,
		Queue: &core.CommandQueueParams{
			Name: "audio",
		},
	}
}

func handle(obj *events.MessageNewObject) (err error) {
	args := core.ExtractArguments(obj)

	bass := "10"
	treble := "0"

	if len(args) == 2 {
		bass = args[0]
		treble = args[1]

		if _, err = strconv.Atoi(bass); err != nil {
			core.ReplySimple(obj, "ошибка: недопустимое значение нижней частоты")

			return
		}

		if _, err = strconv.Atoi(treble); err != nil {
			core.ReplySimple(obj, "ошибка: недопустимое значение верхней частоты")

			return
		}
	}

	atts := core.ExtractAttachments(obj, "audio")

	if len(atts) == 0 {
		core.ReplySimple(obj, core.ERR_NO_AUDIO)

		return
	}

	attachment := atts[0]

	if attachment.Audio.Duration > 15*60 {
		core.ReplySimple(obj, core.ERR_10MIN)

		return
	}
	if attachment.Audio.URL == "" {
		core.ReplySimple(obj, core.ERR_AUDIO_VK_API_BUG)

		return
	}

	response, err := http.Get(attachment.Audio.URL)

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	buffbass := bytes.Buffer{}
	cmdbass := exec.Command("ffmpeg", "-i", "-", "-af", "bass=g="+bass, "-f", "mp3", "-")

	cmdbass.Stdin = response.Body
	cmdbass.Stdout = &buffbass

	if cmdbass.Run() != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	bufftreble := bytes.Buffer{}
	cmdtreble := exec.Command("ffmpeg", "-i", "-", "-af", "treble=g="+treble, "-f", "mp3", "-")

	cmdtreble.Stdin = &buffbass
	cmdtreble.Stdout = &bufftreble

	if cmdtreble.Run() != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	d, err := audiosystem.UploadAudio(&bufftreble,
		attachment.Audio.Artist,
		fmt.Sprintf("%s (bassboosted by deafmute bot, bass=%sdB treble=%sdB)", attachment.Audio.Title, bass, treble))
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
