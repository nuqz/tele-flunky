package middleware

import (
	"errors"

	tg "github.com/nuqz/tele-flunky/telegram"
	"github.com/nuqz/tele-flunky/telegram/access"
)

var ErrAccessDenied = errors.New("Access denied")

func AllowTo(h tg.Handler, roles ...access.Role) tg.Handler {
	return tg.HandlerFunc(func(ctx *tg.Context) error {
		if ctx.User.IsAdmin() {
			return h.HandleUpdate(ctx)
		} else if ctx.User.IsBanned() {
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

func AllowAbove(h tg.Handler, role access.Role) tg.Handler {
	return tg.HandlerFunc(func(ctx *tg.Context) error {
		if ctx.User.IsAdmin() {
			return h.HandleUpdate(ctx)
		} else if ctx.User.IsBanned() {
			return ErrAccessDenied
		}

		if ctx.User.Role = role {
			return h.HandleUpdate(ctx)
		}

		return ErrAccessDenied
	})
}
