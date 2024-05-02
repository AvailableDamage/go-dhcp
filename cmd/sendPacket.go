package main

import(
	"net"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"syscall"
	"log"
)

func sendEthernet(ciaddr net.IP, chaddr net.HardwareAddr, dhcpPayload []byte, iface string) {
	
	nic,_ := net.InterfaceByName(iface)


	if nic.HardwareAddr == nil {
		nic.HardwareAddr,_ = net.ParseMAC("aa:bb:cc:dd:ee:ff")
	}

	log.Println("Source MAC: ", nic.HardwareAddr)
	ethernet := &layers.Ethernet{
		SrcMAC: nic.HardwareAddr, 
		DstMAC: chaddr,
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip := &layers.IPv4{
		SrcIP: dhcpPayload[20:24],
		DstIP: ciaddr,
		Version: 4,
		TTL: 64,
		Protocol: layers.IPProtocolUDP,		
		Flags:    layers.IPv4DontFragment,
	}
	var serverPort layers.UDPPort
	var clientPort layers.UDPPort

	serverPort = 67
	clientPort = 68

	udp := &layers.UDP{
		DstPort: clientPort,
		SrcPort: serverPort,
	}

	err := udp.SetNetworkLayerForChecksum(ip)
	if err != nil {
		log.Printf("Send Ethernet: Couldn't set network layer: %v", err)
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	} 
	gopacket.SerializeLayers(buf, opts,
		ethernet,
		ip,
		udp,
		gopacket.Payload(dhcpPayload),
	)

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, 0)

	defer func() {
		err = syscall.Close(fd)
		if err != nil {
		}
	}()

	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		log.Println("Cant set socket Options: ", err)
	}


	var hwAddr [8]byte
	copy(hwAddr[:], nic.HardwareAddr)

	ethAddr := syscall.SockaddrLinklayer{
		Protocol: 0,
		Ifindex:  nic.Index,
		Halen:    6,
		Addr:     hwAddr,
	}
	err = syscall.Sendto(fd, buf.Bytes(), 0, &ethAddr)
	log.Println("Packet sent...", ciaddr)
}
