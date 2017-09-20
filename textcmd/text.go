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
		message += "`!lennyface` - "( ͡° ͜ʖ ͡°)"
		message += "`!invite` - Discord Invite Link"
		message += "\n"
		message += "*Ubiq specific*\n"
		message += "`!ubqusd [AMOUNT]` - USD conversion\n"
		message += "`!urls` - URLs\n"
		message += "`!apx` - APX Ventures info"
		message += "`!miner / !miners` - Mining Software"
		message += "`!pool / !pools` - UBQ Pool List"
		message += "`!website / !site` - UBQ Website"
		message += "`!explorer` - UBQ Block Explorer"
	case "!wallet":
		message = "`Web:` <https://pyrus.ubiqsmart.com> `Web2:` <https://myetherwallet.com> `GUI:` <https://github.com/ubiq/fusion/releases> `CLI:` <https://github.com/ubiq/go-ubiq/releases>\n"
	case "!lennyface":
		message = "( ͡° ͜ʖ ͡°)"
	case "!apx":
		message = "`Channel:` #apx-ventures      `Website:` apxv.org     `Telegram:` t.me/apxventures     `Roadmap:` https://drive.google.com/file/d/0ByqyVzIU5PtFLXp2UGZPcUFYd1U/view"
	case "!invite":
		message = "https://discord.gg/HF6vEGF"
	case "!miner","!mine","!mining":
		message = "https://bitcointalk.org/index.php?topic=1433925.0"
	case "!pool","!pools":
		message = "http://ubiq.allcanmine.net (CN)   http://ubiqminer.com   http://ubiqmine.ca   https://ubq.kwikpool.party   https://ubiqpool.io http://pool.ubq.tw    http://www.ubiq.cc/minerpool    https://ubiq.suprnova.cc       http://ubiq.minerpool.net    http://ubq.poolcoin.biz     http://mole-pool.net    https://ubiq.coin-miners.info     https://aikapool.com/ubiq/      http://ubq.pool.sexy    https://ubq.poolto.be    http://ubq.minertopia.org https://ubiq.hakopool.com    http://ubiq.epicpool.club    https://ubq.zet-tech.eu   http://ubiq.hodlpool.com"
	case "!website","!site":
		message = "http://ubiqsmart.com/"
	case "!explorer":
		message = "`Explorer`: https://ubiqscan.io `Explorer 2`: http://www.ubiq.cc"		
	}

	return &message
}
