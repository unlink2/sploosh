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

type Cooldown struct {
  Userid string
  EndTime int64
}

type CommandI interface {
  Execute(s *discordgo.Session, m *discordgo.MessageCreate) bool
  GetNames() []string
  GetOutput() []string
  GetID() int
  GetHelp() string

  IsOnCooldown(userid string) bool
  SetCooldown(userid string, duration int64)

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

func (c *DefaultCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate) bool {
  if c.IsOnCooldown(m.Author.ID) {
    return false
  }

  for _, op := range c.Output {
    s.ChannelMessageSend(m.ChannelID, op)
  }

  c.SetCooldown(m.Author.ID, c.CooldownLen)

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

// util stuff

func loadSound(path string) ([][]byte, error) {
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
func playSound(s *discordgo.Session, guildID string, channelID string, buffer [][]byte) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

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

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
