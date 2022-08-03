package slowreverb

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"vkbot/commands/bassboost"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"зр", "слоуреверб"},
		Description: "наложить замедление и реверберацию на аудиозапись",
		Handler:     handle,
		QueueName:   "audio",
	}
}

func handle(_ *context.Context, obj *events.MessageNewObject) {
	args := core.ExtractArguments(obj)

	ratio := "0.75"

	if len(args) == 1 {
		ratio = args[0]

		if _, err := strconv.ParseFloat(ratio, 32); err != nil {
			core.ReplySimple(obj, "ошибка: недопустимое значение коэффицента замедления")

			return
		}
	}

	atts := core.ExtractAttachments(obj, "audio")

	if len(atts) == 0 {
		core.ReplySimple(obj, core.ERR_NO_AUDIO)

		return
	}

	attachment := atts[0]

	if attachment.Audio.Duration > 10*60 {
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

	bufftemp := bytes.Buffer{}
	cmdtemp := exec.Command("ffmpeg", "-i", "-", "-af", "atempo="+ratio, "-f", "mp3", "-")

	cmdtemp.Stdin = response.Body
	cmdtemp.Stdout = &bufftemp

	if cmdtemp.Run() != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	buffreverb := bytes.Buffer{}
	cmdreverb := exec.Command("ffmpeg", "-i", "-", "-af", "ladspa=file=tap_reverb:tap_reverb", "-f", "mp3", "-")

	cmdreverb.Stdin = &bufftemp
	cmdreverb.Stdout = &buffreverb

	if cmdreverb.Run() != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	d, err := bassboost.UploadAudio(&buffreverb,
		attachment.Audio.Artist,
		fmt.Sprintf("%s (slowed and reverbed by deafmute bot, ratio=%s)", attachment.Audio.Title, ratio))
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
}
