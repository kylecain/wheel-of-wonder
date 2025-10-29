package bot

import (
	"database/sql"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/command"
	"github.com/kylecain/wheel-of-wonder/internal/component"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/service"
)

type Bot struct {
	Session *discordgo.Session
	Config  *config.Config
	DB      *sql.DB
}

func NewBot(config *config.Config, db *sql.DB) (*Bot, error) {
	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		return nil, err
	}
	return &Bot{
		Session: dg,
		Config:  config,
		DB:      db,
	}, nil
}

func (b *Bot) Start() error {
	err := b.Session.Open()
	if err != nil {
		return err
	}

	movieRepository := repository.NewMovie(b.DB)
	userRepository := repository.NewUser(b.DB)

	searchMovieService := service.NewMovieSearch(http.DefaultClient)

	command.RegisterAll(b.Session, b.Config, movieRepository, userRepository, searchMovieService)
	component.RegisterAll(b.Session, movieRepository, userRepository)

	return nil
}

func (b *Bot) Stop() error {
	return b.Session.Close()
}
