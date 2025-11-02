package util

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/model"
)

func InteractionResponseError(s *discordgo.Session, i *discordgo.InteractionCreate, err error, message string) {
	slog.Error(message, "error", err)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("failed to respond to interaction", "error", err)
	}
}

func MovieEmbed(movie *model.Movie) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       movie.Title,
		Description: movie.Description,
		Image:       &discordgo.MessageEmbedImage{URL: movie.ImageURL},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Added by",
				Value:  movie.Username,
				Inline: true,
			},
			{
				Name:   "Source",
				Value:  movie.ContentURL,
				Inline: true,
			},
		},
	}
}

func MovieEmbedSlice(movie *model.Movie) []*discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       movie.Title,
		Description: movie.Description,
		Image:       &discordgo.MessageEmbedImage{URL: movie.ImageURL},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Added by",
				Value:  movie.Username,
				Inline: true,
			},
			{
				Name:   "Source",
				Value:  movie.ContentURL,
				Inline: true,
			},
		},
	}

	return []*discordgo.MessageEmbed{embed}
}

func RespondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

const (
	maxNameLength        = 100
	maxDescriptionLength = 1000
)

func ScheduleEvent(movie *model.Movie, imageData string, startTime time.Time, endTime time.Time, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	name, description := movie.Title, movie.Description

	if len(movie.Title) > maxNameLength {
		name = movie.Title[:maxNameLength]
	}

	if len(movie.Description) > maxDescriptionLength {
		description = movie.Description[:maxDescriptionLength]
	}

	scheduledEvent, err := s.GuildScheduledEventCreate(i.GuildID, &discordgo.GuildScheduledEventParams{
		Name:               name,
		Description:        description,
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
