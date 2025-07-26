package component

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type CreateEventPreferredTime struct {
	MovieRepository *repository.Movie
	UserRepository  *repository.User
	Config          *config.Config
}

func NewCreateEventPreferredTime(
	movieRepository *repository.Movie,
	userRepository *repository.User,
	config *config.Config,
) *CreateEventPreferredTime {
	return &CreateEventPreferredTime{
		MovieRepository: movieRepository,
		UserRepository:  userRepository,
		Config:          config,
	}
}

func CreateEventPreferredtimeButton(movieID, movieTitle string) discordgo.Button {
	return discordgo.Button{
		Style:    discordgo.PrimaryButton,
		Label:    "Create Event (Preferred Time)",
		CustomID: fmt.Sprintf("%s:%s:%s", CustomIdCreateEventPreferredTime, movieID, movieTitle),
	}
}

func (c *CreateEventPreferredTime) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := strings.Split(i.MessageComponentData().CustomID, ":")
	movieTitle := args[2]

	user, err := c.UserRepository.UserByUserId(i.Member.User.ID)
	if err != nil {
		InteractionResponseError(s, i, err, "failed to get user settings")
		return
	}

	parsedTime, err := time.Parse("15:04", user.PreferredTimeOfDay)
	if err != nil {
		InteractionResponseError(s, i, err, "failed to get parse time")
		return
	}

	loc, err := time.LoadLocation(user.PreferredTimezone)
	if err != nil {
		InteractionResponseError(s, i, err, "failed to load timezone")
		return
	}

	now := time.Now().In(loc)

	target := time.Date(
		now.Year(), now.Month(), now.Day(),
		parsedTime.Hour(), parsedTime.Minute(), 0, 0,
		loc,
	)

	daysUntil := (dayOfWeekMap[user.PreferredDayOfWeek] - int(now.Weekday()) + 7) % 7
	if daysUntil == 0 && now.After(target) {
		daysUntil = 7
	}

	targetTime := target.AddDate(0, 0, daysUntil)
	endingTime := targetTime.Add(2 * time.Hour)

	scheduledEvent, err := s.GuildScheduledEventCreate(c.Config.GuildID, &discordgo.GuildScheduledEventParams{
		Name:               "Wheel of Wonder",
		Description:        movieTitle,
		ScheduledStartTime: &targetTime,
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
			Content: fmt.Sprintf("You created an event for %s at %v", movieTitle, target),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		InteractionResponseError(s, i, err, "failed to create event modal")
		return
	}
}
