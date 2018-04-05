package access

type Role uint8

const (
	Guest Role = 100 - iota
	Bot

	Admin Role = 255 - iota
	Friend
	Known
)
