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

  return newConn
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
        channel, _ := s.Channel(m.ChannelID)
        guild, _ := s.Guild(channel.GuildID)

        emoji := guild.Emojis

        _, res := command.Execute(bot.MessageWrapper{S: s, M: m,
          Content: m.Content,
          Emoji: emoji,
          Guild: guild,
          Channel: channel,

          DChannelID: channel.ID,
          DGuildID: guild.ID,
          DAuthorID: m.Author.ID,
        })

        s.ChannelMessageSend(m.ChannelID, res.Message)

        if res.Sound != "" {
          // load sound here
          sound, err := bot.LoadSound("./sounds/" + res.Sound)
          if err != nil && sound == nil {
            s.ChannelMessageSend(m.ChannelID, "Error loading sound:" + err.Error())
            fmt.Println("Error loading sound:", err)
            continue
          } else {
            // Find the channel that the message came from.
        		channel, err := s.State.Channel(m.ChannelID)
        		if err != nil {
        			// Could not find channel.
              fmt.Println("Error finding channel:", err)
            } else {
              // Find the guild for that channel.
          		guildSnd, err := s.State.Guild(channel.GuildID)
          		if err != nil {
          			// Could not find guild.
          			fmt.Println("Error finding guild:", err)
              } else {
                // Look for the message sender in that guild's current voice states.
            		for _, vs := range guildSnd.VoiceStates {

            			if vs.UserID == m.Author.ID {
                    //c.SetCooldown(m.Author.ID, c.DefaultCommand.CooldownLen)
                    err = bot.PlaySound(s, guildSnd.ID, vs.ChannelID, sound)
              			if err != nil {
                      s.ChannelMessageSend(m.ChannelID, "Error playing sound:" + err.Error())
              				fmt.Println("Error playing sound:", err)
              			}
                    break
            			}
                }
              }
            }
          }
        }
      }
    }
  }
}
