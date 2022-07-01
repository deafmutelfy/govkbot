package main

import (
	"vkbot/commands/nick"
	"vkbot/commands/online"
	"vkbot/commands/ping"
	"vkbot/core"
)

type poolType []core.Command

func commandPool() poolType {
	return poolType{
		ping.Register(),
		nick.Register(),
		online.Register(),
	}
}
