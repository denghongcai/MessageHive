package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/denghongcai/generalmessagegate/protocol"
	"log"
	"os"
	"strconv"
)

var controller = make(chan bool)

func main() {
	go client("a", "a")
	go client("b", "a")
	go client("c", "a")
	for {
		line, _ := bufio.NewReader(os.Stdin).ReadBytes('\n')
		op := fmt.Sprintf("%s", line)
		if op == "exit\n" {
			return
		}
		n, _ := strconv.Atoi(fmt.Sprintf("%s", line[:len(line)-1]))
		for i := 0; i < n; i++ {
			controller <- true
		}
	}
}

func client(sid string, rid string) {
	config := tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := tls.Dial("tcp", "server06.dhc.house:1430", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	msg := CreatePacket(SYS, sid, rid, `{"token": "a"}`)

	n, err := conn.Write(msg)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: wrote %q (%d bytes)", msg, n)

	readchan := make(chan *[]byte)

	go func() {
		if sid == rid {
			for {
				reply := make([]byte, 256)
				n, err = conn.Read(reply)
				readchan <- &reply
			}
		}
	}()

	for {
		select {
		case <-controller:
			msg = CreatePacket(EXCHANGE, sid, rid, RandSeq(5))

			n, err = conn.Write(msg)
			if err != nil {
				log.Fatalf("client: write: %s", err)
			}
		case reply := <-readchan:
			pkt := new(protocol.Packet)
			pkt.UnPack(reply)
			log.Printf("Receive from %s, body is \"%s\"", pkt.Msg.GetSID(), pkt.Msg.GetBODY())
		}
	}

}