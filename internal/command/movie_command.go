package command

import "github.com/bwmarrin/discordgo"

type MovieCommand struct{}

func NewMovieCommand() *MovieCommand {
	return &MovieCommand{}
}

func (c *MovieCommand) Add() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "add",
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
