package bot

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/command"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/handler"
)

type Bot struct {
	Session *discordgo.Session
	DB      *sql.DB
}

func NewBot(config *config.Config, db *sql.DB) (*Bot, error) {
	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		return nil, err
	}
	return &Bot{
		Session: dg,
		DB:      db,
	}, nil
}

func (b *Bot) Start() error {
	err := b.Session.Open()
	if err != nil {
		return err
	}

	movieRepository := repository.NewMovieRepository(b.DB)
	b.Session.AddHandler(handler.NewMovieHandler(movieRepository).Add)
	_, err = b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", command.NewMovieCommand().Add())

	return nil
}

func (b *Bot) Stop() error {
	return b.Session.Close()
}
