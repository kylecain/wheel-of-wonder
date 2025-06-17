package command

import (
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type SpinCommand struct {
	MovieRepository *repository.MovieRepository
}

func NewSpinCommand(movieRepository *repository.MovieRepository) *SpinCommand {
	return &SpinCommand{
		MovieRepository: movieRepository,
	}
}

func (c *SpinCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "spin",
		Description: "spin the wheel and get a random movie",
	}
}

func (c *SpinCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	movies, err := c.MovieRepository.GetAll(i.GuildID)
	if err != nil {
		slog.Error("failed to get all movies for spin", "error", err)
		return
	}

	if len(movies) == 0 {
		response := "No movies available to spin."
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
		if err != nil {
			slog.Error("failed to respond to spin command", "error", err)
		}
		return
	}

	selectedMovie := movies[rand.Intn(len(movies))]
	response := fmt.Sprintf("You spun the wheel and got: %s", selectedMovie.Title)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		slog.Error("failed to respond to spin command", "error", err)
	}
}

func (c *AddMovieCommand) ComponentHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate) {
	return map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"set_active_movie": c.setActiveMovieHandler,
	}
}

func (c *AddMovieCommand) setActiveMovieHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
}
