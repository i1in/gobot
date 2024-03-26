package internal

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

type UnbanCommand struct {
	name   string
	time   int
	reason []string
}

func (u *UnbanCommand) Run(session *discordgo.Session, message *discordgo.MessageCreate) {
	result := fmt.Sprintf("%s, %s, %s", u.name, strconv.Itoa(u.time), u.reason)
	session.ChannelMessageSend(message.ChannelID, result)
}
