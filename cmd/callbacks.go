package main

import (
	tg "github.com/nuqz/tele-flunky/telegram"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	HomeCallback    = "home"
	UnknownCallback = "unknown"

	x                  = "xs"
	homeKeyboardMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: "Catalog", CallbackData: &x},
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: "Search", CallbackData: &x},
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: "About", CallbackData: &x},
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: "Contacts", CallbackData: &x},
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: "Delivery", CallbackData: &x},
		),
	)
)

func homeCallback(bot *tg.Bot, upd *tg.Update) error {
	if err = bot.UpdateCallbackQueryMessage(
		upd,
		"",
		"homepage text",
		&homeKeyboardMarkup,
	); err != nil {
		return err
	}

	if err := bot.AnswerCallback(upd, "Homepage", false); err != nil {
		return err
	}

	return nil
}
