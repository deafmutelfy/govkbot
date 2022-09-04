package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"vkbot/core"
	"vkbot/subsystems/queuesystem"

	"github.com/SevereCloud/vksdk/v2/events"
)

func handle(ctx context.Context, obj events.MessageNewObject, parentcmd *core.Command) {
	s := core.GetStorage()

	isBlacklisted, _ := s.Db.Get(s.Ctx, fmt.Sprintf("blacklist.%d", obj.Message.FromID)).Result()
	if isBlacklisted == "true" {
		return
	}

	cmds := s.CommandPool

	if parentcmd != nil {
		cmds = parentcmd.Subcommands
	}

	switch obj.Message.Action.Type {
	case "chat_invite_user":
		handleChatInviteUser(&obj)

		return
	}

	tokens := strings.Split(obj.Message.Text, " ")
	if len(tokens) == 0 {
		return
	}
	if len(tokens[0]) <= 1 {
		if parentcmd != nil {
			core.ReplySimple(&obj, generateHelp(parentcmd.Aliases[0], cmds))
		}

		return
	}

	targetcmd := strings.ToLower(tokens[0])

	launcher := func(x *core.Command) {
		if x.Queue == nil {
			go func() {
				if err := x.Handler(&obj); err != nil {
					log.Println(err)
				}
			}()

			return
		}

		queuesystem.Add(&obj, x.Handler)
	}

	launched := false

	for _, x := range *cmds {
		for _, a := range x.Aliases {
			if (((targetcmd == a) && (x.NoPrefix || parentcmd != nil)) || ((targetcmd[1:] == a) && (!x.NoPrefix))) && !x.Hidden {
				if !x.Metacommand {
					launcher(&x)
				} else {
					obj.Message.Text = strings.Join(tokens[1:], " ")

					handle(ctx, obj, &x)
				}

				launched = true

				return
			}
		}
	}

	if !launched && parentcmd != nil {
		core.ReplySimple(&obj, generateHelp(parentcmd.Aliases[0], cmds))
	}

	handleUserRPAction(&obj)
}
