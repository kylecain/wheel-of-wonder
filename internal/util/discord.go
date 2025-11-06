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
	maxEventDescriptionLength = 1000
	maxEventNameLength        = 100
)

func ScheduleEvent(movieTitle, movieDescription, imageData string, startTime time.Time, endTime time.Time, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	name, description := movieTitle, movieDescription

	if len(movieTitle) > maxEventNameLength {
		name = movieTitle[:maxEventNameLength]
	}

	if len(movieDescription) > maxEventDescriptionLength {
		description = movieDescription[:maxEventDescriptionLength]
	}

	fmt.Println(startTime, endTime)

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
			Content: fmt.Sprintf("You created an event for %s", movieTitle),
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

func ScheduleEventModal(shouldIncludeMovie bool) ([]discordgo.MessageComponent, error) {
	tzName := "America/Chicago"
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return nil, err
	}
	now := time.Now().In(loc)

	if !(now.Minute()%30 == 0 && now.Second() == 0 && now.Nanosecond() == 0) {
		add := 30 - (now.Minute() % 30)
		now = now.Add(time.Duration(add) * time.Minute)
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())

	dateDefault := now.Format("2006-01-02")
	timeDefault := now.Format("15:04")

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "date",
					Label:       "Enter the date of the event (YYYY-MM-DD)",
					Style:       discordgo.TextInputShort,
					Placeholder: "YYYY-MM-DD",
					Required:    true,
					MaxLength:   10,
					MinLength:   10,
					Value:       dateDefault,
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
					Value:       timeDefault,
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
					Value:       "America/Chicago",
				},
			},
		},
	}

	if shouldIncludeMovie {
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "movie",
					Label:       "Enter the movie name",
					Style:       discordgo.TextInputShort,
					Placeholder: "Movie title",
					Required:    true,
				},
			},
		})
	}

	return components, nil
}

func MovieSelectMenu(movies []model.Movie, customId string, i *discordgo.InteractionCreate) []discordgo.MessageComponent {
	var menuOptions []discordgo.SelectMenuOption

	for _, movie := range movies {
		option := discordgo.SelectMenuOption{
			Label:       movie.Title,
			Value:       fmt.Sprintf("%d", movie.ID),
			Description: movie.Username,
		}
		menuOptions = append(menuOptions, option)
	}

	menu := discordgo.SelectMenu{
		CustomID:    customId,
		Placeholder: "Movies",
		Options:     menuOptions,
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{menu},
		},
	}
}
