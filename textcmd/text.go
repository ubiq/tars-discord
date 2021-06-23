package textcmd

func Commands(command string) *string {
	message := ""

	switch command {
	case "!ann":
		message = "`BitcoinTalk ANN:` <https://bitcointalk.org/index.php?topic=1763606.0>\n"
	case "!commands":
		message = "Here are some custom Ubiq Discord commands\n"
		message += "\n"
		message += "__**General**__\n\n"
		message += "`!escrow` - Escrow caution message\n"
		message += "`!hide / !hidechannels` - Tutorial on Hiding Discord Channels\n"
		message += "`!invite` - Ubiq Discord Invite Link\n"
		message += "`!price [TICKER]` - Price look ups\n"
		message += "\n"
		message += "__**Ubiq specific**__\n\n"
		message += "`!backup` - Backup your account\n"
		message += "`!bridge / !redshift` - Ubiq Redshift bridge information\n"
		message += "`!bots` - Ubiq Twitter bots\n"
		message += "`!caps` - Correct spelling for Ubiq\n"
		message += "`!compare` - Comparison chart\n"
		message += "`!dojo` - Guides to using Shinobi\n"
		message += "`!escher` - Escher contract info\n"
		message += "`!exchange / !market` - Ubiq exchanges\n"
		message += "`!explorer` - Ubiq Block Explorer\n"
		message += "`!github` - Ubiq GitHub Repository\n"
		message += "`!hide` - How to Hide Channels\n"
		message += "`!miner` - Mining Software\n"
		message += "`!mp / !monetarypolicy` - Monetary Policy\n"
		message += "`!nucleus / !transparency` - Nucleus Transparency Report\n"
		message += "`!onepage` - Ubiq in one page\n"
		message += "`!pools` - Ubiq Pools List\n"
		message += "`!resettabs / !blank` - Reset tabs in Fusion\n"
		message += "`!roadmap / !quarterly` - Roadmap and Quarterly report\n"
		message += "`!shinobi` - Access link and Info page for Shinobi\n"
		message += "`!social` - Ubiq social media links\n"
		message += "`!stats` - Ubiq network stats\n"
		message += "`!ubqeur [AMOUNT]` - EUR conversion\n"
		message += "`!ubqusd [AMOUNT]` - USD conversion\n"
		message += "`!website / !site` - Ubiq Website\n"
		message += "\n"
		message += "__**Programming**__\n\n"
		message += "`!ethunits` - Ethereum units\n"
		message += "`!solidity` - Solidity documentation\n"
	case "!backup":
		message = "To backup your wallet file(s), save every file in the following directories: Mac: `~/Library/Ubiq/keystore` Linux: `~/.ubiq/keystore` Windows: `%APPDATA%/Ubiq/keystore` *Note that each time you create a new account a new file will be created, so you must back up the new file in that directory when you create a new account.* The following video will help you on Windows: https://www.youtube.com/watch?v=x5tNtKpnkMw\n"
	case "!bots":
		message = "`Ubiq new coin bot:` <https://twitter.com/ubiqcoinbot> `Ubiq BCT ANN bot:` <https://twitter.com/ubiqannbot>\n"
	case "!bridge", "!redshift":
		message = "`Wrapped Ubiq (WUBQ) Contract:` <https://ubiqscan.io/address/0x1fa6a37c64804c0d797ba6bc1955e50068fbf362>\n"
		message += "`Polygon WUBQ Contract:` <https://polygonscan.com/token/0xb1c5c9b97b35592777091cd34ffff141ae866abd>\n"
		message += "`Quickswap WETH-WUBQ Pair:` <https://info.quickswap.exchange/pair/0x8dccc1251f93fc4eb71dc7a686b58d3e718a5949>\n"
	case "!caps":
		message = "The correct spelling for Ubiq is 'Ubiq', not all capitals (UBIQ is incorrect). This is similar to Nike where the logo is in all capitals but not the name. The ticker symbol for Ubiq is UBQ.\n"
	case "!compare":
		message = "`Comparison chart provided by user moreexplosions:` https://imgur.com/a/Kr8RW"
	case "!dojo":
		message = "`Introduction and making trades:` <https://blog.ubiqsmart.com/the-ubiq-dex-introducing-shinobi-5433adecc5e3>\n"
		message += "`Using Sparrow:` <https://blog.ubiqsmart.com/the-ubiq-dex-using-sparrow-96bb604a1c89>\n"
		message += "`Liquidity Pools:` <https://blog.ubiqsmart.com/the-ubiq-dex-liquidity-pools-6b1b1982c30a>\n"
		message += "`Using Shinobi - Tips and Tricks to keep you safe:` <https://blog.ubiqsmart.com/using-shinobi-tips-and-tricks-to-keep-you-safe-8116e38b53b5>\n"
		message += "`Using Shinobi - Charts and General Info:` <https://blog.ubiqsmart.com/using-shinobi-charts-and-general-info-f449d6c326ec>\n"
		message += "`Using Shinobi - A Guide on Transaction Nonces:` <https://blog.ubiqsmart.com/using-shinobi-a-guide-on-transaction-nonces-6c8c058c5512>\n"
		message += "`Using Shinobi - Yield Farming and Token Generation Events:` <https://blog.ubiqsmart.com/using-shinobi-yield-farming-and-token-generation-events-e83a48de5824>\n"
	case "!escher":
		message = "`Escher contract address:` 0xcf3222b7FDa7a7563b9E1E6C966Bead04AC23c36  - Use default ABI and default 18 decimals."
	case "!escrow":
		message = "`Caution:` When trading OTC directly between yourself and another user, there's a chance you may get scammed by the other party and lose your funds. Whenever making an OTC trade, _always_ use an escrow. A trusted escrow protects both parties during the trade. Contact Ubiq Discord moderators to arrange a trustworthy escrow for you."
	case "!ethunits":
		message = "`Ethereum units:` <https://github.com/ryepdx/ethereum-units>\n"
	case "!exchange", "!market":
		message = "`Bittrex:` <https://bittrex.com/Market/Index?MarketName=BTC-UBQ>\n"
		message += "`QuickSwap:` <https://info.quickswap.exchange/pair/0x8dccc1251f93fc4eb71dc7a686b58d3e718a5949>\n"
		message += "`Dove Wallet:` <https://dovewallet.com/en/trade/spot/ubq-usdt>\n"
		message += "`BitZ:` <https://www.bitz.cm/en/exchange/ubq_usdt>\n"
		message += "`Catalx:` <https://catalx.io/trade/BTC-UBQ>\n"
		message += "`Asymetrex:` <https://asymetrex.com/markets/ubqbtc>\n"
	case "!explorer":
		message = "`Explorer:` <https://ubiqscan.io> `Explorer 2:` <https://ubiqexplorer.com>\n"
	case "!github":
		message = "`GitHub:` <https://github.com/ubiq>\n"
	case "!hide", "!hidechannels":
		message = "<https://support.discordapp.com/hc/en-us/articles/213599277-How-do-I-hide-channels->\n"
	case "!invite":
		message = "`Ubiq Discord invite link:` <https://discord.gg/XaqzJB4>\n"
	case "!miner":
		message = "`Ubqminer (0% dev fee):` <https://github.com/ubiq/ubqminer/releases> `PhoenixMiner (0.65% dev fee):` <https://bitcointalk.org/index.php?topic=2647654.0> `TT-Miner (1% dev fee):` <https://bitcointalk.org/index.php?topic=5025783.0> `Nanominer (1% dev fee):` <https://nanominer.org/> `SRBMiner (0.65% dev fee):` <https://www.srbminer.com>\n"
	case "!mp", "!monetarypolicy":
		message = "`Monetary policy and mining block rewards scheme in Ubiq:` <https://blog.ubiqsmart.com/ubiq-research-monetary-policy-2e27458983ec>\n"
	case "!nucleus", "!transparency":
		message = "`Nucleus Transparency Report:` <https://blog.ubiqsmart.com/nucleus-transparency-report-6496e444bd85>\n"
	case "!onepage":
		message = "`Ubiq in one page`: <https://blog.ubiqsmart.com/ubiq-in-one-page-df1672fb85dd>\n"
	case "!pools":
		message = "`List of known mining pools:`\n"
		message += "<https://gomine.pro/pool/ubiq>\n"
		message += "<https://ubq.mypool.online>\n"
		message += "<https://comining.io> <https://ubiqpool.maxhash.org>\n"
		message += "<https://ubiqpool.io>		<https://ubiq.clona.ru>\n"
		message += "<https://ubq.solopool.org>     <https://ubiq.wattpool.net>\n"
		message += "<https://ubiq.phoenixmax.org>\n"
	case "!resettabs", "!blank":
		message = "In the Fusion URL bar, enter <https://wallet.ubiqsmart.io/?reset-tabs=true>\n"
	case "!roadmap", "!quarterly":
		message = "`Roadmap and quarterly report:` <https://blog.ubiqsmart.com/tagged/monthly> \n"
	case "!shinobi":
		message = "`You can access Shinobi by visiting:` <https://shinobi.ubiq.ninja>\n"
		message += "`Information can be found at:` <https://info.ubiq.ninja>\n"
	case "!social":
		message = "Follow us on social media `Medium:` <https://blog.ubiqsmart.com> `Twitter:` <https://twitter.com/ubiqsmart> `YouTube:` <https://www.youtube.com/ubiqvideos> `Reddit:` <https://www.reddit.com/r/Ubiq/> `GitHub:` <https://github.com/ubiq> `BitcoinTalk:` <https://bitcointalk.org/index.php?topic=1763606.0> `Telegram:` <https://t.me/Ubiqsmart> `Wikipedia:` <https://en.wikipedia.org/wiki/Ubiq>\n"
	case "!solidity":
		message = "`Solidity documentation:` <http://solidity.readthedocs.io>\n"
	case "!stats":
		message = "`Ubiq network stats:` <https://stats.ubiqscan.io>\n"
	case "!wallet":
		message = "`Web:` <https://pyrus.ubiqsmart.com> `Web2:` <https://myetherwallet.com> `GUI:` <https://github.com/ubiq/fusion/releases> `CLI:` <https://github.com/ubiq/go-ubiq/releases>\n"
	case "!website", "!site":
		message = "`Ubiq website:` <http://ubiqsmart.com/>	`Ubiq Community website:` <https://www.ubiqescher.com/>\n"
	case "!vyper":
		message = "`Vyper documentation:` <http://vyper.readthedocs.io>\n"
	}

	return &message
}
