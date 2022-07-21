package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/callback"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
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

	cmds := commandPool()
	s.CommandPool = &cmds

	res, _ := s.Vk.GroupsGetByID(nil)
	s.Cfg.GroupId = res[0].ID

	var h eventManager
	if s.Cfg.EventManager == "callback" {
		cb := callback.NewCallback()
		cb.ConfirmationKey = s.Cfg.ConfirmationKey

		h = cb
	} else if s.Cfg.EventManager == "longpoll" {
		h, err = longpoll.NewLongPoll(s.Vk, s.Cfg.GroupId)

		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln(errors.New("unknowd event manager: " + s.Cfg.EventManager))
	}

	h.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		if obj.Message.Action.Type == "chat_invite_user" {
			handleChatInviteUser(&obj)

			return
		}

		tokens := strings.Split(obj.Message.Text, " ")
		if len(tokens) == 0 {
			return
		}
		if len(tokens[0]) <= 1 {
			return
		}

		targetcmd := tokens[0]

		for _, x := range cmds {
			for _, a := range x.Aliases {
				if (((targetcmd == a) && (x.NoPrefix)) || ((targetcmd[1:] == a) && (!x.NoPrefix))) && !x.Hidden {
					if (targetcmd[1:] == a) && (!x.NoPrefix) {
						targetcmd = targetcmd[1:]
					}

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

	if s.Cfg.EventManager == "callback" {
		http.HandleFunc("/callback", h.(eventManagerCallback).HandleFunc)
		http.ListenAndServe(s.Cfg.Host+":"+s.Cfg.Port, nil)
	} else {
		h.(eventManagerLongpoll).Run()
	}
}
