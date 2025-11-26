package command

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/component"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type SetPreferredEventTime struct {
	userRepository *repository.User
	logger         *slog.Logger
}

func NewSetPreferredEventTime(userRepository *repository.User, logger *slog.Logger) *SetPreferredEventTime {
	return &SetPreferredEventTime{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (c *SetPreferredEventTime) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameSetPreferredEventTime,
		Description: "Set your preferred event time",
	}
}

func (c *SetPreferredEventTime) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := c.logger.
		With(slog.String("command_name", i.ApplicationCommandData().Name)).
		With(util.InteractionGroup(i))

	l.Info("received command interaction")

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   component.CustomIDSetPreferredTimeModal,
			Title:      "Set Preferred Time Modal",
			Components: component.SetPreferredEventTimeModal(),
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		l.Error("failed to respond to interaction", slog.Any("err", err))
	} else {
		l.Info("successfully responded to command")
		l.Error("failed to respond to interation")
	}
}
