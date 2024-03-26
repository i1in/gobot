package internal

import (
	"Bot/internal/config"
	"github.com/bwmarrin/discordgo"
)

type User struct {
	UserId  string
	RoleIds []string
}

func (u *User) isAuthor(author string) bool {
	if u.UserId == author {
		return true
	}

	return false
}

func (u *User) GetRoles() []config.RolesConfig {
	perms := make([]config.RolesConfig, 0, len(config.C().Roles))

	for _, r := range config.C().Roles {
		for _, roleId := range u.RoleIds {
			if roleId == r.Id {
				perms = append(perms, r)
			}
		}
	}

	return perms
}

func (u *User) HasPerms(perm config.RolePermission) bool {
	perms := u.GetRoles()

	for _, r := range perms {
		for _, p := range r.Perm {
			if p == perm {
				return true
			}
		}
	}

	return false
}

func getUser(session *discordgo.Session, message *discordgo.MessageCreate) *User {
	member, err := session.GuildMember(message.GuildID, message.Author.ID)
	if err != nil {
		return nil
	}

	return &User{
		UserId:  message.Author.ID,
		RoleIds: member.Roles,
	}
}
