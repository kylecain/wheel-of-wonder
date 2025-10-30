package command

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/model"
	"github.com/kylecain/wheel-of-wonder/internal/service"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type AddMovie struct {
	MovieRepository    *repository.Movie
	SearchMovieService *service.MovieSearch
}

func NewAddMovie(movieRepository *repository.Movie, searchMovieService *service.MovieSearch) *AddMovie {
	return &AddMovie{
		MovieRepository:    movieRepository,
		SearchMovieService: searchMovieService,
	}
}

func (c *AddMovie) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameAddMovie,
		Description: "Add a movie",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "movie",
				Description: "movie that will be added to the wheel",
				Required:    true,
			},
		},
	}
}

func (c *AddMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].StringValue()

	movieInfo, err := c.SearchMovieService.FetchMovie(input)
	if err != nil {
		util.InteractionResponseError(s, i, err, "Failed to fetch movie info.")
		return
	}

	movie := &model.Movie{
		GuildID:     i.GuildID,
		UserID:      i.Member.User.ID,
		Username:    i.Member.User.Username,
		Title:       movieInfo.Title,
		Description: movieInfo.Description,
		ImageURL:    movieInfo.ImageURL,
		ContentURL:  movieInfo.ContentURL,
	}

	_, err = c.MovieRepository.AddMovie(movie)
	if err != nil {
		util.InteractionResponseError(s, i, err, "Failed to add movie.")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: util.MovieEmbedSlice(movie),
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("Failed to respond to add-movie command", "error", err)
	}
}
