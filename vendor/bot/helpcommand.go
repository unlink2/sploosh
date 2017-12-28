package bot

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
)

type HelpCommand struct {
  DefaultCommand
}

func (c *HelpCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate) bool {
  var out = "```Commands: \n"
  for _, op := range Commands {
    out = fmt.Sprintf("%s\n%s", out, op.GetHelp())
  }
  s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s\n\n%s", out, "Fork me: https://unlink2.github.io/sploosh/```"))
  return true
}

func (c *HelpCommand) OnGuildCreated(s *discordgo.Session, event *discordgo.GuildCreate) {
}
