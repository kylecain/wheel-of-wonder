package util

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/model"
)

func MovieEmbed(movie *model.Movie) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       movie.Title,
		Description: movie.Description,
		Image:       &discordgo.MessageEmbedImage{URL: movie.ImageURL},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Added by",
				Value:  movie.Username,
				Inline: true,
			},
			{
				Name:   "Source",
				Value:  movie.ContentURL,
				Inline: true,
			},
		},
	}
}
