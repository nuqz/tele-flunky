package access

import (
	"encoding/json"

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

func NewUserFromJSON(data []byte) (*User, error) {
	u := new(User)
	if err := json.Unmarshal(data, u); err != nil {
		return nil, err
	}

	return u, nil

}

func (u *User) IsAdmin() bool  { return u.Role == Admin }
func (u *User) IsBanned() bool { return u.Role == Banned }
