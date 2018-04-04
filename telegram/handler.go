package telegram

type Handler interface {
	HandleUpdate(*Bot, *Update) error
}

type HandlerFunc func(*Bot, *Update) error

func (fn HandlerFunc) HandleUpdate(bot *Bot, upd *Update) error {
	return fn(bot, upd)
}
