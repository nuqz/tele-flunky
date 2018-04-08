package main

import (
	"gopkg.in/telegram-bot-api.v4"

	"github.com/nuqz/tele-flunky/access"
	tg "github.com/nuqz/tele-flunky/telegram"
)

var (
	AloneCallback     = "alone"
	HomeCallback      = "home"
	StickerIDCallback = "sticker_id"
	UnknownCallback   = "unknown"

	HomepageMessage = "What we're going to do now?..."
	ComeBackMessage = "I hope to see you again!"

	homeKeyboardMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:         "I want to know Sticker ID",
				CallbackData: &StickerIDCallback,
			},
		),
	)
)

func aloneCallback(ctx *tg.Context) error {
	if ctx.Update.IsCallbackQuery() {
		if err := ctx.Bot.AnswerCallback(
			ctx.Update,
			ComeBackMessage,
			false,
		); err != nil {
			return err
		}
	}

	sticker := tgbotapi.NewStickerShare(
		ctx.Update.Chat().ID,
		tg.StickerMinecraftForeverAlone,
	)
	_, err = ctx.Bot.Send(sticker)
	return err
}

func homeCallback(ctx *tg.Context) error {
	if ctx.Update.IsCallbackQuery() {
		// When new user come from Let's be friends
		if ctx.User.Role < access.Friend {
			ctx.User.Role = access.Friend
			if err = ctx.Bot.Storage.PutUser(ctx.User); err != nil {
				return err
			}
		}

		if err = ctx.Bot.UpdateCallbackQueryMessage(
			ctx.Update,
			"",
			HomepageMessage,
			&homeKeyboardMarkup,
		); err != nil {
			return err
		}

		return ctx.Bot.AnswerCallback(ctx.Update, "Homepage", false)
	} else if ctx.Update.IsCommand() {
		// When friends come from /start
		sticker := tgbotapi.NewStickerShare(
			ctx.Update.Chat().ID,
			tg.StickerCriminalRaccoonHat,
		)
		if _, err = ctx.Bot.Send(sticker); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(ctx.Update.Chat().ID, HomepageMessage)
		msg.ReplyMarkup = homeKeyboardMarkup
		msg.ParseMode = "markdown"
		_, err = ctx.Bot.Send(msg)
		return err
	}

	return nil
}

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
