package bot

import (
  "github.com/bwmarrin/discordgo"
  "strings"
  "io/ioutil"
  "fmt"
  "config"
  "os"
)

type SoundCommand struct {
  DefaultCommand

}

func (c *SoundCommand) Execute(mw MessageWrapper) (bool, ResponseWrapper) {
  var res ResponseWrapper
  split := strings.Split(mw.Content, " ")

  // get whitelisted arrays
  whitelist := config.Globalcfg.Section("soundwhitelist")
  whitelistKeys := whitelist.Keys()

  var canExecute = false

  for _, key := range whitelistKeys {
    if key.String() == mw.DGuildID {
      canExecute = true
      break;
    }
  }

  if !canExecute && len(whitelistKeys) > 0 {
    res.Message += "Your discord guild is not whitelisted for sounds! Please contact" +
      "lukaskrickl@gmail.com to be added to the whitelist!"
    return false, res
  }

  if len(split) < 1 {
    return false, res
  }

  if split[0] == "~ls" {
    files, err := ioutil.ReadDir("./sounds")
    if err != nil {
      res.Message += "Unable to list sounds!"
      return false, res
    }

    out := "```Sounds:\n\n"

    for _, f := range files {
      out = out + f.Name() + "\n"
    }

    out = out + "```"

    res.Message += out
  } else if split[0] == "~ps" {
    if c.IsOnCooldown(mw.DAuthorID) {
      out := fmt.Sprintf("You are on cooldown for %d seconds!", c.GetRemainingCooldown(mw.DAuthorID))
      res.Message += out
      return false, res
    }

    if len(split) < 2 {
      res.Message += "Usage: ~ps <sound>. Use ~ls to see a list of all sound files."
      return false, res
    }

    if _, err := os.Stat("./sounds/" + split[1]); err != nil {
      res.Message += "Error playing sound: " + err.Error()
    } else {
      res.Sound = split[1]
      res.Message += "Playing sound " + split[1]
      c.SetCooldown(mw.DAuthorID, c.CooldownLen)
    }
  }

  return true, res
}

func (c *SoundCommand) OnGuildCreated(s *discordgo.Session, event *discordgo.GuildCreate) {
}
