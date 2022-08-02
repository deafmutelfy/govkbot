package bassboost

import (
	"bytes"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/valyala/fastjson"
)

func UploadAudio(out *bytes.Buffer, artist string, title string) (*interface{}, error) {
	s := core.GetStorage()

	res := new(interface{})

	if err := s.UserVk.RequestUnmarshal("audio.getUploadServer", res); err != nil {
		return nil, err
	}

	r, err := s.UserVk.UploadFile(
		(*res).(map[string]interface{})["upload_url"].(string),
		out,
		"file",
		"res.mp3",
	)
	if err != nil {
		return nil, err
	}

	v, err := fastjson.ParseBytes(r)
	if err != nil {
		return nil, err
	}

	server := v.GetInt64("server")
	audio := v.GetStringBytes("audio")
	hash := v.GetStringBytes("hash")

	if err := s.UserVk.RequestUnmarshal("audio.save", res, api.Params{
		"server": server,
		"audio":  string(audio),
		"hash":   string(hash),
		"artist": artist,
		"title":  title,
	}); err != nil {
		return nil, err
	}

	return res, nil
}
