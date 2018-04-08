package main

import (
	"gopkg.in/telegram-bot-api.v4"

	tg "github.com/nuqz/tele-flunky/telegram"
)

var (
	StartCommand = "start"

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
		return err
	}

	return nil
}
