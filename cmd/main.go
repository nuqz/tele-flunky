package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/nuqz/tele-flunky/storage"
	tg "github.com/nuqz/tele-flunky/telegram"
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

	s, err := storage.NewLevelDBStorage(os.Getenv("TG_BOT_LEVELDB_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	if bot, err = tg.NewBotEnv(s, *debug); err != nil {
		log.Fatal(err)
	}

	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Kill, os.Interrupt)
}

func main() {
	bot.Command(StartCommand, NewStartCommand())

	bot.CallbackQuery(HomeCallback, tg.HandlerFunc(homeCallback))
	bot.CallbackQuery(AloneCallback, tg.HandlerFunc(aloneCallback))
	bot.CallbackQuery(StickerIDCallback, tg.HandlerFunc(stickerIDCallback))
	bot.Message(StickerMessage, tg.HandlerFunc(stickerIDMessage))

	bot.Serve(bot.DefaultHandler())

	<-quit
	bot.Stop()
}
