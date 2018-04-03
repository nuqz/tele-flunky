package telegram

import (
	"errors"
	"log"
	"net/http"
	"os"
	"runtime"

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

	callbacks map[string]Handler
	commands  map[string]Handler
	queries   map[string]Handler
	config    *BotConfig
	updates   tgbotapi.UpdatesChannel
	done      chan struct{}
}

func NewBot(cfg *BotConfig, debug bool) (*Bot, error) {
	tgBotAPI, err := tgbotapi.NewBotAPI(cfg.token)
	if err != nil {
		return nil, err
	}
	tgBotAPI.Debug = debug

	bot := &Bot{
		BotAPI:    tgBotAPI,
		config:    cfg,
		callbacks: map[string]Handler{},
		commands:  map[string]Handler{},
		queries:   map[string]Handler{},
	}

	if cfg.IsWebhook() {
		_, err := bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://"+cfg.domain+":8443/"+bot.Token, cfg.certPath))
		if err != nil {
			return nil, err
		}

		bot.updates = bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServeTLS("0.0.0.0:8443", cfg.certPath, cfg.keyPath, nil)
	} else {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		if bot.updates, err = bot.GetUpdatesChan(u); err != nil {
			return nil, err
		}
	}

	return bot, nil
}

func NewBotEnv(debug bool) (*Bot, error) {
	cfg, err := ConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return NewBot(cfg, debug)
}

func (bot *Bot) Updates() <-chan *Update {
	out := make(chan *Update)
	go func() {
		for upd := range bot.updates {
			out <- &Update{&upd}
		}
		close(out)
	}()
	return out
}

func (bot *Bot) Callback(name string, h Handler) { bot.callbacks[name] = h }
func (bot *Bot) Command(name string, h Handler)  { bot.commands[name] = h }
func (bot *Bot) Query(query string, h Handler)   { bot.queries[query] = h }

func (bot *Bot) Serve() {
	if bot.done == nil {
		bot.done = make(chan struct{}, numProcs)
		updates := bot.Updates()

		for i := 0; i < numProcs; i++ {
			go func() {
			eternity:
				for {
					select {
					case upd := <-updates:
						if h, ok := bot.commands[upd.Command()]; ok {
							h.HandleUpdate(bot, upd)
							continue
						}

						if h, ok := bot.queries[upd.Query()]; ok {
							h.HandleUpdate(bot, upd)
							continue
						}

						if h, ok := bot.callbacks[upd.Callback()]; ok {
							h.HandleUpdate(bot, upd)
							continue
						}

						if upd.Message != nil {
							log.Println("Sorry, don't understand!")
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
