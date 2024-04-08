package internal

import (
	"Bot/internal/config"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

type BanCommand struct {
	member string
	time   string
	reason string
}

func (b *BanCommand) Run(session *discordgo.Session, message *discordgo.MessageCreate) {

	user := getUser(session, message)

	if !user.HasPerms(config.PermBan) {
		session.ChannelMessageSendReply(message.ChannelID, ":x: Не хватает прав", message.Reference())
		return
	}

	if user.isAuthor(b.member) {
		session.ChannelMessageSendReply(message.ChannelID, ":x: Вы не можете забанить сам себя", message.Reference())
		return
	}

	outputTime := convert(b.time).convertTime()

	if _, err := strconv.Atoi(outputTime.time); err != nil {
		result := fmt.Sprintf("участник \n%s \nпричина\n%s", b.member, b.reason)
		session.ChannelMessageSend(message.ChannelID, result)
		return
	}

	result := fmt.Sprintf("%s, %s, %s", b.member, outputTime.timeToDate(), b.reason)

	session.ChannelMessageSend(message.ChannelID, result)
}
