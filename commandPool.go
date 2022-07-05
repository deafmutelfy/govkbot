package main

import (
	"vkbot/commands/curse"
	"vkbot/commands/getrole"
	"vkbot/commands/help"
	"vkbot/commands/initrole"
	"vkbot/commands/kick"
	"vkbot/commands/linus"
	"vkbot/commands/nick"
	"vkbot/commands/online"
	"vkbot/commands/ping"
	"vkbot/commands/rule34"
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
		rule34.Register(),
		curse.Register(),
		initrole.Register(),
		getrole.Register(),
		kick.Register(),
	}
}
