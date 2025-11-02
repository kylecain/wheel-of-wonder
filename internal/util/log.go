package util

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/model"
)

func InteractionGroup(i *discordgo.InteractionCreate) slog.Attr {
	return slog.Group("interaction",
		slog.String("id", i.ID),
		slog.String("command", i.ApplicationCommandData().Name),
		slog.String("username", i.Member.User.Username),
		slog.String("global_name", i.Member.User.GlobalName),
		slog.String("guild_id", i.GuildID),
	)
}

func MovieGroup(m *model.Movie) slog.Attr {
	return slog.Group("movie",
		slog.Int64("id", m.ID),
		slog.String("guild_id", m.GuildID),
		slog.String("user_id", m.UserID),
		slog.String("username", m.Username),
		slog.String("title", m.Title),
		slog.String("description", m.Description),
		slog.String("image_url", m.ImageURL),
		slog.String("content_url", m.ContentURL),
		slog.Bool("watched", m.Watched),
		slog.Bool("active", m.Active),
		slog.Time("created_at", m.CreatedAt),
		slog.Time("updated_at", m.UpdatedAt),
	)
}

func MovieInfoGroup(m *model.MovieInfo) slog.Attr {
	return slog.Group("movieInfo",
		slog.String("title", m.Title),
		slog.String("description", m.Description),
		slog.String("image_url", m.ImageURL),
		slog.String("content_url", m.ContentURL),
	)
}
