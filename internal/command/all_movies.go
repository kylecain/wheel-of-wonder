package command

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type AllMovies struct {
	MovieRepository *repository.Movie
}

func NewAllMovies(movieRepository *repository.Movie) *AllMovies {
	return &AllMovies{
		MovieRepository: movieRepository,
	}
}

func (c *AllMovies) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameAllMovies,
		Description: "Get all movies in the wheel",
	}
}

func (c *AllMovies) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := c.MovieRepository.GetAll(i.GuildID)
	if err != nil {
		util.InteractionResponseError(s, i, err, "Failed to get all movies.")
		return
	}

	var b strings.Builder
	for _, movie := range movies {
		b.WriteString(fmt.Sprintf("**%s** — %s\n", movie.Title, movie.Username))
	}

	embed := &discordgo.MessageEmbed{
		Title:       "All Movies",
		Description: b.String(),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d movies", len(movies)),
		},
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("failed to respond to getall command", "error", err)
	}
}
