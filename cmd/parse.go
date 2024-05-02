package main

import (
	"bytes"
	"encoding/hex"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type dhcpOption struct {
                tag             []byte
                length          []byte
                value           []byte
}

func worker(buf []byte, locAddr string, poolID int) {

	header := buf[0:240]
	dhcpopt := buf[240:]

	curOptions := parseDhcpOpt(dhcpopt)

	dhcpType := getOptValue(curOptions, 53)


	requestedOpts := getOptValue(curOptions, 55)
	log.Println(requestedOpts)

	chaddr := header[28:44]
	xid := header[4:8]

	clientID := getOptValue(curOptions, 61)

	if getPacketType(dhcpType) == request {
		log.Println("Received Request")
		reqIP := getOptValue(curOptions, 50)
		reqIPInfo := getAddrInfo(reqIP)

		if reqIPInfo.State == offered && bytes.Equal(reqIPInfo.CAddr, chaddr[:header[2]]) {
			log.Println("OFFER accepted")
			delLease(reqIP)
			reservIP(reqIP, used, chaddr[:int(header[2])])
			responsePkt := createPacket(ack, xid, reqIP, chaddr, poolID, requestedOpts, clientID)
			rpkt := responsePkt.Bytes()
			go sendEthernet(rpkt[16:20], rpkt[28:34], rpkt, DhcpConf.Pools[poolID].Interface)
		}else if reqIPInfo.State == used && bytes.Equal(reqIPInfo.CAddr, chaddr){

			log.Println("Lease Renewed")
			responsePkt := createPacket(ack, xid, reqIP, chaddr, poolID, requestedOpts, clientID)
			rpkt := responsePkt.Bytes()
			go sendEthernet(rpkt[16:20], rpkt[28:34], rpkt, DhcpConf.Pools[poolID].Interface)
		}else if reqIPInfo.State == 0 {
			reservIP(reqIP, used, chaddr[:int(header[2])])
			responsePkt := createPacket(ack, xid, reqIP, chaddr, poolID, requestedOpts, clientID)
			rpkt := responsePkt.Bytes()
			go sendEthernet(rpkt[16:20], rpkt[28:34], rpkt, DhcpConf.Pools[poolID].Interface)
		}

		return
	}else if getPacketType(dhcpType) == discover {

		reqIP := getOptValue(curOptions, 50)
		reqIPInfo := getAddrInfo(reqIP)

		var freeAddr net.IP
		
		//var freeAddrCopy net.IP
		freeAddrCopy := make(net.IP, 4)
		if bytes.Equal(reqIP, []byte{0}){

			reservedIP := searchMAC(chaddr[:6])

			if bytes.Equal(reservedIP, []byte{}) {
				freeAddrCopy = make(net.IP, len(freeAddr))
				log.Println("No IP requested; no existing lease")
				freeAddr = getFreeIP()
				freeAddrCopy = make(net.IP, len(freeAddr))
				copy(freeAddrCopy, freeAddr)
			}else {
				copy(freeAddrCopy, reservedIP)
			}
		}else {
			if bytes.Equal(reqIPInfo.CAddr, chaddr)  {
				log.Println("Requested IP is equal to reserved")
				copy(freeAddrCopy, reqIP)
			}else {
				reservedIP := searchMAC(chaddr)
				if bytes.Equal(reservedIP, []byte{}) {
					log.Println("ReqIP not available, no existing reservation")
					freeAddr = getFreeIP()
					freeAddrCopy = make(net.IP, len(freeAddr))
					copy(freeAddrCopy, freeAddr)
					log.Println("FreeAddrCopy: ", freeAddr)
				}else {
					log.Println("ReqIP available")
					copy(freeAddrCopy, reservedIP)
				}
			}
		}
		
		log.Println("FreeAddrCopy: ", freeAddrCopy)
		
		reservIP(freeAddrCopy, offered, chaddr[:int(header[2])])
		responsePkt := createPacket(offer, xid, freeAddrCopy, chaddr, poolID, requestedOpts, clientID)

		rpkt := responsePkt.Bytes()
		
		go sendEthernet(rpkt[16:20], rpkt[28:34], rpkt, DhcpConf.Pools[poolID].Interface)

		time.Sleep(600 * time.Second)

		reqIPInfo = getAddrInfo(getOptValue(curOptions, 50))
		if reqIPInfo.State == offered {
			go func(freeAddrCopy net.IP) {
    			delLease(freeAddrCopy)
			}(freeAddrCopy)
		}
		
	}else if getPacketType(dhcpType) == release {
		log.Println("Received RELEASE for: ", chaddr[:int(header[2])], searchMAC(chaddr[:int(header[2])]))
		delLease(searchMAC(chaddr[:int(header[2])]))
	}


	return
}

func createPacket(dhcptype pType, xid []byte, yiaddr []byte, chaddr []byte, poolID int, reqOpts []byte, clientID []byte) (bytes.Buffer) {

	var newPacket bytes.Buffer



	placeholder := make([]byte, 192)
	//lt := make([]byte, 4)

	newPacket.Write([]byte{2, 1, 6, 0})
	newPacket.Write(xid)
	newPacket.Write([]byte{0,1,0,0,0,0,0,0})
	newPacket.Write(yiaddr)
	newPacket.Write([]byte{10,0,202,2})
	newPacket.Write([]byte{0,0,0,0})
	newPacket.Write(chaddr)
	newPacket.Write(placeholder)
	newPacket.Write([]byte{99,130,83,99})

	newPacket = appendOption(newPacket, 53, []byte{byte(dhcptype)})
	newPacket = appendOption(newPacket, 54, ip2Byte(DhcpConf.ServerID))

	for _,opt := range reqOpts {
		optVal := getOptionValue(optionID(opt))
		log.Printf("-----------------------\nOPTVAL: %v", optVal)
		if !bytes.Equal(optVal, []byte{}) {
			log.Println("OptVal: ", optVal, reflect.TypeOf(optVal))
			newPacket = appendOption(newPacket, opt, []byte(optVal))
		}
	}

	newPacket = appendOption(newPacket, 51, getOptionValue(optLeaseTime))

	if !bytes.Equal(clientID, []byte{}) {
		newPacket = appendOption(newPacket, 61, clientID)
	}

	newPacket.Write([]byte{255})

	return newPacket
}

func appendOption(buf bytes.Buffer, tag byte, value []byte) (bytes.Buffer) {
	buf.Write([]byte{tag})
	buf.Write([]byte{byte(len(value))})
	buf.Write(value)
	return buf
}

func getOptValue(opts []dhcpOption, stag byte) ([]byte) {
	
	for _, option := range opts {
		tag := option.tag[0]
		if tag == stag {
			return option.value
		}
	}
	return []byte{0}
}

func parseDhcpOpt(buf []byte) ([]dhcpOption) {
	var optionList []dhcpOption
	var temp dhcpOption
	optBuf := bytes.NewBuffer(buf)
	lastTag := 0
	for lastTag != 255 {
		temp.tag = optBuf.Next(1)
		lastTag = int(temp.tag[0])
		temp.length = optBuf.Next(1)
		temp.value = optBuf.Next(int(temp.length[0]))
		optionList = append(optionList, temp)
	}
	return optionList
}

func getPacketType(buf []byte) (pType) {
	switch buf[0] {
	case 1:
		return discover
	case 2:
		return offer
	case 3:
		return request
	case 4:
		return decline
	case 5:
		return ack
	case 6:
		return nack
	case 7:
		return release
	case 8:
		return inform
	default:
		return dhcpError
	}
}

func getAddrInfo(reqIP []byte) (lease) {
	pool := readPool()
	for i := range pool {
		if bytes.Equal(pool[i].Addr, reqIP) {
			return pool[i]
		}
	}
	return lease{}
}

func leaseToStr( leases []lease) ([][]string) {
	
	var leaseStr [][]string
	var temp []string

	header := []string{"addr", "state", "cid"}

	leaseStr = append(leaseStr, header)

	for _,v := range leases {

		addrStr := strconv.Itoa(int(v.Addr[0])) + "." + strconv.Itoa(int(v.Addr[1])) + "." + strconv.Itoa(int(v.Addr[2])) + "." + strconv.Itoa(int(v.Addr[3])) 

		temp = append(temp, addrStr)
		temp = append(temp, rparseState(v.State))
		temp = append(temp, hex.EncodeToString(v.CAddr))
		
		leaseStr = append(leaseStr, temp)
		temp = nil
	}
	return leaseStr

}

func ip2Byte (s string) ([]byte) {
	var ip []byte
	addr := strings.Split(s, ".")

	for _, v := range addr {
		y,err := strconv.Atoi(v)

		if err != nil {
        	log.Println(err)
    	}
		ip = append(ip, byte(y))
	}

	return ip
}

