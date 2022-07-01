package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/callback"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/go-redis/redis/v9"
)

func main() {
	log.Println(os.Environ())
	s := core.GetStorage()
	s.Cfg = core.Config{}

	if err := s.Cfg.Load("config.yaml"); err != nil {
		log.Fatalln(err)
	}

	s.Vk = api.NewVK(s.Cfg.Token)

	opt, err := redis.ParseURL(s.Cfg.RedisUrl)
	if err != nil {
		log.Fatalln(err)
	}

	s.Db = redis.NewClient(opt)
	s.Ctx = context.Background()

	cb := callback.NewCallback()
	cb.ConfirmationKey = s.Cfg.ConfirmationKey

	cmds := commandPool()

	cb.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		tokens := strings.Split(obj.Message.Text, " ")
		if len(tokens) == 0 {
			return
		}
		if len(tokens[0]) <= 1 {
			return
		}
		if tokens[0][0] != '/' {
			return
		}

		targetcmd := tokens[0][1:]

		for _, x := range cmds {
			for _, a := range x.Aliases {
				if targetcmd == a {
					go x.Handler(&ctx, &obj)
				}
			}
		}
	})

	http.HandleFunc("/callback", cb.HandleFunc)
	http.ListenAndServe("0.0.0.0:"+s.Cfg.Port, nil)
}
