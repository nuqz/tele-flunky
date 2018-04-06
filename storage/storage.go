package storage

import (
	"gopkg.in/telegram-bot-api.v4"

	"github.com/nuqz/tele-flunky/access"
)

type UserStorage interface {
	PutUser(*access.User) error
	HasUser(*tgbotapi.User) (bool, error)
	GetUser(*tgbotapi.User) (*access.User, error)
	DeleteUser(*access.User) error
}

type BotStorage interface {
	UserStorage
}
