package bot

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
)

type HelpCommand struct {
  ID int
  Names []string
  Output []string
  Help string

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

func (c *HelpCommand) GetNames() []string {
  return c.Names
}

func (c *HelpCommand) GetOutput() []string {
  return c.Output
}

func (c *HelpCommand) GetID() int {
  return c.ID
}

func (c *HelpCommand) GetHelp() string {
  return c.Help
}
