package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device       string = "eth0"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 100 * time.Millisecond
	handle       *pcap.Handle
)

type SyncAckPack struct {
	Time            time.Time
	SourceIP        string
	DestinationIP   string
	SourcePort      string
	DestinationPort string
	AckSeq          uint32
}

func main() {
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	var filter string = "tcp port 80"
	err = handle.SetDirection(pcap.DirectionInOut)
	if err != nil {
		log.Fatal(err)
	}
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Only capturing " + filter + " packets.")

	var session = map[uint32]map[string]SyncAckPack{}

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Process packet here
		iptype := packet.NetworkLayer().LayerType()

		var srcIP, dstIP string

		if iptype == layers.LayerTypeIPv4 {
			ip := packet.NetworkLayer().(*layers.IPv4)
			srcIP = ip.SrcIP.String()
			dstIP = ip.DstIP.String()
		} else if iptype == layers.LayerTypeIPv6 {
			ip := packet.NetworkLayer().(*layers.IPv6)
			srcIP = ip.SrcIP.String()
			dstIP = ip.DstIP.String()
		}

		tcp := packet.TransportLayer().(*layers.TCP)

		//print the map
		fmt.Println(session)

		if tcp.SYN && tcp.ACK {
			if _, ok := session[tcp.Ack]; !ok {
				session[tcp.Ack] = map[string]SyncAckPack{}
			}

			session[tcp.Ack][srcIP] = SyncAckPack{
				Time:            packet.Metadata().Timestamp,
				SourceIP:        srcIP,
				DestinationIP:   dstIP,
				SourcePort:      tcp.SrcPort.String(),
				DestinationPort: tcp.DstPort.String(),
				AckSeq:          tcp.Ack,
			}
		} else if tcp.ACK {
			if _, ok := session[tcp.Seq]; ok {
				if _, ok := session[tcp.Seq][dstIP]; ok {
					syncack := session[tcp.Seq][dstIP]

					fmt.Println("Time:", packet.Metadata().Timestamp.Sub(syncack.Time), "SourceIP:", syncack.SourceIP, "DestinationIP:", syncack.DestinationIP, "SourcePort:", syncack.SourcePort, "DestinationPort:", syncack.DestinationPort, "AckSeq:", syncack.AckSeq)

					delete(session[tcp.Seq], dstIP)
				}
			}
		}

	}
}
