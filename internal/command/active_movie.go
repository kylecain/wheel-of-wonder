package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type ActiveMovie struct {
	MovieRepository *repository.Movie
}

func NewActiveMovie(movieRepository *repository.Movie) *ActiveMovie {
	return &ActiveMovie{
		MovieRepository: movieRepository,
	}
}

func (c *ActiveMovie) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameActiveMovie,
		Description: "Show the active movie",
	}
}

func (c *ActiveMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	activeMovie, err := c.MovieRepository.GetActive(i.GuildID)
	if err != nil || activeMovie == nil {
		util.InteractionResponseError(s, i, err, "Failed to retrieve active movie.")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: util.MovieEmbedSlice(activeMovie),
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		util.InteractionResponseError(s, i, err, "Failed to respond to add command")
	}
}
