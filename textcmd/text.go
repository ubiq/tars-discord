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
		message += "__**General**__\n\n"
		message += "`!price [TICKER]` - Price look ups\n"
		message += "`!invite` - Discord Invite Link\n"
		message += "`!hide / !hidechannels` - Tutorial on Hiding Discord Channels\n"
		message += "`!verify` - How to get Verified"
		message += "\n"
		message += "__**Ubiq specific**__\n\n"
		message += "`!apx` - APX Ventures info\n"
		message += "`!explorer` - UBQ Block Explorer\n"
		message += "`!hide` - How to Hide Channels\n"
		message += "`!miner / !miners` - Mining Software\n"
		message += "`!pool / !pools` - UBQ Pool List\n"
		message += "`!ubqusd [AMOUNT]` - USD conversion\n"
		message += "`!website / !site` - UBQ Website\n"
	case "!apx":
		message = "`Channel:` #apx-ventures      `Website:` <http://apxv.org>     `Telegram:` <http://t.me/apxventures>     `Roadmap:` <https://drive.google.com/file/d/0ByqyVzIU5PtFLXp2UGZPcUFYd1U/view>\n"
	case "!explorer":
		message = "`Explorer`: <https://ubiqscan.io> `Explorer 2`: <http://www.ubiq.cc>\n"
	case "!hide", "!hidechannels":
		message = "<https://support.discordapp.com/hc/en-us/articles/213599277-How-do-I-hide-channels->\n"
	case "!invite":
		message = "<https://discord.gg/HF6vEGF>\n"
	case "!miner", "!mine", "!mining":
		message = "<https://bitcointalk.org/index.php?topic=1433925.0>\n"
	case "!pool", "!pools":
		message = "http://ubiq.allcanmine.net (CN)   http://ubiqminer.com   http://ubiqmine.ca   https://ubq.kwikpool.party   https://ubiqpool.io http://pool.ubq.tw    http://www.ubiq.cc/minerpool    https://ubiq.suprnova.cc       http://ubiq.minerpool.net    http://ubq.poolcoin.biz     http://mole-pool.net    https://ubiq.coin-miners.info     https://aikapool.com/ubiq/      http://ubq.pool.sexy    https://ubq.poolto.be    http://ubq.minertopia.org https://ubiq.hakopool.com    http://ubiq.epicpool.club    https://ubq.zet-tech.eu   http://ubiq.hodlpool.com\n"
	case "!verify","!verified":
		message = "__**Verified**__"
		message += "\n"
		message += "You can request to be #verified in the #general channel. this allows us to see who is a community member\n"
		message += "a) mention your user name from Slack\n"
		message += "b) only people with a profile pic get Verified\n"
		message += "c) once Verified, you are then added to the private #verified channel and other previously hidden channels"
	case "!wallet":
		message = "`Web:` <https://pyrus.ubiqsmart.com> `Web2:` <https://myetherwallet.com> `GUI:` <https://github.com/ubiq/fusion/releases> `CLI:` <https://github.com/ubiq/go-ubiq/releases>\n"
	case "!website", "!site":
		message = "<http://ubiqsmart.com/>\n"
	}

	return &message
}
