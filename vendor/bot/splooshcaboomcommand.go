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

func (c *SplooshKaboomCommand) Execute(mw MessageWrapper) (bool, ResponseWrapper) {
  var res ResponseWrapper
  // load sounds if not done already
  /*var err error
  if len(c.splooshSound) == 0 {
    c.splooshSound, err = loadSound("./sounds/sploosh.dca")
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Println("Loaded sploosh.dca!")
    }
  }
  if len(c.kaboomSound) == 0 {
    c.kaboomSound, err = loadSound("./sounds/kaboom.dca")
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Println("Loaded kaboom.dca!")
    }
  }*/

  // check for what the command started with
  if strings.HasPrefix(mw.Content, "~show") {
    res.Message += fmt.Sprintf("SPLOOSH! KABOOM!\n%s", c.RenderSplooshKaboom(mw.DGuildID, mw.Emoji))
  } else if strings.HasPrefix(mw.Content, "~reset") {
    e := GetEmojiForName("JKanStyle", mw.Emoji)
    res.Message += "Resetting " + EmojiToPrintableString(e, "")

    sk := c.GetSplooshKaboomForID(mw.DGuildID)
    sk.GenerateNewGame()
  } else if strings.HasPrefix(mw.Content, "~target") {
    split := strings.Split(mw.Content, " ")

    if len(split) < 3 {
      res.Message += "Usage: ~target x y"
      return false, res
    }
    x, err := strconv.Atoi(split[1])
    if err != nil {
      res.Message += "Usage: ~target x y"
      return false, res
    }

    y, err := strconv.Atoi(split[2])
    if err != nil {
      res.Message += "Usage: ~target x y"
      return false, res
    }
    if x - 1 < 0 || y - 1 < 0 {
      res.Message += "Usage: ~target x y"
      return false, res
    }

    response, result := c.Target(mw.DGuildID, y - 1, x - 1, mw.Emoji)

    if result < 0 {
      response = fmt.Sprintf("SPLOOSH! KABOOM!\n%s", response)
    } else if result == RESULTKABOOM || result == RESULTSHIPSUNK {
      response = fmt.Sprintf("KABOOM!\n%s", response)
    } else if result == RESULTSPLOOSH {
      response = fmt.Sprintf("SPLOOSH!\n%s", response)
    }

    res.Message += response
    if result == RESULTKABOOM {
      res.Sound = "kaboom.dca"
    } else if result == RESULTSPLOOSH {
      res.Sound = "sploosh.dca"
    }
  } else if strings.HasPrefix(mw.Content, "~cheat") {
    sk := c.GetSplooshKaboomForID(mw.DGuildID)
    sk.GameOver()

    res.Message += "Filthy Cheater!"
  }

  return true, res
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

func (c *SplooshKaboomCommand) RenderSplooshKaboom(gid string, emoji []*discordgo.Emoji) string {
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
