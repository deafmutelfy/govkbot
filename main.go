package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
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

	if s.Cfg.UserToken != "" {
		s.UserVk = api.NewVK(s.Cfg.UserToken)
	}

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
		log.Fatalln(errors.New("unknown event manager: " + s.Cfg.EventManager))
	}

	h.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		handle(ctx, obj, nil)
	})

	if s.Cfg.EventManager == "callback" {
		http.HandleFunc("/callback", h.(eventManagerCallback).HandleFunc)
		http.ListenAndServe(s.Cfg.Host+":"+s.Cfg.Port, nil)
	} else {
		h.(eventManagerLongpoll).Run()
	}
}
