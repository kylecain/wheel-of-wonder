package util

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/model"
)

func InteractionGroup(i *discordgo.InteractionCreate) slog.Attr {
	return slog.Group("interaction",
		slog.String("id", i.ID),
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

func UserGroup(u *model.User) slog.Attr {
	return slog.Group("user",
		slog.Int64("id", u.ID),
		slog.String("user_id", u.UserID),
		slog.String("username", u.Username),
		slog.String("preferred_day_of_week", u.PreferredDayOfWeek),
		slog.String("preferredTimeOfDay", u.PreferredTimeOfDay),
		slog.String("preferred_timezone", u.PreferredTimezone),
		slog.Time("created_at", u.CreatedAt),
		slog.Time("updated_at", u.UpdatedAt),
	)
}
