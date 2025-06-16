package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kylecain/wheel-of-wonder/internal/bot"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db"
)

func main() {
	config := config.NewConfig()
	db := db.NewDatabase(config)
	bot, err := bot.NewBot(config, db)
	if err != nil {
		slog.Error("error creating bot", "error", err)
		os.Exit(1)
	}
	bot.Start()
	defer bot.Stop()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
