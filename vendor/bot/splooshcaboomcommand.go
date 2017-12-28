package bot

import (
  "github.com/bwmarrin/discordgo"
  "strings"
  "fmt"
  "strconv"
)

type SplooshKaboomCommand struct {
  DefaultCommand

  sks []*SplooshKaboom

  splooshSound [][]byte
  kaboomSound [][]byte
}

func (c *SplooshKaboomCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate) bool {
  channel, _ := s.Channel(m.ChannelID)
  guild, _ := s.Guild(channel.GuildID)

  emoji := guild.Emojis

  // load sounds if not done already
  var err error
  if len(c.splooshSound) == 0 {
    c.splooshSound, err = loadSound("./sploosh.dca")
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Println("Loaded sploosh.dca!")
    }
  }
  if len(c.kaboomSound) == 0 {
    c.kaboomSound, err = loadSound("./kaboom.dca")
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Println("Loaded kaboom.dca!")
    }
  }

  // check for what the command started with
  if strings.HasPrefix(m.Content, "~show") {
    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("SPLOOSH! KABOOM!\n%s", c.RenderSplooshKaboom(guild.ID, emoji)))
  } else if strings.HasPrefix(m.Content, "~reset") {
    e := GetEmojiForName("JKanStyle", emoji)
    s.ChannelMessageSend(m.ChannelID, "Resetting " + EmojiToPrintableString(e, ""))

    sk := c.GetSplooshKaboomForID(guild.ID)
    sk.GenerateNewGame()
  } else if strings.HasPrefix(m.Content, "~target") {
    split := strings.Split(m.Content, " ")

    if len(split) < 3 {
      s.ChannelMessageSend(m.ChannelID, "Usage: ~target x y")
      return false
    }
    x, err := strconv.Atoi(split[1])
    if err != nil {
      s.ChannelMessageSend(m.ChannelID, "Usage: ~target x y")
      return false
    }

    y, err := strconv.Atoi(split[2])
    if err != nil {
      s.ChannelMessageSend(m.ChannelID, "Usage: ~target x y")
      return false
    }
    if x - 1 < 0 || y - 1 < 0 {
      s.ChannelMessageSend(m.ChannelID, "Usage: ~target x y")
      return false
    }

    response, result := c.Target(guild.ID, y - 1, x - 1, emoji)

    if result < 0 {
      response = fmt.Sprintf("SPLOOSH! KABOOM!\n%s", response)
    } else if result == RESULTKABOOM || result == RESULTSHIPSUNK {
      response = fmt.Sprintf("KABOOM!\n%s", response)
    } else if result == RESULTSPLOOSH {
      response = fmt.Sprintf("SPLOOSH!\n%s", response)
    }

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
            var buff [][]byte

            if result >= 0 {
              if result == RESULTKABOOM {
                buff = c.kaboomSound
              } else if result == RESULTSPLOOSH {
                buff = c.splooshSound
              }
              err = playSound(s, guildSnd.ID, vs.ChannelID, buff)
      				if err != nil {
      					fmt.Println("Error playing sound:", err)
      				}
            }
    			}
        }
      }
    }


    s.ChannelMessageSend(m.ChannelID, response)
  } else if strings.HasPrefix(m.Content, "~cheat") {
    sk := c.GetSplooshKaboomForID(guild.ID)
    sk.GameOver()

    s.ChannelMessageSend(m.ChannelID, "Filthy Cheater!")
  }


  return true
}

func (c *SplooshKaboomCommand) OnGuildCreated(s *discordgo.Session, event *discordgo.GuildCreate) {
  c.sks = append(c.sks, NewSplooshKaboom(event.Guild.ID))
}

/*
GetSplooshKaboomForID - returns a game for an id
gid - can either be a guild id or a channel
*/
func (c *SplooshKaboomCommand) GetSplooshKaboomForID(gid string) *SplooshKaboom {
  for _, sk := range c.sks {
    if sk.ID == gid {
      return sk
    }
  }

  // if it is not found that is an error! append new game and call again!
  c.sks = append(c.sks, NewSplooshKaboom(gid))
  return c.GetSplooshKaboomForID(gid)
}

func (c *SplooshKaboomCommand) Target(gid string, x int, y int, emoji []*discordgo.Emoji) (string, int) {
  sk := c.GetSplooshKaboomForID(gid)

  if sk.Bombs == 0 {
    return "Game Over!", -1
  }

  result := sk.Target(x, y)

  return c.RenderSplooshKaboom(gid, emoji), result
}

func (c *SplooshKaboomCommand) RenderSplooshKaboom(gid string,emoji []*discordgo.Emoji) string {
  sk := c.GetSplooshKaboomForID(gid)

  middleH := EmojiToPrintableString(GetEmojiForName("skMiddleH", emoji), "⬛")
  middleV := EmojiToPrintableString(GetEmojiForName("skMiddleV", emoji), "⬛")

  frontW := EmojiToPrintableString(GetEmojiForName("skFrontW", emoji), "◀")
  frontN := EmojiToPrintableString(GetEmojiForName("skFrontN", emoji), "🔼")
  frontS := EmojiToPrintableString(GetEmojiForName("skFrontS", emoji), "🔽")
  frontE := EmojiToPrintableString(GetEmojiForName("skFrontE", emoji), "▶")

  backW := EmojiToPrintableString(GetEmojiForName("skBackW", emoji), "▶")
  backN := EmojiToPrintableString(GetEmojiForName("skBackN", emoji), "🔽")
  backS := EmojiToPrintableString(GetEmojiForName("skBackS", emoji), "🔼")
  backE := EmojiToPrintableString(GetEmojiForName("skBackE", emoji), "◀")

  bomb := EmojiToPrintableString(GetEmojiForName("skBomb", emoji), "💣")
  bombUsed := EmojiToPrintableString(GetEmojiForName("skBombUsed", emoji), "🎱")

  notSunk := EmojiToPrintableString(GetEmojiForName("skNotSunk", emoji), "🦑")
  sunk := EmojiToPrintableString(GetEmojiForName("skSunk", emoji), "✅")
  blank := EmojiToPrintableString(GetEmojiForName("skBlank", emoji), "☁")

  sploosh := EmojiToPrintableString(GetEmojiForName("skSploosh", emoji), "❌")
  kaboom := EmojiToPrintableString(GetEmojiForName("skKaboom", emoji), "✅")

  var result = ""

  var bombsRendered = 0

  for x, _ := range sk.GameField {
    for i := 0; i < 3; i++ {
      if bombsRendered >= sk.Bombs {
        result = fmt.Sprintf("%s%s", result, bombUsed)
      } else {
        result = fmt.Sprintf("%s%s", result, bomb)
        bombsRendered++
      }
    }

    //result = fmt.Sprintf("%s %d ", result, x)

    for y, _ := range sk.GameField[x] {
      _, field := sk.IsFieldEmpty(x, y)

      if field == SKEMPTY {
        result = fmt.Sprintf("%s%s", result, blank)
      } else if field == SKFRONTE {
        result = fmt.Sprintf("%s%s", result, frontE)
      } else if field == SKFRONTN {
        result = fmt.Sprintf("%s%s", result, frontN)
      } else if field == SKFRONTS {
        result = fmt.Sprintf("%s%s", result, frontS)
      } else if field == SKFRONTW {
        result = fmt.Sprintf("%s%s", result, frontW)
      } else if field == SKBACKN {
        result = fmt.Sprintf("%s%s", result, backN)
      } else if field == SKBACKW {
        result = fmt.Sprintf("%s%s", result, backW)
      } else if field == SKBACKS {
        result = fmt.Sprintf("%s%s", result, backS)
      } else if field == SKBACKE {
        result = fmt.Sprintf("%s%s", result, backE)
      } else if field == SKMIDDLEH {
        result = fmt.Sprintf("%s%s", result, middleH)
      } else if field == SKMIDDLEV {
        result = fmt.Sprintf("%s%s", result, middleV)
      } else if field == SKSPLOOSH {
        result = fmt.Sprintf("%s%s", result, sploosh)
      } else if field == SKKABOOM {
        result = fmt.Sprintf("%s%s", result, kaboom)
      } else {
        result = fmt.Sprintf("%s%s", result, blank)
      }
    }

    if x < BOATAMOUNT {
      if sk.BoatsLeft - 1 < x {
        result = fmt.Sprintf("%s %s", result, sunk)
      } else {
        result = fmt.Sprintf("%s %s", result, notSunk)
      }
    }

    result += "\n"
  }
  if sk.BoatsLeft == 0 {
    return "You win!"
  }

  return fmt.Sprintf("%s", result)
}
