package internal

import (
	"fmt"
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

func message(session *discordgo.Session, message *discordgo.MessageCreate) {

	if message.Author.ID == session.State.User.ID {
		return
	}

	args := removeByIndex(strings.Split(message.Content, " "))

	log.Printf("[%s]", strings.Join(args, ", "))

	if findPrefix(args[0]) != Prefix {
		return
	}

	switch Command(args[0]) {
	case "ban":
		command := &BanCommand{}
		userId := ""
		isInt := false
		offset := 0

		if message.Message.ReferencedMessage != nil {
			userId = message.Message.ReferencedMessage.Author.ID
		}

		if userId != "" && len(args) == 1 {
			log.Println(".ban sdfjbsdfhbdsjf sdkjf nsdkf s")
			return
		}

		hasUser := false
		if len(args) > 1 && strings.HasPrefix(args[1], "<@") {
			args[1] = strings.Replace(args[1], "&", "", 1)
			r := RegexGetUserId.FindStringSubmatch(args[1])
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

		fmt.Println("%s", hasUser)

		parser := getArgs(args, hasUser, offset)
		command.time, isInt = parser.parseTime()
		command.reason = parser.parseReason(isInt)

		log.Println(userId, command.time, command.reason)

		command.Run(session, message)
		return

	case "add":
		command := &VoiceChannel{}

		command.name = args[1]

		command.Run(session, message)
		return
	}
}
