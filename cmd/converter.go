package main

import(
	"net"
	"strconv"
	"encoding/binary"
	"log"
)

func intToByteArray(num int) []byte {
    var byteArr []byte

    for num > 0 {
        digit := num % 256
        byteArr = append([]byte{byte(digit)}, byteArr...)
        num >>= 8
    }

    return byteArr
}

func parseString(input string) []byte {
    if ip := net.ParseIP(input); ip != nil {
		log.Println("ip")
		return ip[12:]
    }

    if i, err := strconv.ParseInt(input, 10, 32); err == nil {
		bs := make([]byte, 4)
		log.Printf("INT32: %v, %v", i, bs)
    	binary.BigEndian.PutUint32(bs,uint32(i))
		log.Println("int32")
        return bs
    }

    if i, err := strconv.ParseInt(input, 10, 16); err == nil {
		bs := make([]byte, 2)
    	binary.BigEndian.PutUint16(bs, uint16(i))
		log.Println("int16")
        return bs
    }

	log.Println("byte")
    return []byte(input)
}

func parseState(s string) (addrState) {
	switch s{
	case "f":
		return free
	case "r":
		return reserved
	case "u":
		return used
	case "o":
		return offered
	default:
		return missing
	}
}

func rparseState(s addrState) (string) {
	switch s{
	case free:
		return "f"
	case reserved:
		return "r"
	case used:
		return "u"
	case offered:
		return "o"
	default:
		return ""
	}
}
