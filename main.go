package main

import (
	"encoding/binary"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

const (
	ICMP_TYPE_REPLY   = 0
	ICMP_TYPE_REQUEST = 8
	IP_PROTOCOL_ICMP  = 1
)

// // 20b ip header
// ip := make([]byte, 20)

// // version
// ip[0] = 4

// // tos
// ip[1] = 0

// // length
// binary.BigEndian.PutUint16(ip[2:], 0x0000)

// // identification
// binary.BigEndian.PutUint16(ip[4:], 0x0000)

// // flags and offset
// binary.BigEndian.PutUint16(ip[6:], 0x0000)

// // ttl
// ip[8] = byte(ttl)

// // protocol
// ip[9] = 0

// // header checksum
// binary.BigEndian.PutUint16(ip[10:], 0x0000)

// // source ip address
// binary.BigEndian.PutUint32(ip[12:], 0x0000)

// // destination ip address
// binary.BigEndian.PutUint32(ip[16:], 0x0000)

func buildICMPRequest(destination net.IP, ttl int) (*ipv4.Header, []byte) {

	// 12b icmp header + 4b payload
	icmp := make([]byte, 12)

	// type
	icmp[0] = ICMP_TYPE_REQUEST

	// code
	icmp[1] = 0

	// @TODO checksum
	binary.BigEndian.PutUint16(icmp[2:], 0x0000)

	// identifier
	binary.BigEndian.PutUint16(icmp[4:], uint16(ttl))

	// sequence number
	binary.BigEndian.PutUint16(icmp[6:], uint16(ttl))

	// payload
	binary.BigEndian.PutUint32(icmp[8:], 0xdeadbeef)

	ip := ipv4.Header{
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TOS:      0,
		Flags:    ipv4.DontFragment,
		TotalLen: ipv4.HeaderLen + len(icmp),
		TTL:      ttl,
		Protocol: IP_PROTOCOL_ICMP,
		Dst:      destination,
	}

	return &ip, icmp
}

func main() {

	destination := net.IP{8, 8, 8, 8}

	c, err := net.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	r, err := ipv4.NewRawConn(c)
	if err != nil {
		log.Fatal(err)
	}

	ipHeader, icmp := buildICMPRequest(destination, 1)

	cm := &ipv4.ControlMessage{}

	if err := r.WriteTo(ipHeader, icmp, cm); err != nil {
		log.Println(err)
	}

	responseBytes := make([]byte, 1500)
	responseHeader, _, _, err := r.ReadFrom(responseBytes)
	if err != nil {
		log.Println(err)
	}

	log.Println(responseHeader.Src)

}
