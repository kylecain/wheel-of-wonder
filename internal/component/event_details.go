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

type EventDetails struct {
	MovieRepository *repository.Movie
	MovieService    *service.Movie
}

func NewEventDetails(movieRepository *repository.Movie, movieService *service.Movie) *EventDetails {
	return &EventDetails{
		MovieRepository: movieRepository,
		MovieService:    movieService,
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

	dateInput := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	timeInput := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	timezoneInput := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	slog.Info("Modal submitted", "custom_id", data.CustomID, "date", dateInput, "time", timeInput, "timezone", timezoneInput)

	loc, err := time.LoadLocation(timezoneInput)
	if err != nil {
		util.InteractionResponseError(s, i, err, "Invalid timezone provided.")
		return
	}

	startTime, err := time.ParseInLocation("2006-01-02T15:04", fmt.Sprintf("%sT%s", dateInput, timeInput), loc)
	if err != nil {
		util.InteractionResponseError(s, i, err, "Invalid date or time format.")
		return
	}

	endTime := startTime.Add(EventDuration)

	imageData, err := c.MovieService.FetchImageAndEncode(selectedMovie.ImageURL)
	if err != nil {
		slog.Error("Failed to fetch and encode image", "error", err)
	}

	err = util.ScheduleEvent(selectedMovie, imageData, startTime, endTime, s, i)
	if err != nil {
		slog.Error("Failed to schedule and notify for event", "error", err)
	}
}
