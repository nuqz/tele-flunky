package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"

	tg "github.com/nuqz/tele-flunky/telegram"
	"github.com/nuqz/tele-flunky/telegram/access"
)

var (
	StartCommand = "start"
	AdminCommand = "admin"

	Greeting = `️Hello, dear ✌ My name is **Mr. Flunky** and I am the bot. I can do many fun things with my friends:	
🎼 I like to listen music and you can share it with me.
🎯 I follow the situation at the crypto exchange.
😎 I'm very cool bot and 💖 u!`

	startKeyboardMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.InlineKeyboardButton{
			Text:         "Let's be friends 🤝",
			CallbackData: &HomeCallback,
		}),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.InlineKeyboardButton{
			Text:         "Don't bother me 👽",
			CallbackData: &AloneCallback,
		}),
	)
)

func NewStartCommand() tg.Handler { return tg.HandlerFunc(startCommand) }
func startCommand(ctx *tg.Context) error {
	if ctx.User.Role < access.Friend {
		err := ctx.Bot.SendSticker(ctx, tg.StickerMinecraftFabulous)
		if err != nil {
			return err
		}

		return ctx.Bot.SendMessage(ctx, Greeting, &startKeyboardMarkup)
	}

	return homeCallback(ctx)
}
