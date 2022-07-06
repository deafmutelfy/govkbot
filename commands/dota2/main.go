package dota2

import "vkbot/core"

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"д2"},
		Description: "вспомогательные инструменты для Dota 2",
		Metacommand: true,
		Subcommands: &[]core.Command{
			{
				Aliases:     []string{"айди"},
				Description: "установить свой ID игрока",
				Handler:     setId,
			},
			{
				Aliases:     []string{"стат"},
				Description: "получить статистику по профилю",
				Handler:     statistics,
			},
		},
	}
}
