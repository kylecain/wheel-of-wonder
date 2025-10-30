package component

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type CreateEventPreferredTime struct {
	MovieRepository *repository.Movie
	UserRepository  *repository.User
	HttpClient      *http.Client
}

func NewCreateEventPreferredTime(
	movieRepository *repository.Movie,
	userRepository *repository.User,
	httpClient *http.Client,
) *CreateEventPreferredTime {
	return &CreateEventPreferredTime{
		MovieRepository: movieRepository,
		UserRepository:  userRepository,
		HttpClient:      httpClient,
	}
}

func CreateEventPreferredtimeButton(movieID, movieTitle string) discordgo.Button {
	return discordgo.Button{
		Style:    discordgo.PrimaryButton,
		Label:    "Create Event (Preferred Time)",
		CustomID: fmt.Sprintf("%s:%s:%s", CustomIdCreateEventPreferredTime, movieID, movieTitle),
	}
}

const eventDuration = 2 * time.Hour

func nextPreferredEventTime(now time.Time, preferredDay string, preferredTime time.Time, loc *time.Location) time.Time {
	dayOfWeekLower := strings.ToLower(preferredDay)
	weekdayTarget, ok := dayOfWeekMap[dayOfWeekLower]
	if !ok {
		weekdayTarget = int(now.Weekday())
	}
	daysUntil := (weekdayTarget - int(now.Weekday()) + 7) % 7
	targetDay := now.AddDate(0, 0, daysUntil)
	target := time.Date(
		targetDay.Year(), targetDay.Month(), targetDay.Day(),
		preferredTime.Hour(), preferredTime.Minute(), 0, 0,
		loc,
	)
	if daysUntil == 0 && now.After(target) {
		target = target.AddDate(0, 0, 7)
	}
	return target
}

func (c *CreateEventPreferredTime) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := strings.Split(i.MessageComponentData().CustomID, ":")
	movieIdStr := args[1]

	movieId, err := strconv.Atoi(movieIdStr)
	if err != nil {
		util.InteractionResponseError(s, i, err, "Failed to convert movieID")
		return
	}

	selectedMovie, err := c.MovieRepository.GetMovieByID(movieId)
	if err != nil {
		util.InteractionResponseError(s, i, err, "failed to get movie by ID")
		return
	}

	user, err := c.UserRepository.UserByUserId(i.Member.User.ID)
	if err != nil {
		util.InteractionResponseError(s, i, err, "failed to get user settings")
		return
	}

	parsedTime, err := time.Parse("15:04", user.PreferredTimeOfDay)
	if err != nil {
		util.InteractionResponseError(s, i, err, "failed to parse preferred time")
		return
	}

	loc, err := time.LoadLocation(user.PreferredTimezone)
	if err != nil {
		util.InteractionResponseError(s, i, err, "failed to load timezone")
		return
	}

	now := time.Now().In(loc)
	target := nextPreferredEventTime(now, user.PreferredDayOfWeek, parsedTime, loc)
	endingTime := target.Add(eventDuration)

	imageData, err := util.FetchAndEncodeImage(selectedMovie.ImageURL, *c.HttpClient)
	if err != nil {
		slog.Error("Failed to fetch and encode image", "error", err)
	}

	scheduledEvent, err := s.GuildScheduledEventCreate(i.GuildID, &discordgo.GuildScheduledEventParams{
		Name:               selectedMovie.Title,
		Description:        selectedMovie.Description,
		Image:              imageData,
		ScheduledStartTime: &target,
		ScheduledEndTime:   &endingTime,
		EntityType:         discordgo.GuildScheduledEventEntityTypeExternal,
		EntityMetadata: &discordgo.GuildScheduledEventEntityMetadata{
			Location: "Online",
		},
		PrivacyLevel: discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
	})
	if err != nil {
		slog.Error("Failed to create schedule event", "error", err)
		return
	}

	slog.Info("Scheduled event created", "event", scheduledEvent)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You created an event for %s at %v", selectedMovie.Title, target),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("Failed to update user on scheduled event", "error", err)
		return
	}

	eventURL := fmt.Sprintf("https://discord.com/events/%s/%s", i.GuildID, scheduledEvent.ID)
	_, err = s.ChannelMessageSend(i.ChannelID, eventURL)
	if err != nil {
		slog.Error("Failed to send event link to general chat", "error", err)
	}
}
