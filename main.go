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
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/joho/godotenv"
	"github.com/jpatel888/go-bitstamp"
	"github.com/jyap808/go-gemini"
	"github.com/jyap808/go-poloniex"
	"github.com/shopspring/decimal"
	"github.com/toorop/go-bittrex"
	"github.com/ubiq/tars-discord/optionalchannelscmd"
	"github.com/ubiq/tars-discord/textcmd"
	bolt "go.etcd.io/bbolt"
)

// Variables used for command line parameters
var (
	bittrexAPIKey     = ""
	bittrexAPISecret  = ""
	geminiAPIKey      = ""
	geminiAPISecret   = ""
	poloniexAPIKey    = ""
	poloniexAPISecret = ""

	db *bolt.DB

	appHomeDir        = dcrutil.AppDataDir("ubq-tars-discord", false)
	defaultBoldDBFile = filepath.Join(appHomeDir, "exchange_pairs.db")
	defaultConfigFile = filepath.Join(appHomeDir, "secrets.env")

	// Flood
	terminatorMemberFlag = false
)

const (
	// Flood
	floodMemberAddInterval = 3
	floodMilliseconds      = 60000 // 60 seconds
	floodAlertChannel      = "504427120236167207"
	moderatorID            = "348038402148532227"
	terminatorTimerSeconds = 60

	tradingChannelID = "348036278673211392"
)

func checkTradingChannel(channelID string) *string {
	message := ""

	if channelID != tradingChannelID {
		message = fmt.Sprintf("This price related command is only allowed in the <#%s> channel", tradingChannelID)
	}

	return &message
}

func poloPrice(vals *string) *string {
	message := ""

	poloniex := poloniex.New(poloniexAPIKey, poloniexAPISecret)

	rawticker := *vals
	upperTicker := strings.ToUpper(rawticker)
	tickerName := fmt.Sprintf("BTC_%s", upperTicker)

	tickers, err := poloniex.GetTickers()
	if err != nil {
		message = "Error: API not available"
		return &message
	}

	ticker, ok := tickers[tickerName]
	if ok {
		message = fmt.Sprintf("Poloniex - BID: %.8f ASK: %.8f LAST: %.8f HIGH: %.8f LOW: %.8f VOLUME: %.2f %s, %.4f BTC CHANGE: %.2f%%",
			ticker.HighestBid, ticker.LowestAsk, ticker.Last, ticker.High24Hr, ticker.Low24Hr, ticker.QuoteVolume, upperTicker, ticker.BaseVolume, (ticker.PercentChange * 100))
	} else {
		message = "Error: Polo Invalid market"
	}

	return &message
}

func trexPrice(vals *string) *string {
	message := ""

	bittrex := bittrex.New(bittrexAPIKey, bittrexAPISecret)

	rawticker := vals
	upperTicker := strings.ToUpper(*rawticker)
	tickerName := fmt.Sprintf("BTC-%s", upperTicker)

	marketSummary, err := bittrex.GetMarketSummary(tickerName)

	if err != nil {
		message = "Error: Trex invalid market"
	} else {
		y1 := marketSummary[0].PrevDay
		y2 := marketSummary[0].Last
		decimal100 := decimal.NewFromFloat(100.0)
		change := y2.Sub(y1).Div(y1).Mul(decimal100)
		message = fmt.Sprintf("Bittrex  - BID: %s ASK: %s LAST: %s HIGH: %s LOW: %s VOLUME: %s %s, %s BTC CHANGE: %s%%",
			marketSummary[0].Bid.StringFixed(8), marketSummary[0].Ask.StringFixed(8), marketSummary[0].Last.StringFixed(8), marketSummary[0].High.StringFixed(8), marketSummary[0].Low.StringFixed(8), marketSummary[0].Volume.StringFixed(2), upperTicker, marketSummary[0].BaseVolume.StringFixed(4), change.StringFixed(2))
	}

	return &message
}

func ubqEUR(amount *float64) *string {
	message := ""
	fiatErrMessage := "Error retrieving fiat conversion from remote API's"

	// Bittrex lookup
	bittrex := bittrex.New(bittrexAPIKey, bittrexAPISecret)
	upperTicker := "UBQ"
	tickerName := fmt.Sprintf("BTC-%s", upperTicker)
	ticker, err := bittrex.GetTicker(tickerName)

	// BTC lookup
	bitstamp := bitstamp.New()
	btcTickerName := "btceur"
	btcTicker, err := bitstamp.GetTicker(btcTickerName)

	if err != nil {
		log.Println(err)
		return &fiatErrMessage
	}

	btcPrice := btcTicker.Last

	if btcPrice == 0 {
		return &fiatErrMessage
	}
	decimalAmount := decimal.NewFromFloat(*amount)
	decimalBTCPrice := decimal.NewFromFloat(btcPrice)
	eurValue := ticker.Ask.Mul(decimalAmount).Mul(decimalBTCPrice)
	message = fmt.Sprintf("```%.1f UBQ = €%s EUR```", *amount, eurValue.StringFixed(3))
	return &message
}

func ubqUSD(amount *float64) *string {
	message := ""

	// Bittrex lookup
	bittrex := bittrex.New(bittrexAPIKey, bittrexAPISecret)
	upperTicker := "UBQ"
	tickerName := fmt.Sprintf("BTC-%s", upperTicker)
	ticker, err := bittrex.GetTicker(tickerName)

	// BTC lookup
	gemini := gemini.New(geminiAPIKey, geminiAPISecret)
	btcTickerName := "btcusd"
	btcTicker, err := gemini.GetTicker(btcTickerName)

	if err != nil {
		log.Println(err)
		message = "Error retrieving price from remote API's"
		return &message
	}

	btcPrice := btcTicker.Last

	decimalAmount := decimal.NewFromFloat(*amount)
	decimalBTCPrice := decimal.NewFromFloat(btcPrice)
	usdValue := ticker.Ask.Mul(decimalAmount).Mul(decimalBTCPrice)

	message = fmt.Sprintf("```%.1f UBQ = $%s USD```", *amount, usdValue.StringFixed(3))

	return &message
}

func ubqLambo() *string {
	message := ""

	// Bittrex lookup
	bittrex := bittrex.New(bittrexAPIKey, bittrexAPISecret)
	upperTicker := "UBQ"
	tickerName := fmt.Sprintf("BTC-%s", upperTicker)
	ticker, err := bittrex.GetTicker(tickerName)

	// BTC lookup
	gemini := gemini.New(geminiAPIKey, geminiAPISecret)
	btcTickerName := "btcusd"
	btcTicker, err := gemini.GetTicker(btcTickerName)

	if err != nil {
		log.Println(err)
		message = "Error retrieving price from remote API's"
		return &message
	}

	btcPrice := btcTicker.Last

	decimalAmount := decimal.NewFromFloat(1.0)
	decimalBTCPrice := decimal.NewFromFloat(btcPrice)
	usdValue := ticker.Ask.Mul(decimalAmount).Mul(decimalBTCPrice)

	lamboValue := decimal.NewFromFloat(300000.0).Div(usdValue)

	message = fmt.Sprintf("```You would need about %s UBQ to buy a lambo.```", lamboValue.StringFixed(0))

	return &message
}

func initializeBittrex(db *bolt.DB) (err error) {
	// Initial Bittrex table with data

	log.Println("initializeBittrex: START")
	bittrex := bittrex.New(bittrexAPIKey, bittrexAPISecret)
	markets, err := bittrex.GetMarkets()
	if err != nil {
		return err
	}

	bucketname := []byte("bittrex")

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		// Delete bucket
		err := tx.DeleteBucket(bucketname)
		if err != nil {
			log.Println("No bucket deleted")
		}

		// Open bucket
		bucket, err := tx.CreateBucket(bucketname)
		if err != nil {
			return err
		}

		for _, market := range markets {
			log.Println("initializeBittrex: ADD", market.MarketName)

			key := []byte(market.MarketCurrency)
			value := []byte(market.MarketCurrencyLong)
			err = bucket.Put(key, value)
			if err != nil {
				return err
			}
		}

		// Store the last updated time
		err = bucket.Put([]byte("lastupdated"), []byte(fmt.Sprint(time.Now().Unix())))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}

	log.Println("initializeBittrex: END")
	return nil
}

func initializePoloniex(db *bolt.DB) (err error) {
	// Initialize Poloniex table with data

	log.Println("initializePoloniex: START")
	poloniex := poloniex.New(poloniexAPIKey, poloniexAPISecret)
	currencies, err := poloniex.GetCurrencies()
	if err != nil {
		return err
	}

	bucketname := []byte("poloniex")

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		// Delete bucket
		err := tx.DeleteBucket(bucketname)
		if err != nil {
			log.Println("No bucket deleted")
		}

		// Open bucket
		bucket, err := tx.CreateBucket(bucketname)
		if err != nil {
			return err
		}

		for ticker, value := range currencies.Pair {
			if value.Delisted == 1 {
				continue
			}
			log.Println("initializePoloniex: ADD", ticker)
			key := []byte(ticker)
			value := []byte(value.Name)
			err = bucket.Put(key, value)
			if err != nil {
				return err
			}
		}

		// Store the last updated time
		err = bucket.Put([]byte("lastupdated"), []byte(fmt.Sprint(time.Now().Unix())))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}

	log.Println("initializePoloniex: END")
	return nil
}

func generatePriceMessage(prices []string, tickerHeader string) *string {
	message := ""
	beginString := "```"
	endString := "```"
	headerString := tickerHeader

	if len(prices) == 0 {
		message = "Ticker not found"
		return &message
	}

	pricesString := strings.Join(prices, "\n")

	message = fmt.Sprintf("%s%s\n%s%s", beginString, headerString, pricesString, endString)

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

	if err != nil {
		log.Println(err)
		message = "Error retrieving price from remote API's"
	} else {
		message = fmt.Sprintf("```Gemini BTC price: %.2f```", lastPrice)
	}

	return &message
}

func otherPrice(ticker *string) *string {
	message := ""

	var prices []string
	var tickerHeader string
	var bittrexUpdate bool
	var poloniexUpdate bool

	// Check Bittrex - Last updated
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("bittrex"))
		if b == nil {
			log.Println(".. you need to initialize the pairs database by running: price refresh")
		}
		updated := b.Get([]byte("lastupdated"))
		if updated != nil {
			ts, err := strconv.Atoi(string(updated))
			if err != nil {
				return err
			}
			t := time.Unix(int64(ts), 0)
			if time.Since(t).Seconds() > 3600 {
				bittrexUpdate = true
			}
		}
		return nil
	})

	// Check Bittrex
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("bittrex"))
		if b == nil {
			log.Println(".. you need to initialize the pairs database by running: price refresh")
		}
		v := b.Get([]byte(*ticker))
		if v != nil {
			bittrexUpdate = false
			tickerHeader = fmt.Sprintf("%s - %s", *ticker, v)
			price := trexPrice(ticker)
			prices = append(prices, *price)
		}
		return nil
	})

	// Check Poloniex - Last updated
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("poloniex"))
		if b == nil {
			log.Println(".. you need to initialize the pairs database by running: price refresh")
		}
		updated := b.Get([]byte("lastupdated"))
		if updated != nil {
			ts, err := strconv.Atoi(string(updated))
			if err != nil {
				return err
			}
			t := time.Unix(int64(ts), 0)
			if time.Since(t).Seconds() > 3600 {
				poloniexUpdate = true
			}
		}
		return nil
	})

	// Check Poloniex
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("poloniex"))
		if b == nil {
			log.Println(".. you need to initialize the pairs database by running: price refresh")
		}
		v := b.Get([]byte(*ticker))
		if v != nil {
			poloniexUpdate = false
			if tickerHeader == "" {
				tickerHeader = fmt.Sprintf("%s - %s", *ticker, v)
			}
			price := poloPrice(ticker)
			prices = append(prices, *price)
		}
		return nil
	})

	if bittrexUpdate {
		initializeBittrex(db)
	}
	if poloniexUpdate {
		initializePoloniex(db)
	}

	message = *generatePriceMessage(prices, tickerHeader)

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
	if !asciiMatched && httpMatched && len(m.Member.Roles) == 0 {
		go terminateMember(s, m.GuildID, m.Author.ID, "Generic spam")
		return nil
	}
	uniSpamMatched, _ := regexp.MatchString(`[uU]n[iⅰ].*a[iⅰ]r[dԁ]rop`, *vals)
	if uniSpamMatched && len(m.Member.Roles) == 0 {
		go terminateMember(s, m.GuildID, m.Author.ID, "Uniswap spam")
		return nil
	}
	axieInfinitySpamMatched, _ := regexp.MatchString(`[aA]x[iⅰ]e.*[iI]nf[iⅰ]n[iⅰ]ty`, *vals)
	if axieInfinitySpamMatched && len(m.Member.Roles) == 0 {
		go terminateMember(s, m.GuildID, m.Author.ID, "Axie Infinity spam")
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
		"china":            true,
		"decred-stakepool": true,
		"gamers":           true,
		"korea":            true,
		"music":            true,
		"nsfw":             true,
		"russia":           true,
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
			message = *otherPrice(&ticker)
		}
		// Delete originating message
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		message = fmt.Sprintf("%sRequested by: %s", message, m.Author.Mention())
	case "!pricerefresh":
		// Refresh coin pairs in Bolt DB
		initializeBittrex(db)
		initializePoloniex(db)
		message = "Exchange pairs refreshed"
	case "!ubqusd":
		channelCheck := *checkTradingChannel(m.ChannelID)
		if channelCheck != "" {
			message = channelCheck
			break
		}

		usageStr := "Usage: !ubqusd [AMOUNT] eg. !ubqusd 10"
		valueErrStr := fmt.Sprintf("Value error ;_; - %s", usageStr)
		if len(arguments) < 1 {
			message = usageStr
			break
		}

		amount, err := strconv.ParseFloat(arguments[0], 64)
		if err != nil {
			message = valueErrStr
			break
		}
		if amount < 0.1 || amount > 100000000 {
			message = "ERR: Pick an amount greater than 0.1 an less than 100 million"
			break
		}

		// Delete originating message
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		message = *ubqUSD(&amount)
		message = fmt.Sprintf("%sRequested by: %s", message, m.Author.Mention())
	case "!ubqeur":
		channelCheck := *checkTradingChannel(m.ChannelID)
		if channelCheck != "" {
			message = channelCheck
			break
		}

		usageStr := "Usage: !ubqeur [AMOUNT] eg. !ubqeur 10"
		valueErrStr := fmt.Sprintf("Value error ;_; - %s", usageStr)
		if len(arguments) < 1 {
			message = usageStr
			break
		}

		amount, err := strconv.ParseFloat(arguments[0], 64)
		if err != nil {
			message = valueErrStr
			break
		}
		if amount < 0.1 || amount > 100000000 {
			message = "ERR: Pick an amount greater than 0.1 an less than 100 million"
			break
		}

		// Delete originating message
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		message = *ubqEUR(&amount)
		message = fmt.Sprintf("%sRequested by: %s", message, m.Author.Mention())
	case "!ubqlambo":
		channelCheck := *checkTradingChannel(m.ChannelID)
		if channelCheck != "" {
			message = channelCheck
			break
		}

		message = *ubqLambo()
	// Text commands
	// Keep this in alphabetical order. Where possible just use the singular term.
	case "!ann", "!backup", "!blank", "!bots", "!bridge", "!caps", "!commands", "!compare", "!dojo", "!escher", "!escrow", "!ethunits", "!exchange", "!explorer", "!github", "!hide", "!hidechannels", "!invite", "!market", "!miner", "!mp", "!monetarypolicy", "!nucleus", "!onepage", "!pools", "!quarterly", "!redshift", "!resettabs", "!roadmap", "!shinobi", "!site", "!social", "!solidity", "!stats", "!transparency", "!wallet", "!website", "!vyper":
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
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
	banUserMessage := fmt.Sprintf("Terminated: <@%s>, Reason: %s", userID, reason)
	err := s.GuildBanCreateWithReason(guildID, userID, reason, 1)
	if err != nil {
		log.Printf("err: +%v\n", err)
	} else {
		s.ChannelMessageSend(floodAlertChannel, banUserMessage)
	}
}

// This function is called on GuildMemberAdd event
// Currently just performs Flood handling
func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Check and Terminate
	if terminatorMemberFlag {
		go terminateMember(s, m.GuildID, m.User.ID, "Flood join")
		return
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
				floodMessage := fmt.Sprintf("<@&%s> 🚨 Flood detected. %d Joins in %.1f seconds", moderatorID, floodMemberAddInterval, diff.Seconds())
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
}
