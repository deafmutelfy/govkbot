package main

import (
	"vkbot/commands/commemoration"
	"vkbot/commands/curse"
	"vkbot/commands/dota2"
	"vkbot/commands/getrole"
	"vkbot/commands/help"
	"vkbot/commands/initrole"
	"vkbot/commands/kick"
	"vkbot/commands/linus"
	"vkbot/commands/listrole"
	"vkbot/commands/nick"
	"vkbot/commands/online"
	"vkbot/commands/ping"
	"vkbot/commands/rule34"
	"vkbot/commands/setrole"
	"vkbot/commands/top"
	"vkbot/commands/tts"
	"vkbot/commands/who"
	"vkbot/core"
)

func commandPool() []core.Command {
	return []core.Command{
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
		dota2.Register(),
		setrole.RegisterAdmin(),
		setrole.RegisterModerator(),
		setrole.RegisterPurge(),
		listrole.Register(),
		commemoration.Register(),
	}
}
