package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type configStruct struct {
	Token     string
	BotPrefix string
	Squids    map[string]string
	Ask       string
	Err       string
}

type Command string

func (c Command) isAsk() bool {
	return c == "ask"
}

func (c Command) isSquidChoice() bool {
	i, err := strconv.ParseInt(string(c), 10, 32)
	return i > 0 && i < 10 && err == nil
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func ReadConfig() (*configStruct, error) {
	public := viper.New()
	public.SetConfigName("config")
	public.SetConfigType("json")
	public.AddConfigPath(".")
	public.ReadInConfig()

	private := viper.New()
	private.SetConfigName("config")
	private.SetConfigType("json")
	private.AddConfigPath("private")
	private.ReadInConfig()

	return &configStruct{
		Token:     private.GetString("bot_token"),
		BotPrefix: public.GetString("bot_prefix"),
		Squids:    public.GetStringMapString("squids"),
		Ask:       public.GetString("ask"),
		Err:       public.GetString("error"),
	}, nil
}

func main() {
	config, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	closeFunc, err := Start(config)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	defer close(c)
	go func() {
		for range c {
			closeFunc()
			os.Exit(0)
		}
	}()
	<-make(chan struct{})
}

func Start(config *configStruct) (func(), error) {
	storage, err := New("squid.bolt")

	if err != nil {
		return nil, err
	}

	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		return nil, err
	}

	u, err := goBot.User("@me")
	if err != nil {
		return nil, err
	}

	goBot.AddHandler(createMessageHandler(u.ID, config, storage))
	goBot.AddHandler(createReactionHandler(config, storage))

	return func() {
		goBot.Close()
		storage.Close()
	}, goBot.Open()

}

func createSender(s *discordgo.Session) func(channelId, content string) *discordgo.Message {
	return func(channelId, content string) *discordgo.Message {
		msg, err := s.ChannelMessageSend(channelId, content)
		logError(err)
		return msg
	}
}

func createReactionHandler(config *configStruct, storage Storage) func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		send := createSender(s)
		if ok, _ := storage.IsReactableMessage(r.MessageID); !ok {
			return
		}

		squid := Command(r.Emoji.Name[0:1])

		if squid.isSquidChoice() {
			storage.StoreSquidForuserId(r.UserID, string(squid))
		}
		send(r.ChannelID, fmt.Sprintf("@here\n%s FEELS THIS SQUIDWARD TODAY \n %s", r.Member.Mention(), config.Squids[string(squid)]))
	}
}

func createMessageHandler(BotId string, config *configStruct, storage Storage) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {

		send := createSender(s)
		if m.Author.ID == BotId || !strings.HasPrefix(m.Content, config.BotPrefix) {
			return
		}

		command := Command(strings.TrimPrefix(m.Content, config.BotPrefix))

		if command.isAsk() {
			msg := send(m.ChannelID, config.Ask)
			if msg == nil {
				return
			}
			logError(storage.StoreReactableMessage(msg.ID))
		}

		if command.isSquidChoice() {
			logError(storage.StoreSquidForuserId(m.Author.ID, string(command)))
			send(m.ChannelID, fmt.Sprintf("@here\n%s FEELS THIS SQUIDWARD TODAY \n %s", m.Author.Mention(), config.Squids[string(command)]))
		}

		if len(m.Mentions) > 0 {
			for _, user := range m.Mentions {
				if user.ID == BotId {
					continue
				}
				squid, err := storage.GetSquidForUserId(user.ID)
				if err != nil {
					logError(err)
					send(m.ChannelID, fmt.Sprintf("%s\n%s HAS NOT TOLD ME YET\n%s", m.Author.Mention(), user.Mention(), config.Squids["error"]))
				} else {
					send(m.ChannelID, fmt.Sprintf("%s\n%s FEELS THIS SQUIDWARD TODAY\n%s", m.Author.Mention(), user.Mention(), config.Squids[squid]))
				}
			}
		}
	}
}
