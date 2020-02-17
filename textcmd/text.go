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
		message += "`!verify` - How to get Verified on Ubiq Discord\n"
		message += "\n"
		message += "__**Ubiq specific**__\n\n"
		message += "`!backup` - Backup your account\n"
		message += "`!bots` - Ubiq Twitter bots\n"
		message += "`!caps` - Correct spelling for Ubiq\n"
		message += "`!compare` - Comparison chart\n"
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
	case "!caps":
		message = "The correct spelling for Ubiq is 'Ubiq', not all capitals (UBIQ is incorrect). This is similar to Nike where the logo is in all capitals but not the name. The ticker symbol for Ubiq is UBQ.\n"
	case "!compare":
		message = "`Comparison chart provided by user moreexplosions:` https://imgur.com/a/Kr8RW"
	case "!escher":
		message = "`Escher contract address:` 0xcf3222b7FDa7a7563b9E1E6C966Bead04AC23c36  - Use default ABI and default 18 decimals."
	case "!escrow":
		message = "`Caution:` When trading OTC directly between yourself and another user, there's a chance you may get scammed by the other party and lose your funds. Whenever making an OTC trade, _always_ use an escrow. A trusted escrow protects both parties during the trade. Contact Ubiq Discord moderators to arrange a trustworthy escrow for you."
	case "!ethunits":
		message = "`Ethereum units:` <https://github.com/ryepdx/ethereum-units>\n"
	case "!exchange", "!market":
		message = "`Bittrex:` <https://bittrex.com/Market/Index?MarketName=BTC-UBQ> `Litebit:` <https://www.litebit.eu/en/buy/ubiq> `UPBit:` <https://upbit.com/exchange?code=CRIX.UPBIT.BTC-UBQ> `Catalx:` <https://catalx.io/trade/BTC-UBQ> `Changeangel:` <https://changeangel.io> `Xchange.me:` <https://xchange.me> `Cryptofacil:` <https://cryptofacil.com/trading-view/129/BTC-UBQ>\n"
	case "!explorer":
		message = "`Explorer:` <https://ubiqscan.io> `Explorer 2:` <http://www.ubiq.cc> `Explorer 3:` <https://ubiqexplorer.com>\n"
	case "!github":
		message = "`GitHub:` <https://github.com/ubiq>\n"
	case "!hide", "!hidechannels":
		message = "<https://support.discordapp.com/hc/en-us/articles/213599277-How-do-I-hide-channels->\n"
	case "!invite":
		message = "`Ubiq Discord invite link:` <https://discord.gg/HF6vEGF>\n"
	case "!miner":
		message = "`Ubqminer:` <https://github.com/ubiq/ubqminer/releases> `PhoenixMiner:` <https://bitcointalk.org/index.php?topic=2647654.0> `TT-Miner:` <https://bitcointalk.org/index.php?topic=5025783.0> `Nanominer:` <https://nanominer.org/>\n"
	case "!mp", "!monetarypolicy":
		message = "`Monetary policy and mining block rewards scheme in Ubiq:` <https://blog.ubiqsmart.com/ubiq-research-monetary-policy-2e27458983ec>\n"
	case "!nucleus", "!transparency":
		message = "`Nucleus Transparency Report:` <https://blog.ubiqsmart.com/nucleus-transparency-report-6496e444bd85>\n"
	case "!onepage":
		message = "`Ubiq in one page`: <https://blog.ubiqsmart.com/ubiq-in-one-page-3e3d335064fc>\n"
	case "!pools":
		message = "`List of known mining pools:`\n"
		message += "<http://ubq.altpool.pro>            <https://ubq.mypool.online>\n"
		message += "<https://ubq.zet-tech.eu>		<https://ubiqpool.fr/>		  <https://ubiq.mole-pool.net/>       <https://ubiqpool.maxhash.org>\n"
		message += "<https://ubiqpool.io>		<https://ubiq.clona.ru>		<http://ubq.pool.sexy>\n"
		message += "<http://ubiq.minerpool.net>    <https://ubiqpool.com>    <http://www.ubiq.cc/minerpool/>\n"
		message += "<https://ubq.solopool.org>     <https://ubiq.wattpool.net> <https://ubiq.gomine.pro/>\n"
	case "!resettabs", "!blank":
		message = "In the Fusion URL bar, enter <https://wallet.ubiqsmart.io/?reset-tabs=true>\n"
	case "!roadmap", "!quarterly":
		message = "`Roadmap and quarterly report:` `February:` <https://blog.ubiqsmart.com/ubiq-quarterly-report-march-2018-8ac54ca3eb78> `November:` <https://blog.ubiqsmart.com/ubiq-quarterly-report-november-2017-2acc66750490> `August:` <https://blog.ubiqsmart.com/ubiq-quarterly-report-august-2017-e6484a536b8d> `May:` <https://blog.ubiqsmart.com/ubiq-quarterly-report-10-5-2017-9750e297330f>\n"
	case "!social":
		message = "Follow us on social media `Medium:` <https://blog.ubiqsmart.com> `Twitter:` <https://twitter.com/ubiqsmart> `YouTube:` <https://www.youtube.com/ubiqvideos> `Reddit:` <https://www.reddit.com/r/Ubiq/> `GitHub:` <https://github.com/ubiq> `BitcoinTalk:` <https://bitcointalk.org/index.php?topic=1763606.0> `Telegram:` <https://t.me/Ubiqsmart> `Wikipedia:` <https://en.wikipedia.org/wiki/Ubiq>\n"
	case "!solidity":
		message = "`Solidity documentation:` <http://solidity.readthedocs.io>\n"
	case "!stats":
		message = "`Ubiq network stats:` <https://ubiq.darcr.us>\n"
	case "!verify", "!verified":
		message = "__**Verified**__"
		message += "\n"
		message += "You can request to be Verified in the <#348034489655623680> channel. This allows us to see who is a community member and lets you send messages.\n\n"
		message += "This policy is in place to prevent bots, spammers and scammers.\n"
		message += "To get Verified:\n"
		message += "a) Mention you would like to be Verified in the <#348034489655623680> channel. Mentioning where you came from and how you found out about Ubiq helps too.\n"
		message += "b) Only people with a profile pic are allowed to be Verified. It doesn't have to be your actual photo. It can be any appropriate image you like."
	case "!wallet":
		message = "`Web:` <https://pyrus.ubiqsmart.com> `Web2:` <https://myetherwallet.com> `GUI:` <https://github.com/ubiq/fusion/releases> `CLI:` <https://github.com/ubiq/go-ubiq/releases>\n"
	case "!website", "!site":
		message = "`Ubiq website:` <http://ubiqsmart.com/>	`Ubiq Community website:` <https://www.ubiqescher.com/>\n"
	case "!shokku":
		message = "If you are familiar with Infura <https://infura.io> on Ethereum, Shokku (<https://shokku.com> - website will be available soon) provides the same functionality (public API for dapps to interact directly with the chain without running their own nodes, and ipfs). Something to note is that this project is not based on forked code, it has been written from scratch as Infura is not open source. If you want more information or need assitance to use the service, just ping aldoborrero.\n"
	case "!vyper":
		message = "`Vyper documentation:` <http://vyper.readthedocs.io>\n"
	}

	return &message
}
