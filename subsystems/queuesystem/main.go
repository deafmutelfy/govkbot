package queuesystem

import (
	"context"
	"log"
	"strconv"
	"sync"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/golang-queue/queue"
)

var queuePool *queue.Queue = queue.NewPool(3)
var queueIds []int = []int{}

var queuePoolMutex sync.Mutex
var queueIdsMutex sync.Mutex

func Add(obj *events.MessageNewObject, handler core.CommandHandler) {
	if core.IsInArray(queueIds, obj.Message.FromID) {
		core.ReplySimple(obj, "ошибка: запрос от вас уже получен")

		return
	}

	queueIdsMutex.Lock()
	queueIds = append(queueIds, obj.Message.FromID)
	queueIdsMutex.Unlock()

	b := params.NewMessagesSendBuilder()

	b.DisableMentions(true)

	d, _ := core.Send(obj,
		"[id"+
			strconv.Itoa(obj.Message.FromID)+
			"|"+
			core.GetNickname(obj.Message.FromID)+
			"], ваш запрос принят в обработку. Номер в очереди: "+
			strconv.Itoa(
				queuePool.SubmittedTasks()-
					queuePool.FailureTasks()-
					queuePool.SuccessTasks()-
					queuePool.BusyWorkers()+1),
		b)

	queuePoolMutex.Lock()
	queuePool.QueueTask(func(_ context.Context) error {
		if err := handler(obj); err != nil {
			log.Println(err)
		}

		queueIdsMutex.Lock()
		queueIds = core.Remove(queueIds, obj.Message.FromID)
		queueIdsMutex.Unlock()

		bu := params.NewMessagesDeleteBuilder()

		bu.PeerID(obj.Message.PeerID)
		bu.ConversationMessageIDs([]int{d[0].ConversationMessageID})
		bu.DeleteForAll(true)

		core.GetStorage().Vk.MessagesDelete(bu.Params)

		return nil
	})
	queuePoolMutex.Unlock()
}
