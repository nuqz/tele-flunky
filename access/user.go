package access

import (
	"fmt"

	"gopkg.in/telegram-bot-api.v4"
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
