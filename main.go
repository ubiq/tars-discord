package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
	"github.com/decred/dcrd/dcrutil"
	"github.com/joho/godotenv"
	"github.com/jpatel888/go-bitstamp"
	"github.com/jyap808/go-bittrex"
	"github.com/jyap808/go-gemini"
	"github.com/jyap808/go-poloniex"
	"github.com/ubiq/tars-discord/optionalchannelscmd"
	"github.com/ubiq/tars-discord/textcmd"
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

	appHomeDir        = dcrutil.AppDataDir("ubq-tars-discord", false)
	defaultBoldDBFile = filepath.Join(appHomeDir, "exchange_pairs.db")
	defaultConfigFile = filepath.Join(appHomeDir, "secrets.env")

	hClient = &http.Client{Timeout: 8 * time.Second}

	// Flood
	guildMemberAddCount  = 0
	floodMemberTimestamp time.Time
	floodAlertTimestamp  time.Time
)

const (
	// Flood
	floodMemberAddInterval = 3
	floodSeconds           = 60
	floodAlertSeconds      = 600
	floodAlertChannel      = "504427120236167207"
	moderatorID            = "348038402148532227"
	shieldsTimerSeconds    = 600

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

	poloniex := poloniex.New(poloniex_api_key, poloniex_api_secret)

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

	bittrex := bittrex.New(bittrex_api_key, bittrex_api_secret)

	rawticker := vals
	upperTicker := strings.ToUpper(*rawticker)
	tickerName := fmt.Sprintf("BTC-%s", upperTicker)

	marketSummary, err := bittrex.GetMarketSummary(tickerName)

	if err != nil {
		message = "Error: Trex invalid market"
	} else {
		y1 := marketSummary[0].PrevDay
		y2 := marketSummary[0].Last
		change := ((y2 - y1) / y1) * 100
		message = fmt.Sprintf("Bittrex  - BID: %.8f ASK: %.8f LAST: %.8f HIGH: %.8f LOW: %.8f VOLUME: %.2f %s, %.4f BTC CHANGE: %.2f%%",
			marketSummary[0].Bid, marketSummary[0].Ask, marketSummary[0].Last, marketSummary[0].High, marketSummary[0].Low, marketSummary[0].Volume, upperTicker, marketSummary[0].BaseVolume, change)
	}

	return &message
}

func getUSDto(fiat string) float64 {
	r, err := hClient.Get("https://api.fixer.io/latest?base=USD")
	if err != nil {
		log.Fatal(err)
		return -1.0
	}
	b, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
		return -1.0
	}
	defer r.Body.Close()
	var f interface{}
	json.Unmarshal([]byte(b), &f)
	return f.(map[string]interface{})["rates"].(map[string]interface{})[fiat].(float64)
}

func ubqEUR(amount *float64) *string {
	message := ""
	fiatErrMessage := "Error retrieving fiat conversion from remote API's"

	// Bittrex lookup
	bittrex := bittrex.New(bittrex_api_key, bittrex_api_secret)
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
	} else {
		eurValue := *amount * ticker.Ask * btcPrice
		message = fmt.Sprintf("```%.1f UBQ = €%.3f EUR```", *amount, eurValue)
		return &message
	}
}

func ubqUSD(amount *float64) *string {
	message := ""

	// Bittrex lookup
	bittrex := bittrex.New(bittrex_api_key, bittrex_api_secret)
	upperTicker := "UBQ"
	tickerName := fmt.Sprintf("BTC-%s", upperTicker)
	ticker, err := bittrex.GetTicker(tickerName)

	// BTC lookup
	gemini := gemini.New(gemini_api_key, gemini_api_secret)
	btcTickerName := "btcusd"
	btcTicker, err := gemini.GetTicker(btcTickerName)

	if err != nil {
		log.Println(err)
		message = "Error retrieving price from remote API's"
		return &message
	}

	btcPrice := btcTicker.Last

	usdValue := *amount * ticker.Ask * btcPrice

	message = fmt.Sprintf("```%.1f UBQ = $%.3f USD```", *amount, usdValue)

	return &message
}

func ubqLambo() *string {
	message := ""

	// Bittrex lookup
	bittrex := bittrex.New(bittrex_api_key, bittrex_api_secret)
	upperTicker := "UBQ"
	tickerName := fmt.Sprintf("BTC-%s", upperTicker)
	ticker, err := bittrex.GetTicker(tickerName)

	// BTC lookup
	gemini := gemini.New(gemini_api_key, gemini_api_secret)
	btcTickerName := "btcusd"
	btcTicker, err := gemini.GetTicker(btcTickerName)

	if err != nil {
		log.Println(err)
		message = "Error retrieving price from remote API's"
		return &message
	}

	btcPrice := btcTicker.Last
	amount := 1.0

	usdValue := amount * ticker.Ask * btcPrice

	lamboValue := 300000 / usdValue

	message = fmt.Sprintf("```You would need about %.0f UBQ to buy a lambo.```", lamboValue)

	return &message
}

func initializeBittrex(db *bolt.DB) (err error) {
	// Initial Bittrex table with data

	log.Println("initializeBittrex: START")
	bittrex := bittrex.New(bittrex_api_key, bittrex_api_secret)
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
		fmt.Println(err)
	}

	log.Println("initializeBittrex: END")
	return nil
}

func initializePoloniex(db *bolt.DB) (err error) {
	// Initialize Poloniex table with data

	log.Println("initializePoloniex: START")
	poloniex := poloniex.New(poloniex_api_key, poloniex_api_secret)
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

	gemini := gemini.New(gemini_api_key, gemini_api_secret)

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

func KeysString(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) *string {

	vals := &m.Content
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
			break
		}

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
			v := b.Get([]byte(ticker))
			if v != nil {
				bittrexUpdate = false
				tickerHeader = fmt.Sprintf("%s - %s", ticker, v)
				price := trexPrice(&ticker)
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
			v := b.Get([]byte(ticker))
			if v != nil {
				poloniexUpdate = false
				if tickerHeader == "" {
					tickerHeader = fmt.Sprintf("%s - %s", ticker, v)
				}
				price := poloPrice(&ticker)
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
		message = *ubqUSD(&amount)
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
		message = *ubqEUR(&amount)
	case "!ubqlambo":
		channelCheck := *checkTradingChannel(m.ChannelID)
		if channelCheck != "" {
			message = channelCheck
			break
		}

		message = *ubqLambo()
	// Text commands
	// Keep this in alphabetical order. Where possible just use the singular term.
	case "!ann", "!backup", "!blank", "!bots", "!caps", "!commands", "!compare", "!escher", "!escrow", "!ethunits", "!exchange", "!explorer", "!github", "!hide", "!hidechannels", "!invite", "!market", "!miner", "!mp", "!monetarypolicy", "!nucleus", "!onepage", "!pools", "!quarterly", "!resettabs", "!roadmap", "!site", "!social", "!solidity", "!stats", "!transparency", "!verified", "!verify", "!wallet", "!website", "!shokku", "!vyper":
		message = *textcmd.Commands(command)
	case "!join":
		usageStr := "**Usage:** !join [OPTIONAL_CHANNEL]\n\n"
		usageStr += fmt.Sprintf("**Optional Channels:** %s", KeysString(optionalChannels))

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
		usageStr += fmt.Sprintf("**Optional Channels:** %s", KeysString(optionalChannels))

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
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(guildMemberAdd)

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
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) > 0 {
		message := handleMessage(s, m)
		s.ChannelMessageSend(m.ChannelID, *message)
	}
}

// This function is called on GuildMemberAdd event
// Currently just performs Flood handling
func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {

	guildMemberAddCount += 1

	if guildMemberAddCount%floodMemberAddInterval != 0 {
		if floodMemberTimestamp.IsZero() {
			// Initialize value
			floodMemberTimestamp = time.Now()
		}
		return
	}

	// Flood check
	t1 := time.Now()
	diff := t1.Sub(floodMemberTimestamp)
	if diff.Seconds() < floodSeconds {
		if floodAlertTimestamp.IsZero() {
			// Initialize to past time
			floodAlertTimestamp = time.Now().Add(time.Second * time.Duration(-floodAlertSeconds*2))
		}

		diffAlert := t1.Sub(floodAlertTimestamp)

		if diffAlert.Seconds() > floodAlertSeconds {
			floodMessage := fmt.Sprintf("<@&%s> 🚨 Flood detected. %d Joins in %.1f seconds. Shields ON - Increased verification level", moderatorID, floodMemberAddInterval, diff.Seconds())
			s.ChannelMessageSend(floodAlertChannel, floodMessage)
			floodAlertTimestamp = t1

			// Shields on - Increase verification level
			level := discordgo.VerificationLevelVeryHigh
			gp := discordgo.GuildParams{"", "", &level, 0, "", 0, "", "", ""}
			s.GuildEdit(m.GuildID, gp)

			// Shields off - Restore default verification level after a timer
			shieldsOffMessage := fmt.Sprintf("Shields OFF - Restored default verification level")
			shieldsTimer := time.NewTimer(shieldsTimerSeconds * time.Second)

			go func() {
				<-shieldsTimer.C
				level = discordgo.VerificationLevelMedium
				gp = discordgo.GuildParams{"", "", &level, 0, "", 0, "", "", ""}
				s.GuildEdit(m.GuildID, gp)

				s.ChannelMessageSend(floodAlertChannel, shieldsOffMessage)
			}()

		}
	}
	// Set new value
	floodMemberTimestamp = time.Now()
}
