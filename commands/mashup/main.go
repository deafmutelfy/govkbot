package mashup

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"vkbot/core"
	"vkbot/subsystems/audiosystem"

	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"мэшап", "мешап"},
		Description: "слить две аудиозаписи в одну",
		Handler:     handle,
		Queue: &core.CommandQueueParams{
			Name: "audio",
		},
	}
}

func handle(obj *events.MessageNewObject) (err error) {
	args := core.ExtractArguments(obj)

	p, err := parseParams(&args)
	if err != nil {
		core.ReplySimple(obj, fmt.Sprint("ошибка:", err))
		err = nil

		return
	}

	atts := core.ExtractAttachments(obj, "audio")
	if len(atts) == 0 {
		core.ReplySimple(obj, core.ERR_NO_AUDIO)

		return
	}
	if len(atts) != 2 {
		core.ReplySimple(obj, "ошибка: аудиозаписей должно быть две")

		return
	}

	var audios []string
	for _, x := range atts {
		if x.Audio.Duration > 10*60 {
			core.ReplySimple(obj, core.ERR_10MIN)

			return
		}
		if x.Audio.URL == "" {
			core.ReplySimple(obj, core.ERR_AUDIO_VK_API_BUG)

			return
		}

		response, err := http.Get(x.Audio.URL)
		if err != nil {
			core.ReplySimple(obj, core.ERR_UNKNOWN)

			return err
		}

		f, err := os.CreateTemp("", "govkbot-mashup-cache")
		if err != nil {
			core.ReplySimple(obj, core.ERR_UNKNOWN)

			return err
		}
		defer os.Remove(f.Name())

		if _, err := io.Copy(f, response.Body); err != nil {
			core.ReplySimple(obj, core.ERR_UNKNOWN)
		}

		audios = append(audios, f.Name())
	}

	b := []string{"-i", audios[0], "-i", audios[1]}
	if p.Longest {
		b = append(b, "-longest")
	} else {
		b = append(b, "-shortest")
	}
	b = append(b,
		"-filter_complex",
		fmt.Sprintf("[0:a]adelay=%d|%d,bass=g=%d,treble=g=%d[a0]; [1:a]adelay=%d|%d,bass=g=%d,treble=g=%d[a1]; [a0][a1]amerge=inputs=2[out]",
			p.FirstOffset,
			p.FirstOffset,
			p.FirstBass,
			p.FirstTreble,
			p.SecondOffset,
			p.SecondOffset,
			p.SecondBass,
			p.SecondTreble,
		),
		"-map", "[out]", "-ac", "2", "-f", "mp3", "-",
	)

	buff := bytes.Buffer{}
	cmd := exec.Command("ffmpeg", b...)

	cmd.Stdout = &buff

	if cmd.Run() != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	d, err := audiosystem.UploadAudio(&buff,
		atts[0].Audio.Artist+" & "+atts[1].Audio.Artist,
		fmt.Sprintf("%s & %s (bassboosted by deafmute bot)", atts[0].Audio.Title, atts[1].Audio.Title))
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
