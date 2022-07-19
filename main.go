package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/callback"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/go-redis/redis/v9"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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
	s.CommandPool = &cmds

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
				if targetcmd == a && !x.Hidden {
					if !x.Metacommand {
						go x.Handler(&ctx, &obj)
					} else {
						if len(tokens) < 2 {
							core.ReplySimple(&obj, generateHelp(a, x.Subcommands))

							return
						}

						targetcmd := tokens[1]
						launched := false

						for _, v := range *x.Subcommands {
							if targetcmd == v.Aliases[0] && !v.Hidden {
								launched = true
								go v.Handler(&ctx, &obj)
							}
						}

						if !launched {
							core.ReplySimple(&obj, generateHelp(a, x.Subcommands))
						}
					}
				}
			}
		}
	})

	http.HandleFunc("/callback", cb.HandleFunc)
	http.ListenAndServe("0.0.0.0:"+s.Cfg.Port, nil)
}
