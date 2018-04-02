package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
)

type Update struct {
	*tgbotapi.Update
}

func (upd *Update) Callback() string {
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}

	return ""
}

func (upd *Update) Command() string {
	if upd.Message != nil {
		return upd.Message.Command()
	}

	if upd.EditedMessage != nil {
		return upd.EditedMessage.Command()
	}

	return ""
}

func (upd *Update) Query() string {
	if upd.InlineQuery != nil {
		return upd.InlineQuery.Query
	}

	return ""
}
