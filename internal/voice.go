package internal

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type VoiceChannel struct {
	name string
}

func (v *VoiceChannel) Run(session *discordgo.Session, message *discordgo.MessageCreate) {

	result := fmt.Sprintf("Участник <@%s> был добавлен в канал", v.name)
	session.ChannelMessageSend(message.ChannelID, result)
}
