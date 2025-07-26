package command

import (
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/component"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type Spin struct {
	MovieRepository *repository.Movie
}

func NewSpin(movieRepository *repository.Movie) *Spin {
	return &Spin{
		MovieRepository: movieRepository,
	}
}

func (c *Spin) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameSpin,
		Description: "Spin the wheel and get a random movie",
	}
}

func (c *Spin) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := c.MovieRepository.GetAll(i.GuildID)
	if err != nil || len(movies) == 0 {
		InteractionResponseError(s, i, err, "failed to get all movies for spin")
		return
	}

	selectedMovie := movies[rand.Intn(len(movies))]
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You spun the wheel and got: %s", selectedMovie.Title),
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						component.SetActiveButton(selectedMovie),
					},
				},
			},
		},
	})

	if err != nil {
		slog.Error("failed to respond to spin command", "error", err)
	}
}
