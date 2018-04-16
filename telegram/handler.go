package telegram

import (
	"context"

	"github.com/nuqz/tele-flunky/telegram/access"
)

type Context struct {
	context.Context

	Bot    *Bot
	Update *Update
	User   *access.User
}

type Handler interface {
	HandleUpdate(*Context) error
}

type HandlerFunc func(*Context) error

func (fn HandlerFunc) HandleUpdate(ctx *Context) error { return fn(ctx) }
