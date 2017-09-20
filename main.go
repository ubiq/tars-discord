package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
	"github.com/decred/dcrutil"
	"github.com/joho/godotenv"
	"github.com/jyap808/go-bittrex"
	"github.com/jyap808/go-gemini"
	"github.com/jyap808/go-poloniex"
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
)

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

func handleMessage(vals *string) *string {

	valSplit := strings.Split(*vals, " ")
	message := ""

	if len(*vals) == 0 {
		return nil
	}

	command := valSplit[0]
	arguments := valSplit[1:]

	switch command {
	case "!price":
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
	// Text commands
	case "!ann", "!apx", "!commands", "!explorer", "!hide", "!hidechannels", "!invite", "!mine", "!miner", "!mining", "!pool", "!pools", "!site", "!verified", "!verify" , "!wallet", "!website":
		message = *textcmd.Commands(command)
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
		message := handleMessage(&m.Content)
		s.ChannelMessageSend(m.ChannelID, *message)
	}
}
