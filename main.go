package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
	"github.com/decred/dcrutil"
	"github.com/joho/godotenv"
)

// Variables used for command line parameters
var (
	bittrex_api_key     = ""
	bittrex_api_secret  = ""
	gemini_api_key      = ""
	gemini_api_secret   = ""
	poloniex_api_key    = ""
	poloniex_api_secret = ""

	db *bolt.DB

	appHomeDir        = dcrutil.AppDataDir("bittrex-dolores", false)
	defaultBoldDBFile = filepath.Join(appHomeDir, "exchange_pairs.db")
	defaultConfigFile = filepath.Join(appHomeDir, "secrets.env")
)

func handleMessage(vals *string) *string {

	valSplit := strings.Split(*vals, " ")
	message := ""

	if len(*vals) == 0 {
		return nil
	}

	command := valSplit[0]
	arguments := valSplit[1:]

	switch command {
	case "!echo":
		if len(arguments) == 0 {
			message = "Usage: !echo [TEXT]"
		}

		valSplit2 := strings.SplitN(*vals, " ", 2)
		message = fmt.Sprintf("*Echo:* %s", valSplit2[1])
	case "!test1":
		line1 := "```*xyz* bbb"
		line2 := "*<http://ticker.com|123>* ccc```"
		message = fmt.Sprintf("%s\n%s", line1, line2)
	default:
	}

	return &message
}

func main() {
	godotenv.Load(defaultConfigFile)
	token := os.Getenv("DISCORD_API_TOKEN")
	bittrex_api_key = os.Getenv("BITTREX_API_KEY")
	bittrex_api_secret = os.Getenv("BITTREX_API_SECRET")

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

		log.Fatalln("Failed to create config directory")
	}

	// Optionally open database read only
	db, err = bolt.Open(defaultBoldDBFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) > 0 {
		message := handleMessage(&m.Content)
		s.ChannelMessageSend(m.ChannelID, *message)
	}
}
