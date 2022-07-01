package core

import (
	"context"
	"sync"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/go-redis/redis/v9"
)

var once sync.Once

var storageInstance *Storage = nil

type Storage struct {
	Cfg         Config
	Vk          *api.VK
	Db          *redis.Client
	Ctx         context.Context
	CommandPool *PoolType
}

func GetStorage() *Storage {
	once.Do(func() {
		storageInstance = &Storage{}
	})

	return storageInstance
}
