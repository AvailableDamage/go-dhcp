package main

import (
	"log"
	"fmt"
	"sync"
	"time"
)

type pType byte
type addrState byte

const (
	discover 	pType = 1
	offer 		pType = 2
	request		pType = 3
	decline		pType = 4
	ack 		pType = 5
	nack 		pType = 6
	release		pType = 7
	inform		pType = 8
	dhcpError	pType = 0
)

const (
	free		addrState = 1
	reserved	addrState = 2
	used		addrState = 3
	offered		addrState = 4
	missing		addrState = 0
)

type lease struct {
	Addr 	[]byte
	State 	addrState
	CAddr 	[]byte
	Time 	int
}

var DhcpConf = configuration{}

func main() {
	cleanTicker := time.NewTicker(60 * time.Second)
	var wg sync.WaitGroup

	log.Println("Started DHCP Server")
	DhcpConf = loadconf()
	getOptionValue(51)
	log.Println("Nameserver:", DhcpConf.Pools[0].Options["Nameserver"], optionCodeToString[6])

	for i,pool := range DhcpConf.Pools {
		wg.Add(1)
		fmt.Println("interface: ", pool)
		go receiver(&wg, pool.Interface, i)
	}
	
	for {
		select {
    	case <-cleanTicker.C:
    	    log.Println("Clean leases")
			cleanLeases()
    	}
	}

	wg.Wait()
}
