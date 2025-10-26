package component

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/model"
)

func MovieSelectMenu(movies []model.Movie, customId string, i *discordgo.InteractionCreate) []discordgo.MessageComponent {
	var menuOptions []discordgo.SelectMenuOption

	for _, movie := range movies {
		option := discordgo.SelectMenuOption{
			Label:       movie.Title,
			Value:       fmt.Sprintf("%d", movie.ID),
			Description: movie.Username,
		}
		menuOptions = append(menuOptions, option)
	}

	menu := discordgo.SelectMenu{
		CustomID:    customId,
		Placeholder: "Movies",
		Options:     menuOptions,
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{menu},
		},
	}
}
