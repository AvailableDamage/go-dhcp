package main

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var leasefile = "/var/run/goDHCP/leases.csv"
var mutex sync.Mutex

func getFreeIP() (net.IP) {
	log.Println("Get free IP...", DhcpConf.Pools[0].StartIP)

	f,err := os.Open(leasefile)

	defer f.Close()

	if err != nil {
		log.Println("Error reading leases: ", err)
	}

    csvReader := csv.NewReader(f)
    data,err := csvReader.ReadAll()
	var startIP []byte
	var endIP []byte

	var y int
	var i []byte
	var line []string

	startIP = DhcpConf.Pools[0].StartIP[12:]
	endIP = DhcpConf.Pools[0].EndIP[12:]

	for i=startIP; bytes.Compare(i, endIP) != 0;i[3]++ {
		used := false
		for y,line = range data {
			if y == 0{
				continue
			}
			if bytes.Equal(ip2Byte(line[0]), i) {
				used = true
				break
			}
			used = false
		}
		if !used {
			return i
		}else {
			continue
		}
	}
	return net.IP{} 
}

func editLease(ip net.IP) {

    mutex.Lock()
    defer mutex.Unlock()
	f,err := os.Open(leasefile)


	if err != nil {
		log.Println("Error opening leases for reading: ", err)
	}

    csvReader := csv.NewReader(f)
	
	var temp [][]string

    data,err := csvReader.ReadAll()

	for y,line := range data {

		if y == 0 {
			temp = append(temp, line)
		}else if bytes.Compare(ip2Byte(line[0]), ip) == 0 {
			temp = append(temp, line)
			temp[y][1] = rparseState(used)
		}else {
			temp = append(temp, line)
		}
	}

	f.Close()

	f,err = os.OpenFile(leasefile, os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0644)

	defer f.Close()

	if err != nil {
		log.Println("Error opening leases for writing: ", err)
	}


	if err := os.Truncate(leasefile, 0); err != nil {
    	log.Printf("Failed to truncate: %v", err)
	}	

	csvWriter := csv.NewWriter(f)
	csvWriter.WriteAll(temp)
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Fatal(err)
	}
	log.Println(temp)
}

func delLease(ip net.IP) {
    mutex.Lock()
    defer mutex.Unlock()
	f,err := os.Open(leasefile)


	if err != nil {
		log.Println("Error opening leases for reading: ", err)
	}

    csvReader := csv.NewReader(f)
	
	var temp [][]string

    data,err := csvReader.ReadAll()

	for y,line := range data {

		if y == 0 {
			temp = append(temp, line)
		}else if bytes.Compare(ip2Byte(line[0]), ip) == 0 {
		}else {
			temp = append(temp, line)
		}
	}

	f.Close()

	f,err = os.OpenFile(leasefile, os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0644)

	defer f.Close()

	if err != nil {
		log.Println("Error opening leases for writing: ", err)
	}


	if err := os.Truncate(leasefile, 0); err != nil {
    	log.Printf("Failed to truncate: %v", err)
	}	

	csvWriter := csv.NewWriter(f)
	csvWriter.WriteAll(temp)
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Fatal(err)
	}
}

func reservIP(ip net.IP, state addrState, chaddr []byte) {
	f,err := os.OpenFile(leasefile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	defer f.Close()

	if err != nil {
		log.Println("Error opening leases for writing: ", err)
	}

    csvWriter := csv.NewWriter(f)

	var lease []string

	lease = append(lease, ip.String())
	lease = append(lease, rparseState(state))
	lease = append(lease, hex.EncodeToString(chaddr))
	lease = append(lease, strconv.Itoa(int(time.Now().Unix())))
	fmt.Println("Time: ", time.Now().Unix())

	csvWriter.Write(lease)

	csvWriter.Flush()

}

func cleanLeases() {
	
    mutex.Lock()
    defer mutex.Unlock()
	leases := readPool()

	for _,lease := range leases {
		now := time.Now().Unix()
		if int(now) >= lease.Time+DhcpConf.LeaseTime {
			log.Println("lease expired")
			mutex.Unlock()
			delLease(lease.Addr)
			mutex.Lock()
		}
	}

}

func searchMAC(chaddr []byte) ([]byte) {
	leases := readPool()

	for _,line := range leases {
		if bytes.Equal(chaddr, line.CAddr) {
			return line.Addr
		}
	}
	return []byte{}
}

func readPool() ([]lease){

    f, err := os.Open(leasefile)
    if err != nil {
        fmt.Println(err)
    }

    csvReader := csv.NewReader(f)
    data,err := csvReader.ReadAll()
    if err != nil {
        log.Println(err)
    }
	var temp lease
	var leases []lease
    for i,line := range data {
        if i > 0 {
			temp.Addr = ip2Byte(line[0])
			temp.State = parseState(line[1])
			tempchaddr,_ := hex.DecodeString(line[2])
			temp.CAddr = tempchaddr
			temp.Time,_ = strconv.Atoi(line[3])
			leases = append(leases, temp)
        }
    }

	defer f.Close()
	return leases
}
