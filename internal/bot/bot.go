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

	movieRepository := repository.NewMovieRepository(b.DB)

	b.Session.AddHandler(handler.NewMovieHandler(movieRepository).Add)
	b.Session.AddHandler(handler.NewMovieHandler(movieRepository).GetAll)

	_, err = b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.Config.GuildId, command.NewMovieCommand().Add())
	if err != nil {
		return err
	}
	_, err = b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.Config.GuildId, command.NewMovieCommand().GetAll())
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) Stop() error {
	return b.Session.Close()
}
