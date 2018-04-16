package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Update struct {
	*tgbotapi.Update
	user    *tgbotapi.User
	chat    *tgbotapi.Chat
	cmd     string
	cmdArgs string
	cbQuery string

	inlineQuery   string
	isInlineQuery bool
}

func updateFromMessage(update *Update, msg *tgbotapi.Message) {
	update.user = msg.From
	update.chat = msg.Chat
	if msg.IsCommand() {
		update.cmd = msg.Command()
		update.cmdArgs = msg.CommandArguments()
	}
}

func updateFromCallbackQuery(update *Update, query *tgbotapi.CallbackQuery) {
	update.user = query.From
	update.chat = query.Message.Chat
	update.cbQuery = query.Data
}

func updateFromInlineQuery(update *Update, query *tgbotapi.InlineQuery) {
	update.user = query.From
	update.inlineQuery = query.Query
	update.isInlineQuery = true
}

func NewUpdate(upd *tgbotapi.Update) (update *Update) {
	update = new(Update)
	update.Update = upd

	if upd.Message != nil {
		updateFromMessage(update, upd.Message)
		return
	}

	if upd.EditedMessage != nil {
		updateFromMessage(update, upd.EditedMessage)
		return
	}

	if upd.CallbackQuery != nil {
		updateFromCallbackQuery(update, upd.CallbackQuery)
		return
	}

	if upd.InlineQuery != nil {
		updateFromInlineQuery(update, upd.InlineQuery)
	}

	return
}

func (upd *Update) IsCommand() bool       { return upd.cmd != "" }
func (upd *Update) IsCallbackQuery() bool { return upd.cbQuery != "" }
func (upd *Update) IsInlineQuery() bool   { return upd.isInlineQuery }
func (upd *Update) User() *tgbotapi.User  { return upd.user }
func (upd *Update) Chat() *tgbotapi.Chat  { return upd.chat }
func (upd *Update) Command() string       { return upd.cmd }
func (upd *Update) CommandArgs() string   { return upd.cmdArgs }
func (upd *Update) CallbackQuery() string { return upd.cbQuery }
func (upd *Update) InlineQuery() string   { return upd.inlineQuery }
