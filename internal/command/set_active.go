package command

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type SetActive struct {
	MovieRepository *repository.Movie
}

func NewSetActive(movieRepository *repository.Movie) *SetActive {
	return &SetActive{
		MovieRepository: movieRepository,
	}
}

func (c *SetActive) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameSetActive,
		Description: "Set a movie as active by ID",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "movie_id",
				Description: "Movie ID that will be set as active",
				Required:    true,
			},
		},
	}
}

func (c *SetActive) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].IntValue()

	err := c.MovieRepository.UpdateActive(input, true)
	if err != nil {
		InteractionResponseError(s, i, err, "Failed to update active movie.")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You set %d as active", input),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("failed to respond to add command", "error", err)
	}
}
