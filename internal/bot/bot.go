package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/command"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/handler"
)

type Bot struct {
	Session *discordgo.Session
}

func NewBot(config *config.Config) (*Bot, error) {
	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		return nil, err
	}
	return &Bot{Session: dg}, nil
}

func (b *Bot) Start() error {
	err := b.Session.Open()
	if err != nil {
		return err
	}

	b.Session.AddHandler(handler.NewMovieHandler().Add)
	_, err = b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", command.NewMovieCommand().Add())

	return nil
}

func (b *Bot) Stop() error {
	return b.Session.Close()
}
