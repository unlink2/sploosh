package bot

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
)

type SoundCommand struct {
  DefaultCommand

}

func (c *SoundCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate) bool {
  c.DefaultCommand.Execute(s, m)
  var out = "```Commands: \n"
  for _, op := range Commands {
    out = fmt.Sprintf("%s\n%s", out, op.GetHelp())
  }
  s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s\n\n%s", out, "Fork me: https://unlink2.github.io/sploosh/```"))
  return true
}

func (c *SoundCommand) OnGuildCreated(s *discordgo.Session, event *discordgo.GuildCreate) {
}
