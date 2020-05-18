package main

import (
	"github.com/happyxcj/gosocket"
	"github.com/happyxcj/gosocket/protocol"
	"github.com/happyxcj/gosocket/pkts"
	"time"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	c, _ := gosocket.NewAndInitClient(":8080",
		gosocket.QCOptions(
			gosocket.PktHandler(handlePkt),
			gosocket.OnClose(func(e error) {
				fmt.Println("current connection is closed, please subscribe to the services again")
				// Just for test.
				os.Exit(1)
			})))
	defer c.Close()
	go pingLoop(c)
	mockSend(c)
	select {}
}

func pingLoop(c *gosocket.Client) {
	for {
		packet := pkts.NewPingPkt()
		c.Send(packet)
		time.Sleep(5 * time.Second)
	}
}

type Message struct {
	Id      int
	Content string
}

type ErrResp struct {
	Code int
	Msg  string
}

func mockSend(c *gosocket.Client) {
	notification := &Message{Id: 1, Content: "hello"}
	data, _ := json.Marshal(notification)
	pkt := pkts.NewEasyNotifyPkt(1000, data)
	pkt.Props().WithInt64(pkts.PropCreatedTime, time.Now().Unix())
	c.Send(pkt)

	time.Sleep(time.Second)

	subMsg := &Message{Id: 1, Content: "hello"}
	data, _ = json.Marshal(subMsg)
	subPkt := pkts.NewEasySubPkt(10, 1000, data)
	pkt.Props().WithInt64(pkts.PropCreatedTime, time.Now().Unix())
	c.Send(subPkt)
}

func handlePkt(p protocol.Packet) {
	switch v := p.(type) {
	case *pkts.SubPkt:
		handleSubPkt(v)
	case *pkts.PubPkt:
		handlePubPkt(v)
	}
}

func handleSubPkt(p *pkts.SubPkt) {
	switch p.Cmd() {
	case 1000:
		fmt.Printf("subscribe '%v' successfully\n", p.Desc())
	case 110:
		appErr := &ErrResp{}
		json.Unmarshal(p.Body(), appErr)
		fmt.Printf("unable to subscrbie service, request id is: %v, reason: %v\n", p.SeqId(), appErr.Msg)
	}

}

func handlePubPkt(p *pkts.PubPkt) {
	switch p.Cmd() {
	case 1000:
		msg := &Message{}
		json.Unmarshal(p.Body(), msg)
		fmt.Println("receive publish message: ", msg.Content)
	}
}
