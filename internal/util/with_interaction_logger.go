package util

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func WithInteractionLogger(l *slog.Logger, i *discordgo.InteractionCreate) *slog.Logger {
	return l.WithGroup("interaction").With(
		slog.String("id", i.ID),
		slog.String("command", i.ApplicationCommandData().Name),
		slog.String("username", i.Member.User.Username),
		slog.String("global_name", i.Member.User.GlobalName),
		slog.String("guild_id", i.GuildID),
	)
}
