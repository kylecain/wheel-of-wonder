package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/component"
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

func (c *DeleteMovie) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameDeleteMovie,
		Description: "Select movie from list to delete",
	}
}

func (c *DeleteMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := c.MovieRepository.GetAllUnwatched(i.GuildID)
	if err != nil || len(movies) == 0 {
		InteractionResponseError(s, i, err, "no unwatched movies available to delete")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Pick one of the options below:",
			Components: component.MovieSelectMenu(movies, component.CustomIdDeleteMovie, i),
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		InteractionResponseError(s, i, err, "failed to create delete movie select menu")
		return
	}
}
