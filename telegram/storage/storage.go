package storage

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/nuqz/tele-flunky/telegram/access"
)

type UserStorage interface {
	PutUser(*access.User) error
	HasUser(*tgbotapi.User) (bool, error)
	GetUser(*tgbotapi.User) (*access.User, error)
	DeleteUser(*access.User) error
}

type UserStateStorage interface {
	SetUserNextChatMessageHandler(*access.User, *tgbotapi.Chat, string) error
	GetUserNextChatMessageHandler(*access.User, *tgbotapi.Chat) (string, error)
}

type BotStorage interface {
	UserStorage
	UserStateStorage
}
