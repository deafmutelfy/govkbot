package main

import (
	"context"
	"strconv"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/golang-queue/queue"
)

func handle(ctx context.Context, obj events.MessageNewObject, parentcmd *core.Command) {
	s := core.GetStorage()

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
		if x.QueueName == "" {
			go x.Handler(&ctx, &obj)

			return
		}

		q, ok := queuePool[x.QueueName]
		if !ok {
			q = queue.NewPool(3)

			queuePoolMutex.Lock()
			queuePool[x.QueueName] = q
			queuePoolMutex.Unlock()
		}

		if core.IsInArray(queueIds, obj.Message.FromID) {
			core.ReplySimple(&obj, "ошибка: запрос от вас уже получен")

			return
		}

		queueIdsMutex.Lock()
		queueIds = append(queueIds, obj.Message.FromID)
		queueIdsMutex.Unlock()

		b := params.NewMessagesSendBuilder()

		b.DisableMentions(true)

		d, _ := core.Send(&obj,
			"[id"+
				strconv.Itoa(obj.Message.FromID)+
				"|"+
				core.GetNickname(obj.Message.FromID)+
				"], ваш запрос принят в обработку. Номер в очереди: "+
				strconv.Itoa(q.SubmittedTasks()-q.FailureTasks()-q.SuccessTasks()-q.BusyWorkers()+1),
			b)

		queuePoolMutex.Lock()
		q.QueueTask(func(_ context.Context) error {
			x.Handler(&ctx, &obj)

			queueIdsMutex.Lock()
			queueIds = core.Remove(queueIds, obj.Message.FromID)
			queueIdsMutex.Unlock()

			bu := params.NewMessagesDeleteBuilder()

			bu.PeerID(obj.Message.PeerID)
			bu.ConversationMessageIDs([]int{d[0].ConversationMessageID})
			bu.DeleteForAll(true)

			s.Vk.MessagesDelete(bu.Params)

			return nil
		})
		queuePoolMutex.Unlock()
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
}
