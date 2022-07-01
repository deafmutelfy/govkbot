package main

import (
	"vkbot/commands/help"
	"vkbot/commands/nick"
	"vkbot/commands/online"
	"vkbot/commands/ping"
	"vkbot/core"
)

func commandPool() core.PoolType {
	return core.PoolType{
		ping.Register(),
		nick.Register(),
		online.Register(),
		help.Register(),
	}
}
