package bot

import (
  "github.com/bwmarrin/discordgo"
  "strings"
  "io/ioutil"
  "fmt"
)

type SoundCommand struct {
  DefaultCommand

}

func (c *SoundCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate) bool {
  split := strings.Split(m.Content, " ")

  if len(split) < 1 {
    return false
  }

  if split[0] == "~ls" {
    files, err := ioutil.ReadDir("./sounds")
    if err != nil {
      return false
    }

    out := "```Sounds:\n"

    for _, f := range files {
      out = out + f.Name() + "\n"
    }

    out = out + "```"

    s.ChannelMessageSend(m.ChannelID, out)
  } else if split[0] == "~ps" {
    if c.IsOnCooldown(m.Author.ID) {
      return false
    }

    if len(split) < 2 {
      return false
    }

    // load sound here
    sound, err := loadSound("./sounds/" + split[1])
    if err != nil {
      fmt.Println("Error loading sound:", err)
      return false
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
              err = playSound(s, guildSnd.ID, vs.ChannelID, sound)
        			if err != nil {
        				fmt.Println("Error playing sound:", err)
        			}
              break
      			}
          }
        }
      }
    }

    c.SetCooldown(m.Author.ID, c.DefaultCommand.CooldownLen)
  }

  return true
}

func (c *SoundCommand) OnGuildCreated(s *discordgo.Session, event *discordgo.GuildCreate) {
}
