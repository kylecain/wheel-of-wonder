package command

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/component"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type BonusMovie struct {
	logger *slog.Logger
}

func NewBonusMovie(logger *slog.Logger) *BonusMovie {
	return &BonusMovie{
		logger: logger,
	}
}

func (c *BonusMovie) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameBonusMovie,
		Description: "Schedule a movie not tracked by the wheel of wonder.",
	}
}

func (c *BonusMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := c.logger.
		With(slog.String("command_name", i.ApplicationCommandData().Name)).
		With(util.InteractionGroup(i))

	l.Info("received command interaction")

	components, err := util.ScheduleEventModal(true)
	if err != nil {
		l.Error("error creating schedule event modal", slog.Any("err", err))
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   component.CustomIdBonusMovieModal,
			Title:      "Create Event For Bonus Movie",
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		l.Error("failed to respond to interaction", slog.Any("err", err))
	} else {
		l.Info("successfully responded to command")
	}
}
