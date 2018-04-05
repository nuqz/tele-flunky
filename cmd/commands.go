package main

import (
	tg "github.com/nuqz/tele-flunky/telegram"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	StartCommand = "start"

	Greeting = `️Hello, dear ✌ My name is Mr. Flunky and I am the bot. I can do many fun things with my friends:	
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
			CallbackData: &UnknownCallback,
		}),
	)
)

func NewStartCommand() tg.Handler { return tg.HandlerFunc(startCommand) }

func startCommand(bot *tg.Bot, upd *tg.Update) error {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, Greeting)
	msg.ReplyMarkup = startKeyboardMarkup
	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

