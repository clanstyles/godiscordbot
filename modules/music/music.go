package music

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

type Music struct {
	vc *discordgo.VoiceConnection
}

var (
	channel = flag.String("music-channel-id", "", "this music channel id")
)

func New() *Music {
	if *channel == "" {
		flag.Usage()
		os.Exit(0)
	}

	return &Music{}
}

func (m *Music) Register(s *discordgo.Session) error {
	s.AddHandler(m.OnReady)
	s.AddHandler(m.OnDisconnect)
	s.AddHandler(m.OnMessage)

	return nil
}

func (m *Music) OnReady(s *discordgo.Session, event *discordgo.Ready) {
	c, err := s.Channel(*channel)
	if err != nil {
		log.Printf("failed to find channel: %s: %s", *channel, err)
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		log.Printf("failed to find guild: %s", err)
		return
	}

	voice, err := s.ChannelVoiceJoin(g.ID, c.ID, false, true)
	if err != nil {
		log.Printf("failed to join voice channel %s: %s", *channel, err)
		return
	}

	// save the voice channel for later
	m.vc = voice
}

func (m *Music) OnDisconnect(s *discordgo.Session, event *discordgo.Disconnect) {
	if err := m.vc.Disconnect(); err != nil {
		log.Printf("failed to leave voice channel %s: %s", *channel, err)
	}
}

func (m *Music) OnMessage(s *discordgo.Session, mc *discordgo.MessageCreate) {

	// ignore messages that aren't commands
	if !strings.HasPrefix(mc.Content, "!song") {
		return
	}

	// split the incoming message into two parts.
	// the second part should be a valid youtube link
	parts := strings.Split(mc.Content, " ")
	if len(parts) != 2 {
		return
	}

	if err := m.play(s, parts[1]); err != nil {
		log.Printf("failed to play song: %s", err)
	}
}

func (m *Music) play(s *discordgo.Session, video string) error {
	info, err := ytdl.GetVideoInfo(video)
	if err != nil {
		return err
	}

	log.Println("playing the song", info.Title)

	formats := info.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)
	if len(formats) == 0 {
		return errors.New("no supported video formats found")
	}

	download, err := info.GetDownloadURL(formats[0])
	if err != nil {
		return fmt.Errorf("failed to get youtube download url: %s", err)
	}

	// DCA encoding options
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 128
	options.Application = dca.AudioApplicationLowDelay

	session, err := dca.EncodeFile(download.String(), options)
	if err != nil {
		return fmt.Errorf("failed to encode mp3 file: %s", err)
	}
	defer session.Cleanup()

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking
	m.vc.Speaking(true)

	done := make(chan error)
	dca.NewStream(session, m.vc, done)

	// block until the conversion is complete
	err = <-done

	// check for an error
	if err != nil && err != io.EOF {
		return err
	}

	// Stop speaking
	m.vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	return nil
}
