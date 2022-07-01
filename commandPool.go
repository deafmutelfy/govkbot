package main

import (
	"vkbot/commands/ping"
	"vkbot/core"
)

type poolType []core.Command

func commandPool() poolType {
	return poolType{
		ping.Register(),
	}
}
