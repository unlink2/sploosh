package discord

import (
  "github.com/bwmarrin/discordgo"
  "github.com/go-ini/ini"
  "log"
  "fmt"
  "strings"
  "bot"
)

type Bot interface {

}

type DiscordConnection struct {
  Session *discordgo.Session
}

func NewDiscordBot(cfg *ini.File) *DiscordConnection {
  var newConn = new(DiscordConnection)
  var err error

  token, err := cfg.Section("discord").GetKey("token")
  if err != nil {
    log.Fatalln(err)
  }

  s, err := discordgo.New("Bot " + token.String())

  if err != nil {
    log.Fatalln(err)
  }

  s.AddHandler(newConn.onMessage)
  s.AddHandler(ready)
  s.AddHandler(guildCreate)

  // open connection
  err = s.Open()
  if err != nil {
    log.Fatalln(err)
  }

  fmt.Println("Bot is running! Press CTR+C to quit.")

  newConn.Session = s

  createCommands()

  return newConn
}

func createCommands() {
  // create commands
  bot.Commands = append(bot.Commands, &bot.SplooshKaboomCommand{DefaultCommand: bot.DefaultCommand {
    ID: 0,
    Names: []string{"~reset", "~target", "~show", "~cheat"},
    Output: []string{},
    Help: "~reset -> resets game\n~target x y -> targets field\n~show -> shows current Sploosh Kaboom game",
  }})

  bot.Commands = append(bot.Commands, &bot.HelpCommand{DefaultCommand: bot.DefaultCommand {
    ID: 0,
    Names: []string{"~help"},
    Help: "~help -> prints help text",
    Output: []string{""},
  }})

  bot.Commands = append(bot.Commands, &bot.SoundCommand{DefaultCommand: bot.DefaultCommand {
    ID: 0,
    Names: []string{"~ps", "~ls"},
    Help: "~ps <sound name> -> plays sound\n~ls -> lists all sounds",
    Output: []string{},
    CooldownLen: 30,
  }})
}

func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateStatus(0, "SPLOOSH! KABOOM!")

  fmt.Println("Bot ready!")
}


// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

  fmt.Printf("Joined guild %s\n", event.Guild.ID)

  for _, command := range bot.Commands {
    command.OnGuildCreated(s, event)
  }
}

func (*DiscordConnection) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
  // Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
    bot.CleanMessages(s, m.ChannelID)
    bot.MessagesToCleanBuffer = append(bot.MessagesToCleanBuffer, bot.MessageBuffer{
      MessageID: m.ID,
      ChannelID: m.ChannelID,
    })
		return
  }

  for _, command := range bot.Commands {
    for _, name := range command.GetNames() {
      if strings.HasPrefix(m.Content, name) {
        command.Execute(s, m)
      }
    }
  }
}
