package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	yo "github.com/nuqz/go-yobit"
	tg "github.com/nuqz/tele-flunky/telegram"
)

var (
	YobitQuery = "yobit"
)

func yobitQuery(ctx *tg.Context) error {
	tickerName := "btc_usd"

	parts := strings.Split(ctx.Update.InlineQuery(), " ")
	if len(parts) > 1 {
		tickerName = parts[1]
	}

	api := yo.NewPublicAPI()
	ticker, err := api.Ticker(tickerName)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(`
***Yobit*** %s (%s)

Sell: %.8g	Buy: %.8g
Min: %.8g	Average: %.8g	Max: %.8g	Last: %.8g

Volume: %.8g	Current: %.8g
`,
		tickerName, time.Unix(ticker.Updated, 0),
		ticker.Sell, ticker.Buy,
		ticker.Low, ticker.Avg, ticker.High, ticker.Last,
		ticker.Vol, ticker.VolCur,
	)

	_, err = ctx.Bot.AnswerInlineQuery(tgbotapi.InlineConfig{
		InlineQueryID: ctx.Update.Update.InlineQuery.ID,
		Results: []interface{}{
			tgbotapi.NewInlineQueryResultArticle(tickerName, tickerName, msg),
		},
	})

	return err
}
