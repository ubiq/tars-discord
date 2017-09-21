package optionalchannelscmd

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var channels = map[string]string{
	"decred-stakepool": "360173511538638859",
	"music":            "359986790217678848",
	"nsfw":             "360239464884469760",
	"sports":           "360144092832989184",
}

var guildID = "274796263713538049"

func Join(s *discordgo.Session, m *discordgo.MessageCreate, command string) *string {
	message := ""

	s.GuildMemberRoleAdd(guildID, m.Author.ID, channels[command])
	message = fmt.Sprintf("%s added to %s channel\n", m.Author.Mention(), command)

	return &message
}

func Leave(s *discordgo.Session, m *discordgo.MessageCreate, command string) *string {
	message := ""

	s.GuildMemberRoleRemove(guildID, m.Author.ID, channels[command])
	message = fmt.Sprintf("%s removed from %s channel\n", m.Author.Mention(), command)

	return &message
}
