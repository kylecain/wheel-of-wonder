package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kylecain/wheel-of-wonder/internal/bot"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db"
)

func main() {
	handlerOptions := slog.HandlerOptions{Level: slog.LevelInfo}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &handlerOptions))
	httpClient := http.DefaultClient

	config, err := config.NewConfig(logger)
	if err != nil {
		logger.Error("failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	db, err := db.NewDatabase(config)
	if err != nil {
		logger.Error("failed to initalize database", slog.Any("err", err))
		os.Exit(1)
	}

	bot, err := bot.NewBot(config, db, httpClient, logger)
	if err != nil {
		logger.Error("failed to create bot", slog.Any("err", err))
		os.Exit(1)
	}
	bot.Start()
	defer bot.Stop()

	logger.Info("bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
