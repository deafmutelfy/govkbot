package main

import (
	"vkbot/commands/actions"
	"vkbot/commands/base"
	"vkbot/commands/bassboost"
	"vkbot/commands/cm"
	"vkbot/commands/commemoration"
	"vkbot/commands/curse"
	"vkbot/commands/dota2"
	"vkbot/commands/help"
	"vkbot/commands/linus"
	"vkbot/commands/mashup"
	"vkbot/commands/nick"
	"vkbot/commands/online"
	"vkbot/commands/ping"
	"vkbot/commands/rptool"
	"vkbot/commands/rule34"
	"vkbot/commands/slowreverb"
	"vkbot/commands/soyjack"
	"vkbot/commands/tacticalpic"
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
		dota2.Register(),
		commemoration.Register(),
		soyjack.Register(),
		tacticalpic.Register(),
		base.Register(),
		actions.Register(),
		cm.Register(),
		bassboost.Register(),
		slowreverb.Register(),
		mashup.Register(),
		rptool.Register(),
	}
}
