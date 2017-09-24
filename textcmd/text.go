package textcmd

var guestsChannelID = "348034489655623680"

func Commands(command string) *string {
	message := ""

	switch command {
	case "!ann":
		message = "`BitcoinTalk ANN:` <https://bitcointalk.org/index.php?topic=1763606.0>\n"
	case "!commands":
		message = "Here are some custom Ubiq Discord commands\n"
		message += "\n"
		message += "__**General**__\n\n"
		message += "`!price [TICKER]` - Price look ups\n"
		message += "`!invite` - Discord Invite Link\n"
		message += "`!hide / !hidechannels` - Tutorial on Hiding Discord Channels\n"
		message += "`!verify` - How to get Verified\n"
		message += "\n"
		message += "__**Ubiq specific**__\n\n"
		message += "`!apx` - APX Ventures info\n"
		message += "`!exchange / !market` - Ubiq exchanges\n"
		message += "`!explorer` - Ubiq Block Explorer\n"
		message += "`!hide` - How to Hide Channels\n"
		message += "`!miner` - Mining Software\n"
		message += "`!onepage` - Ubiq in one page"
		message += "`!pool` - Ubiq Pool List\n"
		message += "`!stats` - Ubiq network stats\n"
		message += "`!ubqusd [AMOUNT]` - USD conversion\n"
		message += "`!website / !site` - Ubiq Website\n"
	case "!apx":
		message = "`Channel:` #apx-ventures      `Website:` <http://apxv.org>     `Telegram:` <http://t.me/apxventures>     `Roadmap:` <https://drive.google.com/file/d/0ByqyVzIU5PtFLXp2UGZPcUFYd1U/view>\n"
	case "!exchange", "!market":
		message = "`Bittrex:` https://bittrex.com/Market/Index?MarketName=BTC-UBQ `Cryptopia:` https://www.cryptopia.co.nz/Exchange/?market=UBQ_BTC `Litebit:` https://www.litebit.eu/en/buy/ubiq\n"
	case "!explorer":
		message = "`Explorer:` <https://ubiqscan.io> `Explorer 2:` <http://www.ubiq.cc>\n"
	case "!hide", "!hidechannels":
		message = "<https://support.discordapp.com/hc/en-us/articles/213599277-How-do-I-hide-channels->\n"
	case "!invite":
		message = "`Ubiq Discord invite link:` <https://discord.gg/HF6vEGF>\n"
	case "!miner":
		message = "`Claymore:` <https://bitcointalk.org/index.php?topic=1433925.0>\n"
	case "!onepage":
		message = "`Ubiq in one page`: <https://medium.com/the-ubiq-report/ubiq-in-one-page-3e3d335064fc>\n"
	case "!pool":
		message = "`Ubiq mining pools:` http://ubiq.allcanmine.net (CN)   http://ubiqminer.com   http://ubiqmine.ca   https://ubq.kwikpool.party   https://ubiqpool.io http://pool.ubq.tw    http://www.ubiq.cc/minerpool    https://ubiq.suprnova.cc       http://ubiq.minerpool.net    http://ubq.poolcoin.biz     http://mole-pool.net    https://ubiq.coin-miners.info     https://aikapool.com/ubiq/      http://ubq.pool.sexy    https://ubq.poolto.be    http://ubq.minertopia.org https://ubiq.hakopool.com    http://ubiq.epicpool.club    https://ubq.zet-tech.eu   http://ubiq.hodlpool.com\n"
	case "!stats":
		message = "`Ubiq network stats:` <https://ubiq.darcr.us>\n"
	case "!verify", "!verified":
		message = "__**Verified**__"
		message += "\n"
		message += "You can request to be Verified in the <#348034489655623680> channel. This allows us to see who is a community member and gives you extra privileges such as Voice chat and access to more channels.\n\n"
		message += "To get Verified:\n"
		message += "a) Mention you would like to be Verified in the <#348034489655623680> channel. Mentioning where you came from and how you found out about Ubiq helps too.\n"
		message += "b) Only people with a profile pic are allowed to be Verified."
	case "!wallet":
		message = "`Web:` <https://pyrus.ubiqsmart.com> `Web2:` <https://myetherwallet.com> `GUI:` <https://github.com/ubiq/fusion/releases> `CLI:` <https://github.com/ubiq/go-ubiq/releases>\n"
	case "!website", "!site":
		message = "`Ubiq website:` <http://ubiqsmart.com/>\n"
	}

	return &message
}
