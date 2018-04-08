package main

import (
	"fmt"

	"gopkg.in/telegram-bot-api.v4"

	tg "github.com/nuqz/tele-flunky/telegram"
)

var (
	StickerMessage = "sticker"
)

func stickerIDMessage(ctx *tg.Context) error {
	if err := ctx.Bot.Storage.SetUserNextChatMessageHandler(ctx.User,
		ctx.Update.Chat().ID, ""); err != nil {
		return err
	}

	if ctx.Update.Message != nil && ctx.Update.Message.Sticker != nil {
		msg := tgbotapi.NewMessage(
			ctx.Update.Chat().ID,
			fmt.Sprintf("Telegram sticker ID: %s",
				ctx.Update.Update.Message.Sticker.FileID),
		)
		_, err = ctx.Bot.Send(msg)
		return err
	}

	return nil
}
