package setrole

import (
	"context"
	"vkbot/core"
	"vkbot/core/rolesystem"

	"github.com/SevereCloud/vksdk/v2/events"
)

func RegisterAdmin() core.Command {
	return core.Command{
		Aliases:     []string{"админ"},
		Description: "выдать роль администратора пользователю",
		Handler: func(ctx *context.Context, obj *events.MessageNewObject) {
			handle(ctx, obj, rolesystem.ROLE_ADMINISTRATOR)
		},
	}
}

func RegisterModerator() core.Command {
	return core.Command{
		Aliases:     []string{"модератор"},
		Description: "выдать роль модератора пользователю",
		Handler: func(ctx *context.Context, obj *events.MessageNewObject) {
			handle(ctx, obj, rolesystem.ROLE_MODERATOR)
		},
	}
}

func RegisterPurge() core.Command {
	return core.Command{
		Aliases:     []string{"отозвать"},
		Description: "отозвать роль пользователя",
		Handler: func(ctx *context.Context, obj *events.MessageNewObject) {
			handle(ctx, obj, rolesystem.ROLE_MEMBER)
		},
	}
}
