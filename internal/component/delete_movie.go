package component

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type DeleteMovie struct {
	MovieRepository *repository.Movie
}

func NewDeleteMovie(movieRepository *repository.Movie) *DeleteMovie {
	return &DeleteMovie{
		MovieRepository: movieRepository,
	}
}

func DeleteMovieSelectMenu(movieRepository *repository.Movie, i *discordgo.InteractionCreate) []discordgo.MessageComponent {
	unwatchedMovies, err := movieRepository.GetAllUnwatched(i.GuildID)
	if err != nil {
		slog.Error("Error getting unwatched movies", "error", err)
	}

	var menuOptions []discordgo.SelectMenuOption

	for _, movie := range unwatchedMovies {
		option := discordgo.SelectMenuOption{
			Label:       movie.Title,
			Value:       fmt.Sprintf("%d", movie.ID),
			Description: movie.Username,
		}
		menuOptions = append(menuOptions, option)
	}

	menu := discordgo.SelectMenu{
		CustomID:    CustomIdDeleteMovie,
		Placeholder: "Select movie from list to delete",
		Options:     menuOptions,
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{menu},
		},
	}
}

func (c *DeleteMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selected := i.MessageComponentData().Values[0]
	movieId, err := strconv.Atoi(selected)
	if err != nil {
		InteractionResponseError(s, i, err, "failed to convert movie ID")
		return
	}
	c.MovieRepository.DeleteMovie(int64(movieId))

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    "Movie deleted successfully.",
			Components: []discordgo.MessageComponent{},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		InteractionResponseError(s, i, err, "failed to create event modal")
		return
	}
}
