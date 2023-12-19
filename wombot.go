package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/mozoarella/wombot/internal/commands/batlibs"
	"github.com/mozoarella/wombot/internal/commands/magic8ball"
	"github.com/mozoarella/wombot/internal/commands/translations"
	"github.com/mozoarella/wombot/internal/types"
	"github.com/mozoarella/wombot/internal/utils"
)

type Randomizer = types.Randomizer

var (
	GuildID            = ""
	log                = utils.Log
	unregisterCommands = flag.Bool("unregister", false, "Use this flag to uunregister all registered bot commands")
)

var (
	botCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "batlibs",
			Description: "Replaces certain tokens with words",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "msg",
					Description: "Message with the tokens you want to batlib",
					Required:    true,
				},
			},
		},
		{
			Name:        "8ball",
			Description: "Shake a magic 8-ball",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "Optional question for the 8-ball",
					Required:    false,
				},
			},
		},
		{
			Name: "Translate to English",
			Type: discordgo.MessageApplicationCommand,
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"batlibs": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			//log.Info(fmt.Sprintf("User '%s' called the nounverb command with text '%s'", i.Member.User.Username, optionMap["msg"].StringValue()))
			libbedMsg := batlibs.DoBatlibs(optionMap["msg"].StringValue())

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: libbedMsg,
							Author: &discordgo.MessageEmbedAuthor{
								Name:    "Wombot Batlibs!",
								IconURL: s.State.User.AvatarURL("1024"),
							},
						},
					},
				},
			})
		},
		"8ball": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			//log.Info(fmt.Sprintf("User '%s' called the 8ball command", i.Member.User.Username))
			ballResponse := magic8ball.ShakeTheBall(i)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: ballResponse,
							Author: &discordgo.MessageEmbedAuthor{
								Name:    "Wombot Magic 8-Ball",
								IconURL: "https://wombot-files.mozoa.nl/icons/8-ball.png",
							},
						},
					},
				},
			})
		},
		"Translate to English": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: translations.MsgTranslationToEnglish(i),
							Author: &discordgo.MessageEmbedAuthor{
								Name:    "Wombot Translator",
								IconURL: s.State.User.AvatarURL("1024"),
							},
						},
					},
				},
			})
		},
	}
)

func init() {
	flag.Parse()
}

func main() {

	if os.Getenv("WB_DEBUG") == "true" {
		utils.LogLevel.Set(slog.LevelDebug)
	}

	log.Info("Starting Wombot")

	err := godotenv.Load()
	if err != nil {
		log.Debug("Couldn't load .env, this is probably fine")
	} else {
		log.Debug("Loaded .env file(s)")
	}

	/*testString := "I, the $ADJ $NOUN wish to $VERB many $ADJ $ADJ $NOUNS"
	for i := 0; i < 10; i++ {
		log.Info(commands.ProcessMessage((testString)))

	}*/

	dc, err := discordgo.New("Bot " + os.Getenv("WB_DC_TOKEN"))
	if err != nil {
		log.Error("Couldn't set up the Discord session", err)
		return
	}

	dc.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	dc.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentsDirectMessages

	err = dc.Open()
	if err != nil {
		log.Error("Error opening connection", err)
		return
	}

	log.Info(fmt.Sprintf("Wombot is running with the username '%s' and ID '%s'", dc.State.User.Username, dc.State.User.ID))

	log.Info("Registering commands with the Discord API")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(botCommands))
	for i, v := range botCommands {
		cmd, err := dc.ApplicationCommandCreate(dc.State.User.ID, "", v)
		if err != nil {
			log.Error(fmt.Sprintf("Cannot create '%s' command: %s", v.Name, err))
		}
		registeredCommands[i] = cmd
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if *unregisterCommands {
		log.Info("Unregistering commands...")
		commandsToRemove, err := dc.ApplicationCommands(dc.State.User.ID, "")
		if err != nil {
			log.Error("Could not get registered commands")
		}
		for _, c := range commandsToRemove {
			err := dc.ApplicationCommandDelete(dc.State.User.ID, "", c.ID)
			if err != nil {
				log.Error(fmt.Sprintf("Cannot delete '%s' command: %s", c.Name, err))
			}
		}
	}

	log.Info("Gracefully shutting down")
	err = dc.Close()
	if err != nil {
		log.Error(err.Error())
	}
}