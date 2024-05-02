package main

import(
)

type optionID byte
type optionType byte

const (
	optIP			optionType = 1
	optStr 			optionType = 2
	optByte 		optionType = 3
	optInt 			optionType = 4
)

const (
	optNetMask		optionID = 1
	optRouter		optionID = 3
	optDnsSrv		optionID = 6
	optDomain		optionID = 15
	optBroadcast	optionID = 28
	optNtpSrv		optionID = 42
	optLeaseTime	optionID = 51
	optType			optionID = 53
	optServerID		optionID = 54
	optT1			optionID = 58
	optT2			optionID = 59
	optEnd			optionID = 255
)

var optionCodeToString = map[optionID]string{
	optNetMask:		"Subnetmask",
	optRouter: 		"Gateway",
	optDnsSrv: 		"Nameserver",
	optDomain:	 	"Domain",
	optBroadcast: 	"Broadcast Address",
	optNtpSrv: 		"Network Time Server",
	optLeaseTime: 	"Lease Time",
	optType: 		"DHCP Type",
	optServerID: 	"Server ID",
	optT1: 			"Renewal Time",
	optT2: 			"Rebind Time",
	optEnd:			"End",
}


var optionCodeToType = map[optionID]optionType {
	optNetMask:		optIP,
	optRouter: 		optIP,
	optDnsSrv: 		optIP,
	optDomain: 		optStr,
	optBroadcast: 	optIP,
	optNtpSrv: 		optIP,
	optLeaseTime: 	optInt,
	optType: 		optByte,
	optServerID: 	optIP,
	optT1: 			optInt,
	optT2: 			optInt,
	optEnd:			optByte,
}

