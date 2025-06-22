package command

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/model"
)

type AddMovie struct {
	MovieRepository *repository.MovieRepository
}

func NewAddMovie(movieRepository *repository.MovieRepository) *AddMovie {
	return &AddMovie{
		MovieRepository: movieRepository,
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

	movie := &model.Movie{
		GuildID:  i.GuildID,
		UserID:   i.Member.User.ID,
		Username: i.Member.User.Username,
		Title:    input,
	}

	_, err := c.MovieRepository.AddMovie(movie)
	if err != nil {
		InteractionResponseError(s, i, err, "Failed to add movie.")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You added: %s", input),
		},
	})
	if err != nil {
		slog.Error("failed to respond to add command", "error", err)
	}
}
