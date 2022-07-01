package core

import (
	"sync"

	"github.com/SevereCloud/vksdk/v2/api"
)

var once sync.Once

var storageInstance *Storage = nil

type Storage struct {
	Cfg Config
	Vk  *api.VK
}

func GetStorage() *Storage {
	once.Do(func() {
		storageInstance = &Storage{}
	})

	return storageInstance
}
