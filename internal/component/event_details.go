package component

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/config"
)

type EventDetails struct {
	Config *config.Config
}

func NewEventDetails(config *config.Config) *EventDetails {
	return &EventDetails{
		Config: config,
	}
}

func EventDetailsModal() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "date",
					Label:       "Enter the date of the event",
					Style:       discordgo.TextInputShort,
					Placeholder: "YYYY-MM-DD",
					Required:    true,
					MaxLength:   10,
					MinLength:   10,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "time",
					Label:       "Enter the time of the event (24-hour format)",
					Style:       discordgo.TextInputShort,
					Placeholder: "HH:mm",
					Required:    true,
					MaxLength:   5,
					MinLength:   5,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "timezone",
					Label:       "Enter the timezone of the event",
					Style:       discordgo.TextInputShort,
					Placeholder: "America/Chicago",
					Required:    true,
					MaxLength:   50,
					MinLength:   8,
				},
			},
		},
	}
}

func (c *EventDetails) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	customId := data.CustomID
	args := strings.Split(customId, ":")
	movieTitle := args[2]

	dateInput := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	timeInput := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	timezoneInput := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	slog.Info("Modal submitted", "custom_id", data.CustomID, "date", dateInput, "time", timeInput, "timezone", timezoneInput)

	loc, err := time.LoadLocation(timezoneInput)
	if err != nil {
		InteractionResponseError(s, i, err, "Invalid timezone provided.")
		return
	}
	parsedTime, err := time.ParseInLocation("2006-01-02T15:04", fmt.Sprintf("%sT%s", dateInput, timeInput), loc)
	if err != nil {
		InteractionResponseError(s, i, err, "Invalid date or time format.")
		return
	}

	endingTime := parsedTime.Add(2 * time.Hour)

	scheduledEvent, err := s.GuildScheduledEventCreate(c.Config.GuildId, &discordgo.GuildScheduledEventParams{
		Name:               "Wheel of Wonder",
		Description:        movieTitle,
		ScheduledStartTime: &parsedTime,
		ScheduledEndTime:   &endingTime,
		EntityType:         discordgo.GuildScheduledEventEntityTypeExternal,
		EntityMetadata: &discordgo.GuildScheduledEventEntityMetadata{
			Location: "Online",
		},
		PrivacyLevel: discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
	})

	if err != nil {
		InteractionResponseError(s, i, err, "Failed to create scheduled event.")
		return
	}

	slog.Info("Scheduled event created", "event", scheduledEvent)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You created an event for %s", movieTitle),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("failed to respond to create event modal", "error", err)
		return
	}
}
