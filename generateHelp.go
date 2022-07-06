package main

import (
	"strings"
	"vkbot/core"
)

func generateHelp(parentcmd string, cmds *[]core.Command) string {
	msg := "список встроенных команд:\n"

	for _, x := range *cmds {
		if x.Hidden {
			continue
		}

		aliases := []string{}

		for _, v := range x.Aliases {
			aliases = append(aliases, "/"+parentcmd+" "+v)
		}

		msg += strings.Join(aliases, ", ") + " - " + x.Description + "\n"
	}

	return msg
}
