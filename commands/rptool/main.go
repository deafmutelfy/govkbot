package rptool

import (
	"vkbot/core"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"рп"},
		Description: "управление пользовательскими РП-командами",
		Metacommand: true,
		Subcommands: &[]core.Command{
			{
				Aliases:     []string{"создать"},
				Description: "создать РП-команду",
				Handler:     create,
			},
			{
				Aliases:     []string{"лист"},
				Description: "получить список созданных вами команд",
				Handler:     list,
			},
			{
				Aliases:     []string{"удалить"},
				Description: "удалить РП-команду",
				Handler:     remove,
			},
		},
	}
}
