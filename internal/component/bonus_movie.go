package component

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/service"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type BonusMovie struct {
	movieService *service.Movie
	logger       *slog.Logger
}

func NewBonusMovie(movieService *service.Movie, logger *slog.Logger) *BonusMovie {
	return &BonusMovie{
		movieService: movieService,
		logger:       logger,
	}
}

func (c *BonusMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := c.logger.
		With(slog.String("component", CustomIdBonusMovieModal)).
		With(util.InteractionGroup(i))

	data := i.ModalSubmitData()

	dateInput := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	timeInput := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	timezoneInput := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	movieInput := data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	loc, err := time.LoadLocation(timezoneInput)
	if err != nil {
		l.Error("invalid timezone provided", slog.String("timezone_input", timezoneInput), slog.Any("err", err))
		util.InteractionResponseError(s, i, err, "Invalid timezone provided.")
		return
	}

	startTime, err := time.ParseInLocation("2006-01-02T15:04", fmt.Sprintf("%sT%s", dateInput, timeInput), loc)
	if err != nil {
		l.Error("unable to parse time", slog.String("date_input", dateInput), slog.String("time_input", timeInput), slog.Any("err", err))
		util.InteractionResponseError(s, i, err, "Invalid date or time provided.")
		return
	}

	endTime := startTime.Add(EventDuration)

	movieData, err := c.movieService.FetchMovie(movieInput)
	if err != nil {
		l.Error("failed to fetch movie", slog.String("movie_input", movieInput), slog.Any("err", err))
		util.InteractionResponseError(s, i, err, "Failed to fetch movie data.")
		return
	}

	imageData, err := c.movieService.FetchImageAndEncode(movieData.ImageURL)
	if err != nil {
		l.Error("failed to fetch and encode movie data", slog.String("movie_input", movieInput), slog.Any("err", err))
		util.InteractionResponseError(s, i, err, "Failed to encode movie image.")
	}

	err = util.ScheduleEvent(movieData.Title, movieData.Description, imageData, startTime, endTime, s, i)
	if err != nil {
		l.Error("failed to schedule and notify for event", slog.Any("err", err))
	} else {
		l.Info("successfully responded to interaction")
	}

}
