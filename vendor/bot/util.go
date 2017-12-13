package bot

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
)

var Commands []CommandI

type CommandI interface {
  Execute(s *discordgo.Session, m *discordgo.MessageCreate) bool
  GetNames() []string
  GetOutput() []string
  GetID() int

  OnGuildCreated(s *discordgo.Session, event *discordgo.GuildCreate)
}

type DefaultCommand struct {
  ID int
  Names []string
  Output []string

}

func EmojiToPrintableString(emoji *discordgo.Emoji, fallback string) string {
  if emoji == nil {
    return fallback
  }

  return fmt.Sprintf("<:%s:%s>", emoji.Name, emoji.ID)
}

func GetEmojiForName(name string ,emoji []*discordgo.Emoji) *discordgo.Emoji {
  for _, e := range emoji {
    if e.Name == name {
      return e
    }
  }

  return nil
}

func (c *DefaultCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate) bool {
  for _, op := range c.Output {
    s.ChannelMessageSend(m.ChannelID, op)
  }

  return true
}

func (c *DefaultCommand) OnGuildCreated(s *discordgo.Session, event *discordgo.GuildCreate) {
}

func (c *DefaultCommand) GetNames() []string {
  return c.Names
}

func (c *DefaultCommand) GetOutput() []string {
  return c.Output
}

func (c *DefaultCommand) GetID() int {
  return c.ID
}
