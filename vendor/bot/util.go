package bot

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
  "encoding/binary"
	"io"
	"os"
  "time"
  "errors"
)

var Commands []CommandI

var MessagesToCleanBuffer []MessageBuffer

var soundPlayingInChannel []string

/*
MessageWrapper - This struct wraps around messages to allow different connection types
*/
type MessageWrapper struct {
  S *discordgo.Session
  M *discordgo.MessageCreate
  Content string

  DGuildID string
  DChannelID string
  DAuthorID string


  Guild *discordgo.Guild
  Channel *discordgo.Channel
  Emoji []*discordgo.Emoji
}

/*
ResponseWrapper - This struct represents a response from Execute
*/
type ResponseWrapper struct {
  Message string
  Sound string
}

type Cooldown struct {
  Userid string
  EndTime int64
}

type CommandI interface {
  Execute(mw MessageWrapper) (bool, ResponseWrapper)
  GetNames() []string
  GetOutput() []string
  GetID() int
  GetHelp() string

  IsOnCooldown(userid string) bool
  SetCooldown(userid string, duration int64)
  GetRemainingCooldown(userid string) int64

  OnGuildCreated(s *discordgo.Session, event *discordgo.GuildCreate)
}

type MessageBuffer struct {
  MessageID string
  ChannelID string
}

type DefaultCommand struct {
  ID int
  Names []string
  Output []string
  Help string

  CooldownLen int64
  cooldowns []Cooldown
}

func removeIndexMessage(s []MessageBuffer, index int) []MessageBuffer {
    return append(s[:index], s[index+1:]...)
}

func CleanMessages(s *discordgo.Session, channelid string) {
  for i, m := range MessagesToCleanBuffer {
    if m.ChannelID == channelid {
      s.ChannelMessageDelete(m.ChannelID, m.MessageID)
      MessagesToCleanBuffer = removeIndexMessage(MessagesToCleanBuffer, i)
    }
  }
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

func (c *DefaultCommand) Execute(mw MessageWrapper) (bool, ResponseWrapper) {
  var res ResponseWrapper
  if c.IsOnCooldown(mw.DAuthorID) {
    return false, res
  }

  for _, op := range c.Output {
    res.Message += op
  }

  c.SetCooldown(mw.DAuthorID, c.CooldownLen)

  return true, res
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

func (c *DefaultCommand) GetHelp() string {
  return c.Help
}

func (c *DefaultCommand) IsOnCooldown(userid string) bool {
  for index, cd := range c.cooldowns {
    if cd.Userid == userid && time.Now().Unix() < cd.EndTime {
      return true
    } else if cd.Userid == userid {
      c.cooldowns = append(c.cooldowns[:index], c.cooldowns[index+1:]...)
    }
  }
  return false
}

func (c *DefaultCommand) SetCooldown(userid string, duration int64) {
  c.cooldowns = append(c.cooldowns, Cooldown{Userid: userid, EndTime: time.Now().Unix() + duration})
}

func (c *DefaultCommand)  GetRemainingCooldown(userid string) int64 {
  for _, cd := range c.cooldowns {
    if cd.Userid == userid {
      return cd.EndTime - time.Now().Unix()
    }
  }

  return -1
}

// util stuff

func CreateCommands() {
  // create commands
  Commands = append(Commands, &SplooshKaboomCommand{DefaultCommand: DefaultCommand {
    ID: 0,
    Names: []string{"~reset", "~target", "~show", "~cheat"},
    Output: []string{},
    Help: "~reset -> resets game\n~target x y -> targets field\n~show -> shows current Sploosh Kaboom game",
  }})

  Commands = append(Commands, &HelpCommand{DefaultCommand: DefaultCommand {
    ID: 0,
    Names: []string{"~help"},
    Help: "~help -> prints help text",
    Output: []string{""},
  }})

  Commands = append(Commands, &SoundCommand{DefaultCommand: DefaultCommand {
    ID: 0,
    Names: []string{"~ps", "~ls"},
    Help: "~ps <sound name> -> plays sound\n~ls -> lists all sounds",
    Output: []string{},
    CooldownLen: 30,
  }})
}

func LoadSound(path string) ([][]byte, error) {
  var buffer = make([][]byte, 0)
  file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return nil, err
	}

	var opuslen int16

  fmt.Println("Reading dca file: ", path)

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

    if opuslen <= 0 {
      fmt.Println("Error reading from dca file : Unexpected opuslen")
      return buffer, errors.New("unexpected opuslen")
    }

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
  }
}

// playSound plays the current buffer to the provided channel.
func PlaySound(s *discordgo.Session, guildID string, channelID string, buffer [][]byte) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

  // check if channel is already joined
  for _, channel := range soundPlayingInChannel {
    if channel == channelID {
      return errors.New("Channel is already playing sound!")
    }
  }

  // add channel
  soundPlayingInChannel = append(soundPlayingInChannel, channelID)

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(500 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(500 * time.Millisecond)

  // remove index
  for index, channel := range soundPlayingInChannel {
    if channel == channelID {
      soundPlayingInChannel = append(soundPlayingInChannel[:index], soundPlayingInChannel[index+1:]...)
    }
  }

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
