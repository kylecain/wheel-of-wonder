package util

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/model"
)

func ScheduleEvent(movie *model.Movie, imageData string, startTime time.Time, endTime time.Time, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	scheduledEvent, err := s.GuildScheduledEventCreate(i.GuildID, &discordgo.GuildScheduledEventParams{
		Name:               movie.Title,
		Description:        movie.Description,
		Image:              imageData,
		ScheduledStartTime: &startTime,
		ScheduledEndTime:   &endTime,
		EntityType:         discordgo.GuildScheduledEventEntityTypeExternal,
		EntityMetadata: &discordgo.GuildScheduledEventEntityMetadata{
			Location: "Online",
		},
		PrivacyLevel: discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
	})

	if err != nil {
		return err
	}

	slog.Info("Scheduled event created", "event", scheduledEvent)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You created an event for %s", movie.Title),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return err
	}

	eventURL := fmt.Sprintf("https://discord.com/events/%s/%s", i.GuildID, scheduledEvent.ID)

	_, err = s.ChannelMessageSend(i.ChannelID, eventURL)
	if err != nil {
		return err
	}

	return nil
}
