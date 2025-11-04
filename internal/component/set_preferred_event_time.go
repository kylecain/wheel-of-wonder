package component

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/model"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

var dayOfWeekMap = map[string]int{
	"sunday":    int(time.Sunday),
	"monday":    int(time.Monday),
	"tuesday":   int(time.Tuesday),
	"wednesday": int(time.Wednesday),
	"thursday":  int(time.Thursday),
	"friday":    int(time.Friday),
	"saturday":  int(time.Saturday),
}

type SetPreferredEventTime struct {
	userRepository *repository.User
	logger         *slog.Logger
}

func NewSetPreferredEventTime(userRepository *repository.User, logger *slog.Logger) *SetPreferredEventTime {
	return &SetPreferredEventTime{
		userRepository: userRepository,
	}
}

func SetPreferredEventTimeModal() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "preferred_day_of_week",
					Label:       "Preferred Day",
					Style:       discordgo.TextInputShort,
					Placeholder: "Thursday",
					Required:    true,
					MaxLength:   20,
					MinLength:   5,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "preferred_time_of_day",
					Label:       "Preferred Time (HH:mm)",
					Style:       discordgo.TextInputShort,
					Placeholder: "19:00",
					Required:    true,
					MaxLength:   5,
					MinLength:   5,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "preferred_timezone",
					Label:       "Preferred Timezone",
					Style:       discordgo.TextInputShort,
					Placeholder: "America/Chicago",
					Required:    true,
					MaxLength:   50,
					MinLength:   8,
				},
			},
		},
	}
}

func (c *SetPreferredEventTime) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	l := c.logger.
		With(slog.String("component", CustomIDSetPreferredTimeModal)).
		With(util.InteractionGroup(i))
	l.Info("received modal submit interaction")

	data := i.ModalSubmitData()

	dayOfWeekInput := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	timeInput := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	timezoneInput := data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	dayOfWeekInputLower := strings.ToLower(dayOfWeekInput)

	_, ok := dayOfWeekMap[dayOfWeekInputLower]

	if !ok {
		l.Warn("invalid day provided", slog.String("day_input", dayOfWeekInput), slog.String("day_input_normalized", dayOfWeekInputLower))
		util.InteractionResponseError(s, i, fmt.Errorf("invalid day"), "Invalid day")
		return
	}

	_, err := time.Parse("15:04", timeInput)
	if err != nil {
		l.Warn("invalid time provided", slog.String("time_input", timeInput), slog.Any("err", err))
		util.InteractionResponseError(s, i, err, "Invalid time provided.")
		return
	}

	_, err = time.LoadLocation(timezoneInput)
	if err != nil {
		l.Warn("invalid timezone provided", slog.String("timezone_input", timezoneInput), slog.Any("err", err))
		util.InteractionResponseError(s, i, err, "Invalid timezone provided.")
		return
	}

	user := &model.User{
		UserID:             i.Member.User.ID,
		Username:           i.Member.User.Username,
		PreferredDayOfWeek: dayOfWeekInput,
		PreferredTimeOfDay: timeInput,
		PreferredTimezone:  timezoneInput,
	}

	_, err = c.userRepository.AddUser(user)
	if err != nil {
		l.Error("error saving settings", slog.Any("err", err), slog.String("user_id", user.UserID))
		util.InteractionResponseError(s, i, err, "Error saving settings")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Settings Saved:\n- %s, %s, %s", dayOfWeekInput, timeInput, timezoneInput),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		l.Error("failed to respond to interaction", slog.Any("err", err))
	} else {
		l.Info("successfully responded to interaction")
	}
}
