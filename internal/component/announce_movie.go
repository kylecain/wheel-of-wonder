package component

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type AnnounceMovie struct {
	MovieRepository *repository.Movie
}

func NewAnnounceMovie(movieRepository *repository.Movie) *AnnounceMovie {
	return &AnnounceMovie{
		MovieRepository: movieRepository,
	}
}

func AnnounceMovieButton(movieID, movieTitle string) discordgo.Button {
	return discordgo.Button{
		Style:    discordgo.PrimaryButton,
		Label:    "Announce Movie",
		CustomID: fmt.Sprintf("%s:%s:%s", CustomIdAnnounceMovie, movieID, movieTitle),
	}
}

func (c *AnnounceMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := strings.Split(i.MessageComponentData().CustomID, ":")
	_, movieTitle := args[1], args[2]

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Announced: %s", movieTitle),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("Failed to respond to interaction", "error", err)
	}

	message := fmt.Sprintf("Now Playing: %s", movieTitle)
	_, err = s.ChannelMessageSend(i.ChannelID, message)
	if err != nil {
		slog.Error("Failed to send message", "error", err)
	}
}
