package telegram

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nuqz/tele-flunky/telegram/access"
	"github.com/nuqz/tele-flunky/telegram/storage"
)

const (
	EnvCert     = "TG_BOT_CERT"
	EnvDomain   = "TG_BOT_DOMAIN"
	EnvKey      = "TG_BOT_KEY"
	EnvSecret   = "TG_BOT_SECRET"
	EnvToken    = "TG_BOT_TOKEN"
	EnvProxyURL = "TG_BOT_PROXY_URL"
)

var (
	numProcs = runtime.GOMAXPROCS(0)

	ErrNoBotCert   = errors.New("Please, provide path to valid SSL certificate for Telegram bot.")
	ErrNoBotKey    = errors.New("Please, provide path to valid SSL key for Telegram bot.")
	ErrNoBotSecret = errors.New("Please, provide Telegram bot secret with TG_BOT_SECRET environment variable.")
	ErrNoBotToken  = errors.New("Please, provide Telegram bot token with TG_BOT_TOKEN environment variable.")
)

type BotConfig struct {
	certPath  string
	debug     bool
	domain    string
	keyPath   string
	secret    string
	token     string
	proxy_url string
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

	cfg.proxy_url = os.Getenv(EnvProxyURL)
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
	messages      map[string]Handler

	Storage storage.BotStorage
}

func NewBot(s storage.BotStorage, cfg *BotConfig, debug bool) (*Bot, error) {
	var (
		tgbot *tgbotapi.BotAPI
		err   error
	)
	if cfg.proxy_url != "" {
		tgbot, err = tgbotapi.NewBotAPIWithClient(cfg.token, &http.Client{
			Transport: &http.Transport{
				Proxy: func(*http.Request) (*url.URL, error) {
					return url.Parse(cfg.proxy_url)
				},
			},
		})
	} else {
		tgbot, err = tgbotapi.NewBotAPI(cfg.token)
	}
	if err != nil {
		return nil, err
	}
	tgbot.Debug = debug

	bot := &Bot{
		BotAPI: tgbot,

		config: cfg,

		cmds:          map[string]Handler{},
		cbQueries:     map[string]Handler{},
		inlineQueries: map[string]Handler{},
		messages:      map[string]Handler{},

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
func (bot *Bot) Message(msg string, h Handler)        { bot.messages[msg] = h }

func (bot *Bot) DefaultHandler() Handler {
	return HandlerFunc(func(ctx *Context) error {
		var handler Handler
		if h, ok := bot.cmds[ctx.Update.Command()]; ok {
			handler = h
		} else if h, ok := bot.inlineQueries[ctx.Update.InlineQuery()]; ok {
			handler = h
		} else if ctx.Update.IsInlineQuery() {
			parts := strings.Split(ctx.Update.InlineQuery(), " ")
			if len(parts) > 1 {
				if h, ok := bot.inlineQueries[parts[0]]; ok {
					handler = h
				}
			}
		} else if h, ok := bot.cbQueries[ctx.Update.CallbackQuery()]; ok {
			handler = h
		} else {
			handlerName, err := bot.getUserNextChatMessageHandler(ctx)
			if err != nil {
				return err
			}

			if h, ok := bot.messages[handlerName]; ok {
				handler = h
			}
		}

		if handler != nil {
			return handler.HandleUpdate(ctx)
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
						ctx, err := bot.newContext(upd)
						if err != nil {
							log.Println(err)
							continue
						}

						log.Println(ctx.Update)

						if err := h.HandleUpdate(ctx); err != nil {
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

func (bot *Bot) getUserNextChatMessageHandler(ctx *Context) (string, error) {
	log.Printf("%+v", ctx.Update)
	handlerName, err := bot.Storage.GetUserNextChatMessageHandler(
		ctx.User,
		ctx.Update.Chat(),
	)

	if err != nil {
		return "", err
	}

	return handlerName, nil
}

func (bot *Bot) user(upd *Update) (*access.User, error) {
	tgUser := upd.User()

	var user *access.User
	if has, err := bot.Storage.HasUser(tgUser); err != nil {
		return nil, err
	} else if has {
		if user, err = bot.Storage.GetUser(tgUser); err != nil {
			return nil, err
		}
	} else {
		user = access.NewUser(tgUser)
		user.FromChatID = upd.Update.Message.Chat.ID

		user.Role = access.Known
		if user.IsBot {
			user.Role = access.Bot
		}

		if err = bot.Storage.PutUser(user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (bot *Bot) newContext(upd *Update) (*Context, error) {
	user, err := bot.user(upd)
	if err != nil {
		return nil, err
	}

	return &Context{
		Context: context.Background(),
		Bot:     bot,
		Update:  upd,
		User:    user,
	}, nil
}

func (bot *Bot) SendMessage(
	ctx *Context,
	text string,
	markup *tgbotapi.InlineKeyboardMarkup,
) error {
	msg := tgbotapi.NewMessage(ctx.Update.Chat().ID, text)
	if markup != nil {
		msg.ReplyMarkup = markup
	}
	_, err := bot.Send(msg)
	return err
}

func (bot *Bot) SendSticker(ctx *Context, stickerID string) error {
	sticker := tgbotapi.NewStickerShare(
		ctx.Update.Chat().ID,
		StickerCriminalRaccoonHat,
	)
	_, err := bot.Send(sticker)
	return err
}
