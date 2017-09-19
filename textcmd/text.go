package textcmd

// Returns Text commands

func Commands(command string) *string {
	message := ""

	switch command {
	case "!ann":
		message = "`BitcoinTalk ANN:` <https://bitcointalk.org/index.php?topic=1763606.0>\n"
	case "!commands":
		message = "Here are some custom Ubiq Slack commands\n"
		message += "\n"
		message += "*General*\n"
		message += "`!price [TICKER]` - Price look ups\n"
		message += "\n"
		message += "*Ubiq specific*\n"
		message += "`!ubqusd [AMOUNT]` - USD conversion\n"
		message += "`!urls` - URLs\n"
	case "!wallet":
		message = "`Web:` <https://pyrus.ubiqsmart.com> `Web2:` <https://myetherwallet.com> `GUI:` <https://github.com/ubiq/fusion/releases> `CLI:` <https://github.com/ubiq/go-ubiq/releases>\n"
	}

	return &message
}
