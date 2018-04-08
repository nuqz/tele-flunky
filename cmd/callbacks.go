package main

import (
	tg "github.com/nuqz/tele-flunky/telegram"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	HomeCallback    = "home"
	UnknownCallback = "unknown"
	StickerIDCallback = "sticker_id"

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
			tgbotapi.InlineKeyboardButton{
				Text:         "I want to know Sticker ID",
				CallbackData: &StickerIDCallback,
			},
		),
	)
)

	if err = bot.UpdateCallbackQueryMessage(
		upd,
		"",
		"homepage text",
		&homeKeyboardMarkup,
func aloneCallback(ctx *tg.Context) error {
func homeCallback(ctx *tg.Context) error {
func stickerIDCallback(ctx *tg.Context) error {
	if err := ctx.Bot.Storage.SetUserNextChatMessageHandler(
		ctx.User,
		ctx.Update.Chat().ID,
		"sticker",
	); err != nil {
		return err
	}

	if err = ctx.Bot.UpdateCallbackQueryMessage(
		ctx.Update,
		"",
		"Send a sticker to me and I'll send it's ID back to you.",
		nil,
	); err != nil {
		return err
	}

	return ctx.Bot.AnswerCallback(ctx.Update, "Send a sticker to me.", false)
}
