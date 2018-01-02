package bot

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
)

type HelpCommand struct {
  DefaultCommand
}

func (c *HelpCommand) Execute(mw MessageWrapper) (bool, ResponseWrapper) {
  var res ResponseWrapper
  var out = "```Commands: \n"
  for _, op := range Commands {
    out = fmt.Sprintf("%s\n%s", out, op.GetHelp())
  }
  res.Message += fmt.Sprintf("%s\n\n%s", out, "Fork me: https://unlink2.github.io/sploosh/```")
  return true, res
}

func (c *HelpCommand) OnGuildCreated(s *discordgo.Session, event *discordgo.GuildCreate) {
}
