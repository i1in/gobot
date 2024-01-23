package internal

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	BotToken string
	Prefix   = "."
)

var r = regexp.MustCompile(`^(.*?[.,-]{1})(.*)`)

func main() {
	godotenv.Load()

	token := os.Getenv("BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err)
		return
	}

	dg.AddHandler(message)
	dg.Identify.Intents = discordgo.IntentsAll

	err = dg.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("bot running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func findPrefix(s string) string {
	split := r.FindStringSubmatch(s)
	if len(split) < 2 {
		return ""
	}
	return split[1]
}

func Command(s string) string {
	split := r.FindStringSubmatch(s)
	if len(split) < 2 {
		return ""
	}
	return split[2]
}

func convertTime(input string) (int, error) {
	unit := input[len(input)-1:]

	value, err := strconv.Atoi(strings.TrimSuffix(input, unit))
	if err != nil {
		return 0, fmt.Errorf("ошибка преобразования числа: %v", err)
	}

	switch unit {
	case "m":
		return value, nil
	case "h":
		return value * 60, nil
	case "d":
		return value * 24 * 60, nil
	default:
		return 0, fmt.Errorf("неизвестная единица измерения времени: %s", unit)
	}
}

func message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")
	if findPrefix(args[0]) != Prefix {
		return
	}

	if Command(args[0]) == "ban" {
		reg := regexp.MustCompile("[^0-9]")
		name := reg.ReplaceAllString(args[1], "")
		reason := append(args[:0], args[3:]...)
		banTime, err := convertTime(args[2])

		if err != nil {
			fmt.Println("Ошибка парсинга времени:", err)
			return
		}

		description := fmt.Sprintf(":white_check_mark: Участник **<@%s>** был забанен.", name)

		embed := discordgo.MessageEmbed{
			Description: description,
			Color:       0x004444,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Модератор: ",
					Value:  "<@" + m.Author.ID + ">",
					Inline: true,
				},
				{
					Name:   "Время бана: ",
					Value:  strconv.Itoa(int(banTime)) + "дн.",
					Inline: true,
				},
				{
					Name:   "Причина: ",
					Value:  strings.Join(reason, " "),
					Inline: false,
				},
			},
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &embed)

		s.GuildBanCreateWithReason(
			m.GuildID,
			name,
			strings.Join(reason, " "),
			0,
		)

		time.AfterFunc(time.Second*time.Duration(banTime*60), func() {
			err := s.GuildBanDelete(m.GuildID, name)
			if err != nil {
				fmt.Println("Error deleting ban:", err)
				return
			}
		})

		if err != nil {
			fmt.Println("Ошибка создания бана:", err)
			return
		}
	}

	if Command(args[0]) == "unban" {
		reg := regexp.MustCompile("[^0-9]")
		name := reg.ReplaceAllString(args[1], "")
		reason := append(args[:0], args[2:]...)

		err := s.GuildBanDelete(m.GuildID, name)
		if err != nil {
			result := fmt.Sprintf("Error deleting ban: %s", err)
			s.ChannelMessageSend(m.ChannelID, result)
			return
		} else {
			description := fmt.Sprintf(":white_check_mark: Участник **<@%s>** был разбанен.", name)

			embed := discordgo.MessageEmbed{
				Description: description,
				Color:       0x004444,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Модератор: ",
						Value:  "<@" + m.Author.ID + ">",
						Inline: true,
					},
					{
						Name:   "Причина: ",
						Value:  strings.Join(reason, " "),
						Inline: true,
					},
				},
			}

			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		}
	}

	if Command(args[0]) == "kick" {
		reg := regexp.MustCompile("[^0-9]")
		name := reg.ReplaceAllString(args[1], "")
		reason := append(args[:0], args[2:]...)

		err := s.GuildMemberDelete(m.GuildID, name)
		if err != nil {
			result := fmt.Sprintf("Error kicking: %s", err)
			s.ChannelMessageSend(m.ChannelID, result)
			return
		} else {

			description := fmt.Sprintf(":white_check_mark: Участник **<@%s>** был кикнут.", name)

			embed := discordgo.MessageEmbed{
				Description: description,
				Color:       0x004444,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Модератор: ",
						Value:  "<@" + m.Author.ID + ">",
						Inline: true,
					},
					{
						Name:   "Причина: ",
						Value:  strings.Join(reason, " "),
						Inline: true,
					},
				},
			}

			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		}
	}

	if Command(args[0]) == "ping" {
		minutes, err := convertTime(args[1])
		if err != nil {
			fmt.Printf("Ошибка конвертации времени для %s: %v\n", m.Content, err)
		} else {
			s.ChannelMessageSend(m.ChannelID, strconv.Itoa(minutes*60))
		}
	}

	if Command(args[0]) == "say" {
		s.ChannelMessageDelete(m.ChannelID, m.Message.ID)

		said := append(args[:0], args[1:]...)
		result := fmt.Sprintf("%s", strings.Join(said, " "))
		s.ChannelMessageSend(m.ChannelID, result)
	}

	if Command(args[0]) == "mute" {
		member := args[1]
		time := args[2]
		reason := append(args[:0], args[3:]...)
		result := fmt.Sprintf("Участник %s был замучен на %s.\nПричина: %s", member, time, strings.Join(reason, " "))
		s.ChannelMessageSend(m.ChannelID, result)
	}
}
