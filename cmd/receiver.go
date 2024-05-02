package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"syscall"
)

func receiver(wg *sync.WaitGroup, iface string, poolID int) {

	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	//f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	syscall.SetsockoptString(fd, syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, iface)
	addr := syscall.SockaddrInet4{Port: 67}

	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	if err := syscall.Bind(fd, &addr); err != nil {
		log.Fatal("Error binding:", err)
	}

	log.Println("Waiting...")
	for {
		buf := make([]byte, 1024)
		_, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			log.Println("Error reading packet: ", err)
		}
		fmt.Printf("Received on %v \n",iface)
		go worker(buf, iface, poolID)
	}
	wg.Done()
}
