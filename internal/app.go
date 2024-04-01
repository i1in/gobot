package internal

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

type Application struct {
	session *discordgo.Session
}

var (
	BotToken       string
	Prefix         string
	Reg            = regexp.MustCompile(`^(.*?[.,-]{1})(.*)`)
	RegexGetUserId = regexp.MustCompile(`<@&?(\d+)>`)
	Args           []string
)

func New() *Application {
	session, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	session.Identify.Intents = discordgo.IntentsAll

	return &Application{
		session: session,
	}
}

func (a *Application) Start() {

	app := New()

	if err := app.session.Open(); err != nil {
		log.Fatal(err)
	}

	log.Println("Bot running")

	app.session.AddHandler(message)
	app.waitStop()
}

func (a *Application) waitStop() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	a.session.Close()
}

func findPrefix(s string) string {
	split := Reg.FindStringSubmatch(s)
	if len(split) < 2 {
		return ""
	}
	return split[1]
}

func Command(s string) string {
	split := Reg.FindStringSubmatch(s)
	if len(split) < 2 {
		return ""
	}
	return split[2]
}

func removeByIndex(array []string) []string {
	result := make([]string, 0)
	for _, str := range array {
		if str != "" && strings.TrimSpace(str) != "" {
			result = append(result, str)
		}
	}
	return result
}

func parseTime(hasUser bool, offset int) string {
	switch hasUser {
	case true:
		if len(Args) >= 2+offset {
			return Args[1+offset]
		}
	case false:
		return Args[1]
	}
	return "Навсегда"
}

func parseReason(hasUser bool, offset int) string {
	switch hasUser {
	case true:
		if len(Args) >= 3+offset {
			return strings.Join(Args[2+offset:], " ")
		}
	case false:
		if len(Args) == 2 && Args[1] != "" { // заменить Args[1] == "" на isInt boolean в случае если удалось сконвертировать значение, иначе
			return ""
		}

		if len(Args) >= 3 { // <- иначе
			return strings.Join(Args[2:], " ")
		}
	}
	return "Без причины"
}

func message(session *discordgo.Session, message *discordgo.MessageCreate) {

	if message.Author.ID == session.State.User.ID {
		return
	}

	Args = removeByIndex(strings.Split(message.Content, " "))

	log.Printf("[%s]", strings.Join(Args, ", "))

	if findPrefix(Args[0]) != Prefix {
		return
	}

	switch Command(Args[0]) {
	case "ban":
		command := &BanCommand{}
		userId := ""

		offset := 0

		if message.Message.ReferencedMessage != nil {
			userId = message.Message.ReferencedMessage.Author.ID
		}

		if userId != "" && len(Args) == 1 {
			log.Println(".ban sdfjbsdfhbdsjf sdkjf nsdkf s")
			return
		}

		hasUser := false
		if strings.HasPrefix(Args[1], "<@") {
			Args[1] = strings.Replace(Args[1], "&", "", 1)
			r := RegexGetUserId.FindStringSubmatch(Args[1])
			log.Println(r)
			if len(r) > 0 {
				userId = r[1]
			}
			hasUser = true
		}

		if userId == "" {
			// юзер не найден
			log.Println("юзер не найден")
			return
		}

		command.member = userId

		if hasUser {
			offset = 1
		}

		command.time = parseTime(hasUser, offset)
		command.reason = parseReason(hasUser, offset)

		log.Println(userId, command.time, command.reason)

		command.Run(session, message)
		return

	case "add":
		command := &VoiceChannel{}

		command.name = Args[1]

		command.Run(session, message)
		return
	}
}
