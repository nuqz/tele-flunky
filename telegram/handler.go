package telegram

type HandlerFunc func(*Bot, *Update)

func (fn HandlerFunc) HandleUpdate(bot *Bot, upd *Update) { fn(bot, upd) }

type Handler interface {
	HandleUpdate(*Bot, *Update)
}
