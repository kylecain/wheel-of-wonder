package bot

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/command"
	"github.com/kylecain/wheel-of-wonder/internal/component"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/service"
)

type Bot struct {
	session    *discordgo.Session
	config     *config.Config
	db         *sql.DB
	httpClient *http.Client
	logger     *slog.Logger
}

func NewBot(config *config.Config, db *sql.DB, httpClient *http.Client, logger *slog.Logger) (*Bot, error) {
	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	return &Bot{
		session:    dg,
		config:     config,
		db:         db,
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

func (b *Bot) Start() error {
	err := b.session.Open()
	if err != nil {
		return err
	}

	movieRepository := repository.NewMovie(b.db, b.logger)
	userRepository := repository.NewUser(b.db, b.logger)

	movieService := service.NewMovie(b.httpClient, b.logger)

	command.RegisterAll(b.session, b.config, movieRepository, userRepository, movieService, b.logger)
	component.RegisterAll(b.session, movieRepository, userRepository, movieService, b.logger)

	return nil
}

func (b *Bot) Stop() error {
	return b.session.Close()
}
