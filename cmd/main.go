package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	tg "github.com/nuqz/tele-flunky/telegram"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	dbPath = flag.String("dbpath", "/tmp/tele-flunky", "Path to bot's database")
	debug  = flag.Bool("debug", false, "Debug Telegram bot APIs")

	err  error
	quit chan os.Signal
	bot  *tg.Bot
)

func init() {
	flag.Parse()

	if *debug {
		log.SetFlags(log.Ltime | log.Lshortfile)
	}

	if bot, err = tg.NewBotEnv(*debug); err != nil {
		log.Fatal(err)
	}

	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Kill, os.Interrupt)
}

func main() {
	bot.Command("authorize", tg.HandlerFunc(func(bot *tg.Bot, upd *tg.Update) {
		log.Printf("%+v", upd)
	}))


	bot.Serve()

	<-quit
	bot.Stop()
}
