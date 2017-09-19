package textcmd

// Returns Text commands

func Commands(command string) *string {
	message := ""

	switch command {
	case "!ann":
		message = "`BitcoinTalk ANN:` https://bitcointalk.org/index.php?topic=1763606.0\n"
	case "!commands":
		message = "Here are some custom Ubiq Slack commands\n"
		message += "\n"
		message += "*General*\n"
		message += "`!lennyface` - Lenny face\n"
		message += "`!price [TICKER]` - Price look ups\n"
		message += "\n"
		message += "*Ubiq specific*\n"
		message += "`!ubqusd [AMOUNT]` - USD conversion\n"
		message += "`!urls` - URLs\n"
	}

	return &message
}
