package access

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type User struct {
	*tgbotapi.User

	FromChatID  int64  `json:"from_chat_id"`
	FriendToken string `json:"friend_token,omitempty"`
	Role        Role   `json:"role"`
}

func NewUser(user *tgbotapi.User) *User {
	role := Guest
	if user.IsBot {
		role = Bot
	}
	return &User{user, 0, "", role}
}

func (u *User) StorageKey() []byte {
	return []byte(fmt.Sprintf("user_%d", u.ID))
}
