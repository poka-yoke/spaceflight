package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/olorin/nagiosplugin"
)

// This list of BLs was updated from http://www.blalert.com/dnsbls on Nov 2016.
var totalList = []string{
	"images.rbl.msrbl.net",
	// "http.opm.blitzed.org",
	// "duinv.aupads.org",
	// "bl.mailspike.org",
	"bl.mailspike.net",
	"mail-abuse.blacklist.jippg.org",
	"rbl.efnetrbl.org",
	"rbl.polarcomm.net",
	"cbl.anti-spam.org.cn",
	"sbl.spamhaus.org",
	// "fnrbl.fast.net",
	"bad.psky.me",
	// "ip.v4bl.org",
	"free.v4bl.org",
	"relays.nether.net",
	// "dev.null.dk",
	"spamsources.fabel.dk",
	// "lookup.dnsbl.iip.lu",
	// "rbl.choon.net",
	"cblplus.anti-spam.org.cn",
	// "rbl.orbitrbl.com",
	"all.rbl.jp",
	// "pss.spambusters.org.ar",
	"spamlist.or.kr",
	"list.blogspambl.com",
	// "cart00ney.surriel.com",
	// "spamsites.dnsbl.net.au",
	"blacklist.sci.kun.nl",
	"psbl.surriel.com",
	// "blackholes.five-ten-sg.com",
	"owfs.dnsbl.net.au",
	// "msgid.bl.gweep.ca",
	// "blackholes.wirehub.net",
	// "dnsbl.solid.net",
	"ips.backscatterer.org",
	// "dnsbl.ahbl.org",
	"spam.abuse.ch",
	"0spam.fusionzero.com",
	// "dnsbl.njabl.org",
	// "relays.bl.gweep.ca",
	// "singlebl.spamgrouper.com",
	// "bl.deadbeef.com",
	"st.technovision.dk",
	// "t1.dnsbl.net.au",
	// "dnsbl.antispam.or.id",
	"drone.abuse.ch",
	"http.dnsbl.sorbs.net",
	"block.dnsbl.sorbs.net",
	"bl.spamcannibal.org",
	"l2.apews.org",
	// "probes.dnsbl.net.au",
	"combined.abuse.ch",
	// "blocked.hilli.dk",
	"bl.spamcop.net",
	"dnsrbl.swinog.ch",
	"ix.dnsbl.manitu.net",
	// "spamtrap.drbl.drand.net",
	// "bhnc.njabl.org",
	// "hil.habeas.com",
	"dnsbl-2.uceprotect.net",
	"rbl.interserver.net",
	// "httpbl.abuse.ch",
	// "dnsbl.burnt-tech.com",
	"short.rbl.jp",
	"dnsbl-3.uceprotect.net",
	"dnsbl.dronebl.org",
	"spam.rbl.msrbl.net",
	// "will-spam-for-food.eu.org",
	"dnsbl.rv-soft.info",
	"socks.dnsbl.sorbs.net",
	"bl.blocklist.de",
	// "rsbl.aupads.org",
	"orvedb.aupads.org",
	// "csi.cloudmark.com",
	"b.barracudacentral.org",
	// "sbl-xbl.spamhaus.org",
	"spam.dnsbl.sorbs.net",
	// "blackholes.mail-abuse.org",
	"virus.rbl.jp",
	// "bl.csma.biz",
	// "unconfirmed.dsbl.org",
	// "combined.njabl.org",
	"relays.bl.kundenserver.de",
	// "l2.bbfh.ext.sorbs.net",
	"zombie.dnsbl.sorbs.net",
	"ubl.unsubscore.com",
	// "wingate.opm.blitzed.org",
	// "dialups.visi.com",
	"netblock.pedantic.org",
	// "bl.score.senderscore.com",
	"korea.services.net",
	// "virbl.bit.nl",
	// "dialups.mail-abuse.org",
	// "l2.spews.dnsbl.sorbs.net",
	"misc.dnsbl.sorbs.net",
	"virus.rbl.msrbl.net",
	// "multihop.dsbl.org",
	// "dialup.blacklist.jippg.org",
	"dnsbl.sorbs.net",
	// "rmst.dnsbl.net.au",
	// "omrs.dnsbl.net.au",
	// "ucepn.dnsbl.net.au",
	// "owps.dnsbl.net.au",
	"zen.spamhaus.org",
	// "fl.chickenboner.biz",
	"smtp.dnsbl.sorbs.net",
	// "spam.olsentech.net",
	// "spamtrap.trblspam.com",
	// "bl.technovision.dk",
	"bl.emailbasura.org",
	// "proxy.bl.gweep.ca",
	// "all.spamrats.com",
	"spam.spamrats.com",
	"wormrbl.imp.ch",
	// "work.drbl.gremlin.ru",
	// "intruders.docs.uu.se",
	// "relays.mail-abuse.org",
	"dul.dnsbl.sorbs.net",
	// "rdts.dnsbl.net.au",
	"cbl.abuseat.org",
	"cdl.anti-spam.org.cn",
	// "query.senderbase.org",
	// "rbl2.triumf.ca",
	// "orid.dnsbl.net.au",
	// "multi.surbl.org",
	"dnsbl.justspam.org",
	"combined.rbl.msrbl.net",
	"virbl.dnsbl.bit.nl",
	// "spam.wytnij.to",
	"tor.dnsbl.sectoor.de",
	"pbl.spamhaus.org",
	// "whois.rfc-ignorant.org",
	// "osrs.dnsbl.net.au",
	"bl.spameatingmonkey.net",
	"cidr.bl.mcafee.com",
	// "dul.ru",
	// "no-more-funn.moensted.dk",
	"dnsbl-1.uceprotect.net",
	"db.wpbl.info",
	// "rbl.suresupport.com",
	// "dsbl.dnsbl.net.au",
	// "sbl.csma.biz",
	"phishing.rbl.msrbl.net",
	"rbl.dns-servicios.com",
	// "spews.dnsbl.net.au",
	// "access.redhawk.org",
	"spamrbl.imp.ch",
	// "proxy.block.transip.nl",
	// "torserver.tor.dnsbl.sectoor.de",
	// "residential.block.transip.nl",
	"truncate.gbudb.net",
	// "map.spam-rbl.com",
	"rbl.schulte.org",
	"xbl.spamhaus.org",
	"spamguard.leadmon.net",
	"noptr.spamrats.com",
	// "forbidden.icm.edu.pl",
	// "opm.tornevall.org",
	// "ohps.dnsbl.net.au",
	// "dnsbl.webequipped.com",
	// "rbl.snark.net",
	"dnsbl.inps.de",
	"spam.pedantic.org",
	// "sorbs.dnsbl.net.au",
	"blacklist.woody.ch",
	// "t3direct.dnsbl.net.au",
	// "ubl.lashback.com",
	// "ricn.dnsbl.net.au",
	"dnsbl.kempt.net",
	"dynip.rothen.com",
	// "dsn.rfc-ignorant.org",
	// "tor.dan.me.uk",
	// "osps.dnsbl.net.au",
	// "rbl.triumf.ca",
	"web.dnsbl.sorbs.net",
	// "socks.opm.blitzed.org",
	// "l1.spews.dnsbl.sorbs.net",
	// "rbl.spamlab.com",
	"bogons.cymru.com",
	"dyna.spamrats.com",
	// "rbl-plus.mail-abuse.org",
	// "dnsbl.cyberlogic.net",
	"all.s5h.net",
}

// Reverse reverses slice of string elements.
func Reverse(original []string) {
	for i := len(original)/2 - 1; i >= 0; i-- {
		opp := len(original) - 1 - i
		original[i], original[opp] = original[opp], original[i]
	}
}

// ReverseAddress converts IP address in string to reversed address for query.
func ReverseAddress(ipAddress string) (reversedIPAddress string) {
	ipAddressValues := strings.Split(ipAddress, ".")
	Reverse(ipAddressValues)
	reversedIPAddress = strings.Join(ipAddressValues, ".")
	return
}

// DNSBLQuery queries a DNSBL and returns true if the argument gets a match
// in the BL.
func DNSBLQuery(ipAddress, bl string, addresses chan int) {
	reversedIPAddress := fmt.Sprintf(
		"%v.%v",
		ReverseAddress(ipAddress),
		bl,
	)
	result, _ := net.LookupHost(reversedIPAddress)
	if len(result) > 0 {
		log.Printf("%v present in %v(%v)", reversedIPAddress, bl, result)
	}
	addresses <- len(result)
}

func main() {
	ipAddress := flag.String(
		"ip",
		"127.0.0.1",
		"IP Address to look for in the BLs",
	)
	warning := flag.Int("w", 90, "Warning threshold")
	critical := flag.Int("c", 95, "Critical threshold")

	flag.Parse()

	check := nagiosplugin.NewCheck()
	defer check.Finish()
	responses := make(chan int, len(totalList))

	queried := 0
	positive := 0
	for i := range totalList {
		go DNSBLQuery(*ipAddress, totalList[i], responses)
	}
	for i := 0; i < len(totalList); i++ {
		response := <-responses
		if response > 0 {
			positive += response
		}
		queried++
	}
	warningAmount := len(totalList) * (*warning) / 100
	criticalAmount := len(totalList) * (*critical) / 100
	checkLevel := nagiosplugin.OK
	if positive > warningAmount {
		checkLevel = nagiosplugin.WARNING
		if positive > criticalAmount {
			checkLevel = nagiosplugin.CRITICAL
		}
	}
	check.AddResult(
		checkLevel,
		fmt.Sprintf(
			"%v present in %v(%v%%) out of %v BLs | %v",
			*ipAddress,
			positive,
			positive*100/len(totalList),
			len(totalList),
			fmt.Sprintf(
				"queried=%v positive=%v",
				queried,
				positive,
			),
		),
	)
}
