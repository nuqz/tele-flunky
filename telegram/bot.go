package telegram

import (
	"errors"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/nuqz/tele-flunky/storage"
	"gopkg.in/telegram-bot-api.v4"
)

const (
	EnvCert   = "TG_BOT_CERT"
	EnvDomain = "TG_BOT_DOMAIN"
	EnvKey    = "TG_BOT_KEY"
	EnvSecret = "TG_BOT_SECRET"
	EnvToken  = "TG_BOT_TOKEN"
)

var (
	numProcs = runtime.GOMAXPROCS(0)

	ErrNoBotCert   = errors.New("Please, provide path to valid SSL certificate for Telegram bot.")
	ErrNoBotKey    = errors.New("Please, provide path to valid SSL key for Telegram bot.")
	ErrNoBotSecret = errors.New("Please, provide Telegram bot secret with TG_BOT_SECRET environment variable.")
	ErrNoBotToken  = errors.New("Please, provide Telegram bot token with TG_BOT_TOKEN environment variable.")
)

type BotConfig struct {
	certPath string
	debug    bool
	domain   string
	keyPath  string
	secret   string
	token    string
}

func ConfigFromEnv() (*BotConfig, error) {
	cfg := new(BotConfig)

	cfg.token = os.Getenv(EnvToken)
	if cfg.token == "" {
		return nil, ErrNoBotToken
	}

	cfg.secret = os.Getenv(EnvSecret)
	if cfg.secret == "" {
		return nil, ErrNoBotSecret
	}

	cfg.domain = os.Getenv(EnvDomain)
	if cfg.domain != "" {
		cfg.certPath = os.Getenv(EnvCert)
		if cfg.certPath == "" {
			return nil, ErrNoBotCert
		}

		cfg.keyPath = os.Getenv(EnvKey)
		if cfg.keyPath == "" {
			return nil, ErrNoBotCert
		}
	}

	return cfg, nil
}

func (cfg *BotConfig) IsWebhook() bool { return cfg.domain != "" }

type Bot struct {
	*tgbotapi.BotAPI

	config  *BotConfig
	done    chan struct{}
	updates tgbotapi.UpdatesChannel

	cmds          map[string]Handler
	cbQueries     map[string]Handler
	inlineQueries map[string]Handler

	Storage storage.BotStorage
}

func NewBot(s storage.BotStorage, cfg *BotConfig, debug bool) (*Bot, error) {
	tgBotAPI, err := tgbotapi.NewBotAPI(cfg.token)
	if err != nil {
		return nil, err
	}
	tgBotAPI.Debug = debug

	bot := &Bot{
		BotAPI: tgBotAPI,

		config: cfg,

		cmds:          map[string]Handler{},
		cbQueries:     map[string]Handler{},
		inlineQueries: map[string]Handler{},

		Storage: s,
	}

	if cfg.IsWebhook() {
		_, err := bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://"+cfg.domain+":8443/"+bot.Token, cfg.certPath))
		if err != nil {
			return nil, err
		}

		bot.updates = bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServeTLS("0.0.0.0:8443", cfg.certPath, cfg.keyPath, nil)
	} else {
		info, err := bot.GetWebhookInfo()
		if err != nil {
			return nil, err
		}

		if info.IsSet() {
			_, err = bot.RemoveWebhook()
			if err != nil {
				return nil, err
			}
		}

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		if bot.updates, err = bot.GetUpdatesChan(u); err != nil {
			return nil, err
		}
	}

	return bot, nil
}

func NewBotEnv(s storage.BotStorage, debug bool) (*Bot, error) {
	cfg, err := ConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return NewBot(s, cfg, debug)
}

func (bot *Bot) Updates() <-chan *Update {
	out := make(chan *Update)
	go func() {
		for upd := range bot.updates {
			out <- NewUpdate(&upd)
		}
		close(out)
	}()
	return out
}

func (bot *Bot) Command(name string, h Handler)       { bot.cmds[name] = h }
func (bot *Bot) CallbackQuery(name string, h Handler) { bot.cbQueries[name] = h }
func (bot *Bot) InlineQuery(query string, h Handler)  { bot.inlineQueries[query] = h }

func (bot *Bot) DefaultHandler() Handler {
	return HandlerFunc(func(bot *Bot, upd *Update) error {
		if h, ok := bot.cmds[upd.Command()]; ok {
			return h.HandleUpdate(bot, upd)
		}

		if h, ok := bot.inlineQueries[upd.InlineQuery()]; ok {
			return h.HandleUpdate(bot, upd)
		}

		if h, ok := bot.cbQueries[upd.CallbackQuery()]; ok {
			return h.HandleUpdate(bot, upd)
		}

		return nil
	})
}

func (bot *Bot) Serve(h Handler) {
	if bot.done == nil {
		bot.done = make(chan struct{}, numProcs)
		updates := bot.Updates()

		for i := 0; i < numProcs; i++ {
			go func() {
			eternity:
				for {
					select {
					case upd := <-updates:
						if err := h.HandleUpdate(bot, upd); err != nil {
							log.Println(err)
						}
					case <-bot.done:
						break eternity
					}
				}
			}()
		}
	}
}

func (bot *Bot) Stop() {
	for i := 0; i < numProcs; i++ {
		bot.done <- struct{}{}
	}
}

func (bot *Bot) AnswerCallback(upd *Update, text string, alert bool) error {
	if _, err := bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
		CallbackQueryID: upd.Update.CallbackQuery.ID,
		Text:            text,
		ShowAlert:       alert,
	}); err != nil {
		return err
	}

	return nil
}

func (bot *Bot) UpdateCallbackQueryMessage(
	upd *Update,
	caption, text string,
	markup *tgbotapi.InlineKeyboardMarkup,
) error {
	if caption != "" {
		if _, err := bot.Send(tgbotapi.NewEditMessageCaption(
			upd.Update.CallbackQuery.Message.Chat.ID,
			upd.Update.CallbackQuery.Message.MessageID,
			caption,
		)); err != nil {
			return err
		}
	}

	if text != "" {
		if _, err := bot.Send(tgbotapi.NewEditMessageText(
			upd.Update.CallbackQuery.Message.Chat.ID,
			upd.Update.CallbackQuery.Message.MessageID,
			text,
		)); err != nil {
			return err
		}
	}

	if markup != nil {
		if _, err := bot.Send(tgbotapi.NewEditMessageReplyMarkup(
			upd.Update.CallbackQuery.Message.Chat.ID,
			upd.Update.CallbackQuery.Message.MessageID,
			*markup,
		)); err != nil {
			return err
		}
	}

	return nil
}
