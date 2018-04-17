package middleware

import (
	"errors"

	tg "github.com/nuqz/tele-flunky/telegram"
	"github.com/nuqz/tele-flunky/telegram/access"
)

var ErrAccessDenied = errors.New("Access denied")

func AllowTo(h tg.Handler, roles ...access.Role) tg.Handler {
	return tg.HandlerFunc(func(ctx *tg.Context) error {
		if ctx.User.Role == access.Admin {
			return h.HandleUpdate(ctx)
		} else if ctx.User.Role == access.Banned {
			return ErrAccessDenied
		}

		for _, role := range roles {
			if ctx.User.Role == role {
				return h.HandleUpdate(ctx)
			}
		}

		return ErrAccessDenied
	})
}

