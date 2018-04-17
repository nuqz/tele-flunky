package access

type Role uint8

const (
	Guest Role = 100 - iota
	Bot
	Banned

	Admin Role = 255 - iota
	Friend
	Known
)
