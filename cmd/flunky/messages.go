package main

import (
	"fmt"

	tg "github.com/nuqz/tele-flunky/telegram"
)

var (
	StickerMessage = "sticker"
)

func NewStickerIDMessage() tg.Handler { return tg.HandlerFunc(stickerIDMessage) }
func stickerIDMessage(ctx *tg.Context) error {
	if err := ctx.Bot.Storage.SetUserNextChatMessageHandler(ctx.User,
		ctx.Update.Chat(), ""); err != nil {
		return err
	}

	if ctx.Update.Message != nil && ctx.Update.Message.Sticker != nil {
		if err := ctx.Bot.SendMessage(ctx, fmt.Sprintf(
			"Telegram sticker ID: %s",
			ctx.Update.Update.Message.Sticker.FileID,
		), nil); err != nil {
			return err
		}
	}

	return nil
}
