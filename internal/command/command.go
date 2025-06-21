package command

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type Command interface {
	Definition() *discordgo.ApplicationCommand
	HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type ComponentHandler interface {
	ComponentHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var commands = map[string]Command{}
var components = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

func RegisterAll(s *discordgo.Session, config *config.Config, repository *repository.MovieRepository) {
	commands[commandNameAddMovie] = NewAddMovieCommand(repository)
	commands[commandNameAllMovies] = NewAllMoviesCommand(repository)
	commands[commandNameSpin] = NewSpinCommand(repository)
	commands[commandNameActiveMovie] = NewActiveMovieCommand(repository)
	commands[commandNameSetActive] = NewSetActiveCommand(repository)
	commands[commandNameSetWatched] = NewSetWatchedCommand(repository)

	for _, cmd := range commands {
		s.ApplicationCommandCreate(s.State.User.ID, config.GuildId, cmd.Definition())

		if ch, ok := cmd.(ComponentHandler); ok {
			for customID, handler := range ch.ComponentHandlers() {
				components[customID] = handler
			}
		}
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			cmd := i.ApplicationCommandData().Name
			if handler, ok := commands[cmd]; ok {
				handler.HandleCommand(s, i)
			}
		case discordgo.InteractionMessageComponent:
			customId := i.MessageComponentData().CustomID
			key := strings.Split(customId, ":")[0]

			if len(key) == 0 {
				return
			}

			if handler, ok := components[key]; ok {
				handler(s, i)
			}
		case discordgo.InteractionModalSubmit:
			customId := i.ModalSubmitData().CustomID
			args := strings.Split(customId, ":")
			movieTitle := args[2]

			data := i.ModalSubmitData()
			dateInput := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			timeInput := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			timezoneInput := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

			slog.Info("Modal submitted", "custom_id", data.CustomID, "date", dateInput, "time", timeInput, "timezone", timezoneInput)

			loc, err := time.LoadLocation(timezoneInput)
			if err != nil {
				InteractionResponseError(s, i, err, "Invalid timezone provided.")
				return
			}
			parsedTime, err := time.ParseInLocation("2006-01-02T15:04:05", fmt.Sprintf("%sT%s", dateInput, timeInput), loc)
			if err != nil {
				InteractionResponseError(s, i, err, "Invalid date or time format. Please use RFC3339 format.")
				return
			}

			endingTime := parsedTime.Add(2 * time.Hour)

			scheduledEvent, err := s.GuildScheduledEventCreate(config.GuildId, &discordgo.GuildScheduledEventParams{
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
	})
}
