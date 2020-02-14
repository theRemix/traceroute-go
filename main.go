package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
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

func sendICMP(r *ipv4.RawConn, destination net.IP, ttl int) {

	res := make(chan net.IP, 1)
	// timeout := make(chan string, 1)
	go func() {
		ipHeader, icmpReq := buildICMPRequest(destination, ttl)

		cmReq := &ipv4.ControlMessage{}

		if err := r.WriteTo(ipHeader, icmpReq, cmReq); err != nil {
			log.Println(err)
		}

		responseBytes := make([]byte, 1500)
		responseHeader, cmRes, peer, err := r.ReadFrom(responseBytes)
		if err != nil {
			log.Println(err)
		}

		rm, err := icmp.ParseMessage(1, cmRes)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("type", rm.Type)
		switch rm.Type {
		case ipv4.ICMPTypeTimeExceeded:
			log.Println("hop", peer.String())
			// names, _ := net.LookupAddr(peer.String())
			// fmt.Printf("%d\t%v %+v %v\n\t%+v\n", i, peer, names, rtt, cm)
		case ipv4.ICMPTypeEchoReply:
			log.Println("end", peer.String())
			// names, _ := net.LookupAddr(peer.String())
			// fmt.SPrintf("%d\t%v %+v %v\n\t%+v\n", i, peer, names, rtt, cm)
			// return
		case ipv4.ICMPTypeDestinationUnreachable:
			log.Println("dest unreachable", peer.String())
		default:
			log.Printf("unknown ICMP message: %+v\n", rm)
		}

		res <- responseHeader.Src
	}()

	select {
	case res := <-res:
		fmt.Println(res)
	case <-time.After(5 * time.Second):
		fmt.Println("timeout", ttl)
	}

}

func main() {

	// destination := net.IP{8, 8, 8, 8}
	destination := net.IP{104, 20, 40, 243}
	// destination := net.IP{104, 244, 42, 65}

	listener, err := net.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	rawCon, err := ipv4.NewRawConn(listener)
	if err != nil {
		log.Fatal(err)
	}
	if err := rawCon.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
		log.Fatal(err)
	}

	// t := time.Now().Add(5 * time.Second)
	// rawCon.SetReadDeadline(t)

	for ttl := 1; ttl <= 64; ttl++ {
		log.Println("sending", ttl)
		sendICMP(rawCon, destination, ttl)
	}

}
