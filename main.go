package main

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/joho/godotenv"
	"github.com/jyap808/go-gemini"
	"github.com/ubiq/tars-discord/optionalchannelscmd"
	"github.com/ubiq/tars-discord/textcmd"
)

// Variables used for command line parameters
var (
	geminiAPIKey    = ""
	geminiAPISecret = ""

	appHomeDir        = dcrutil.AppDataDir("ubq-tars-discord", false)
	defaultConfigFile = filepath.Join(appHomeDir, "secrets.env")

	// Flood
	terminatorMemberFlag = false
)

const (
	// Flood
	floodMemberAddInterval = 3
	floodMilliseconds      = 60000 // 60 seconds
	floodAlertChannel      = "504427120236167207"
	moderatorRoleID        = "348038402148532227"
	terminatorTimerSeconds = 60

	tradingChannelID = "348036278673211392"
	turdRoleID       = "1072302088232636436"
	verifiedRoleID   = "350127755079319552"
)

func checkTradingChannel(channelID string) *string {
	message := ""

	if channelID != tradingChannelID {
		message = fmt.Sprintf("This price related command is only allowed in the <#%s> channel", tradingChannelID)
	}

	return &message
}

func btcPrice() *string {
	message := ""

	gemini := gemini.New(geminiAPIKey, geminiAPISecret)

	tickerName := "btcusd"

	ticker, err := gemini.GetTicker(tickerName)

	if err != nil {
		log.Println(err)
		message = "Error retrieving price from remote API's"
		return &message
	}

	lastPrice := ticker.Last
	message = fmt.Sprintf("```Gemini BTC price: %.2f```", lastPrice)

	return &message
}

func keysString(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) *string {

	vals := &m.Content
	asciiMatched := isASCII(*vals)
	httpMatched, _ := regexp.MatchString(`http`, *vals)
	uniSpamMatched, _ := regexp.MatchString(`[uU]n[iⅰ].*a[iⅰ]r[dԁ]rop`, *vals)
	axieInfinitySpamMatched, _ := regexp.MatchString(`[aA]x[iⅰ]e.*[iI]nf[iⅰ]n[iⅰ]ty`, *vals)
	if ((!asciiMatched && httpMatched) || uniSpamMatched || axieInfinitySpamMatched) && len(m.Member.Roles) == 0 {
		go terminateMember(s, m.GuildID, m.Author.ID, "Link spam")
		return nil
	}
	valSplit := strings.Split(*vals, " ")
	message := ""

	if len(*vals) == 0 {
		return nil
	}

	command := valSplit[0]
	arguments := valSplit[1:]

	var optionalChannels = map[string]bool{
		"decred-stakepool": true,
		"gamers":           true,
		"music":            true,
		"nsfw":             true,
		"sports":           true,
	}

	switch command {
	case "!price":
		channelCheck := *checkTradingChannel(m.ChannelID)
		if channelCheck != "" {
			message = channelCheck
			break
		}

		if len(arguments) == 0 {
			message = "Usage: !price [TICKER]"
			break
		}

		ticker := strings.ToUpper(arguments[0])

		// Special case to handle BTC price
		if ticker == "BTC" {
			message = *btcPrice()
		} else {
			message = "```Currently only supports BTC price lookups```"
		}
		// Delete originating message
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		message = fmt.Sprintf("%sRequested by: %s", message, m.Author.Mention())
	// Text commands
	// Keep this in alphabetical order. Where possible just use the singular term.
	case "!ann", "!backup", "!bots", "!bridge", "!caps", "!commands", "!compare", "!dojo", "!escher", "!escrow", "!ethunits", "!exchange", "!explorer", "!github", "!hide", "!hidechannels", "!invite", "!market", "!miner", "!mp", "!monetarypolicy", "!nft", "!nfts", "!nucleus", "!odin", "!onepage", "!pools", "!quarterly", "!redshift", "!roadmap", "!shinobi", "!site", "!social", "!solidity", "!stats", "!transparency", "!wallet", "!website", "!vyper":
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		message = fmt.Sprintf("%sRequested by: %s", *textcmd.Commands(command), m.Author.Mention())
	case "!join":
		usageStr := "**Usage:** !join [OPTIONAL_CHANNEL]\n\n"
		usageStr += fmt.Sprintf("**Optional Channels:** %s", keysString(optionalChannels))

		if len(arguments) == 0 {
			message = usageStr
			break
		}

		channel := strings.ToLower(arguments[0])
		channel = strings.TrimPrefix(channel, "#")

		if optionalChannels[channel] {
			message = *optionalchannelscmd.Join(s, m, channel)
		} else {
			message = usageStr
		}
	case "!leave":
		usageStr := "**Usage:** !leave [OPTIONAL_CHANNEL]\n\n"
		usageStr += fmt.Sprintf("**Optional Channels:** %s", keysString(optionalChannels))

		if len(arguments) == 0 {
			message = usageStr
			break
		}

		channel := strings.ToLower(arguments[0])
		channel = strings.TrimPrefix(channel, "#")

		if optionalChannels[channel] {
			message = *optionalchannelscmd.Leave(s, m, channel)
		} else {
			message = usageStr
		}
	case "!echo":
		if len(arguments) == 0 {
			message = "Usage: !echo [TEXT]"
		}

		valSplit2 := strings.SplitN(*vals, " ", 2)
		message = fmt.Sprintf("*Echo:* %s", valSplit2[1])
	default:
	}

	return &message
}

func main() {
	godotenv.Load(defaultConfigFile)
	token := os.Getenv("DISCORD_API_TOKEN")

	// Create the home directory if it doesn't already exist.
	err := os.MkdirAll(appHomeDir, 0700)
	if err != nil {
		// Show a nicer error message if it's because a symlink is
		// linked to a directory that does not exist (probably because
		// it's not mounted).
		if e, ok := err.(*os.PathError); ok && os.IsExist(err) {
			if link, lerr := os.Readlink(e.Path); lerr == nil {
				str := "is symlink %s -> %s mounted?"
				err = fmt.Errorf(str, e.Path, link)
			}
		}

		log.Fatalf("Failed to create config directory: %s", err)
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(guildMemberAdd)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) > 0 {
		message := handleMessage(s, m)
		if message != nil {
			s.ChannelMessageSend(m.ChannelID, *message)
		}
	}
}

var floodStack list.List

type floodCheck struct {
	userID  string
	addTime time.Time
}

func terminateMember(s *discordgo.Session, guildID string, userID string, reason string) {
	message := fmt.Sprintf("Terminated: <@%s>, Reason: %s", userID, reason)
	err := s.GuildBanCreateWithReason(guildID, userID, reason, 7)
	if err != nil {
		log.Printf("err: +%v\n", err)
	} else {
		s.ChannelMessageSend(floodAlertChannel, message)
	}
}

func turdifyMember(s *discordgo.Session, m *discordgo.GuildMemberAdd, reason string) {
	if len(m.Roles) == 0 {
		// Perform actions for users without roles
		message := fmt.Sprintf("Turdified: <@%s>, Reason: %s", m.User.ID, reason)
		s.GuildMemberRoleAdd(m.GuildID, m.User.ID, turdRoleID)
		s.ChannelMessageSend(floodAlertChannel, message)
	}
}

// Originally ported from https://github.com/hugonun/discordid2date/blob/master/main.js#L5
// Binary conversion and manipulation steps have been replaced with a bit shift operation
// to extract the relevant bits representing the timestamp.
func convertIDtoCreationTime(id string) time.Time {
	idInt, _ := strconv.ParseInt(id, 10, 64)
	unixCreationTime := (idInt >> 22) + 1420070400000
	creationTime := time.Unix(0, unixCreationTime*int64(time.Millisecond))
	return creationTime
}

// This function is called on GuildMemberAdd event
// Currently just performs Flood handling
func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Check and Terminate
	if terminatorMemberFlag {
		go terminateMember(s, m.GuildID, m.User.ID, "Flood join")
		return
	}

	if checkSpamName(s, m) {
		return
	}

	// Check new user Creation Time
	memberCreationTime := convertIDtoCreationTime(m.User.ID)
	t1 := time.Now()
	diff := t1.Sub(memberCreationTime)
	if diff.Hours() < 24*30 {
		turdMessage := fmt.Sprintf("Smelly turd: <@%s> account created %f hours ago", m.User.ID, diff.Hours())
		s.GuildMemberRoleAdd(m.GuildID, m.User.ID, turdRoleID)
		s.ChannelMessageSend(floodAlertChannel, turdMessage)
	}

	if floodStack.Len() < floodMemberAddInterval {
		floodStack.PushBack(floodCheck{m.User.ID, time.Now()})
	}
	if floodStack.Len() == floodMemberAddInterval {
		t1 := time.Now()
		floodStackFirst := floodStack.Front().Value.(floodCheck)
		diff := t1.Sub(floodStackFirst.addTime)
		if diff.Milliseconds() < floodMilliseconds {
			// Terminate all members in floodStack
			terminatorMemberFlag = true
			for member := floodStack.Front(); member != nil; member = member.Next() {
				go terminateMember(s, m.GuildID, member.Value.(floodCheck).userID, "Flood join")
			}

			// Set TerminatorTimer
			terminatorTimer := time.NewTicker(terminatorTimerSeconds * time.Second)
			go func() {
				floodMessage := fmt.Sprintf("<@&%s> 🚨 Flood detected. %d Joins in %.1f seconds", moderatorRoleID, floodMemberAddInterval, diff.Seconds())
				s.ChannelMessageSend(floodAlertChannel, floodMessage)
				<-terminatorTimer.C
				terminatorMemberFlag = false
				terminatorQuotes := []string{"You just can't go around killing people",
					"Your foster parents are dead",
					"Hasta la vista, baby",
					"My mission is to protect you",
					"Come with me if you want to live"}
				terminatorOffMessage := fmt.Sprintf("Terminator OFF - %s", terminatorQuotes[rand.Intn(len(terminatorQuotes))])
				s.ChannelMessageSend(floodAlertChannel, terminatorOffMessage)
			}()

			// Clear floodStack
			floodStack.Init()
		} else {
			floodStack.Remove(floodStack.Front())
			floodStack.PushBack(floodCheck{m.User.ID, time.Now()})
		}
	}

	// Run a 24-hour timer for the user, calling turdTimer after 24 hours
	go func() {
		time.Sleep(24 * time.Hour)
		turdifyMember(s, m, "Account Unverified after 24 hours")
	}()
}
