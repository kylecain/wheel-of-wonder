package component

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/service"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type CreateEventPreferredTime struct {
	MovieRepository *repository.Movie
	UserRepository  *repository.User
	MovieService    *service.Movie
}

func NewCreateEventPreferredTime(
	movieRepository *repository.Movie,
	userRepository *repository.User,
	movieService *service.Movie,
) *CreateEventPreferredTime {
	return &CreateEventPreferredTime{
		MovieRepository: movieRepository,
		UserRepository:  userRepository,
		MovieService:    movieService,
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

	startTime, endTime, err := c.getEventStartAndEndTime(i)
	if err != nil {
		util.InteractionResponseError(s, i, err, "Failed to get event start and end time")
		return
	}

	imageData, err := c.MovieService.FetchImageAndEncode(selectedMovie.ImageURL)
	if err != nil {
		slog.Error("Failed to fetch and encode image", "error", err)
	}

	err = util.ScheduleEvent(selectedMovie.Title, selectedMovie.Description, imageData, *startTime, *endTime, s, i)
	if err != nil {
		slog.Error("Failed to schedule and notify for event", "error", err)
	}
}

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

func (c *CreateEventPreferredTime) getEventStartAndEndTime(i *discordgo.InteractionCreate) (*time.Time, *time.Time, error) {

	user, err := c.UserRepository.UserByUserId(i.Member.User.ID)
	if err != nil {
		return nil, nil, err
	}

	parsedTime, err := time.Parse("15:04", user.PreferredTimeOfDay)
	if err != nil {
		return nil, nil, err
	}

	loc, err := time.LoadLocation(user.PreferredTimezone)
	if err != nil {
		return nil, nil, err
	}

	now := time.Now().In(loc)
	startTime := nextPreferredEventTime(now, user.PreferredDayOfWeek, parsedTime, loc)
	endTime := startTime.Add(EventDuration)

	return &startTime, &endTime, nil
}
