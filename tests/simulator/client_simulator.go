package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/denghongcai/MessageHive/modules/protocol"
)

var controller = make(chan bool)

func main() {
	sid := flag.String("s", "a", "sid")
	rid := flag.String("r", "a", "rid")
	group := flag.Bool("g", false, "group")
	flag.Parse()
	fmt.Printf("sid: %s\n", *sid)
	fmt.Printf("rid: %s\n", *rid)
	go client(*sid, *rid, *group)
	go func() {
		for {
			time.Sleep(10 * time.Second)
			controller <- true
		}
	}()
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

func client(sid string, rid string, group bool) {
	config := tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := tls.Dial("tcp", "server01.dhc.house:1430", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	msg := CreatePacket(AUTH, sid, rid, `{"token": "hehe"}`)

	n, err := conn.Write(msg)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: wrote %q (%d bytes)", msg, n)

	if group {
		msg = CreatePacket(JOINGROUP, sid, "d", `{"action": "join", "data": ["a", "b", "c"]}`)

		n, err = conn.Write(msg)
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		log.Printf("client: wrote %q (%d bytes)", msg, n)

	}
	readchan := make(chan *[]byte)

	go func() {
		for {
			reply := make([]byte, 256)
			n, err = conn.Read(reply)
			readchan <- &reply
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
