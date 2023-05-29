package cm

import (
	"errors"
	"vkbot/core"
	"vkbot/subsystems/rolesystem"

	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"чм"},
		Description: "команды для чат-менеджмента",
		Metacommand: true,
		Subcommands: &[]core.Command{
			{
				Aliases:     []string{"кик"},
				Description: "исключить участника беседы",
				Handler:     kick,
			},
			{
				Aliases:     []string{"бан"},
				Description: "забанить участника беседы",
				Handler:     ban,
			},
			{
				Aliases:     []string{"разбан"},
				Description: "разбанить пользователя в беседе",
				Handler:     unban,
			},
			{
				Aliases:     []string{"инит"},
				Description: "инициализировать систему ролей в беседе",
				Handler:     initrole,
			},
			{
				Aliases:     []string{"роль"},
				Description: "узнать свою роль в беседе",
				Handler:     getrole,
			},
			{
				Aliases:     []string{"состав"},
				Description: "состав участников беседы по ролям",
				Handler:     listrole,
			},
			{
				Aliases:     []string{"админ"},
				Description: "выдать роль администратора пользователю",
				Handler: func(obj *events.MessageNewObject) (err error) {
					return setrole(obj, rolesystem.ROLE_ADMINISTRATOR)
				},
			},
			{
				Aliases:     []string{"модератор"},
				Description: "выдать роль модератора пользователю",
				Handler: func(obj *events.MessageNewObject) (err error) {
					return setrole(obj, rolesystem.ROLE_MODERATOR)
				},
			},
			{
				Aliases:     []string{"снять"},
				Description: "отозвать роль пользователя",
				Handler: func(obj *events.MessageNewObject) (err error) {
					return setrole(obj, rolesystem.ROLE_MEMBER)
				},
			},
			{
				Aliases:     []string{"рп"},
				Description: "включить/отключить RP команды (обнять, т.п.) в этом чате",
				Handler:     rp,
			},
			{
				Aliases:     []string{"приветствие"},
				Description: "изменить текущее приветствие",
				Handler:     greeting,
			},
		},
	}
}

func cmInit(obj *events.MessageNewObject) error {
	if obj.Message.PeerID == obj.Message.FromID {
		return errors.New(core.ERR_NO_DM)
	}

	return nil
}
