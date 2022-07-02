package main

import (
	"vkbot/commands/help"
	"vkbot/commands/linus"
	"vkbot/commands/nick"
	"vkbot/commands/online"
	"vkbot/commands/ping"
	"vkbot/commands/top"
	"vkbot/commands/tts"
	"vkbot/commands/who"
	"vkbot/core"
)

func commandPool() core.PoolType {
	return core.PoolType{
		ping.Register(),
		nick.Register(),
		online.Register(),
		help.Register(),
		linus.Register(),
		top.Register(),
		who.Register(),
		tts.Register(),
	}
}
