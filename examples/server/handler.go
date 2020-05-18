package main

import (
	"github.com/happyxcj/gosocket/route"
	"github.com/happyxcj/gosocket/pkts"
	"encoding/json"
	"fmt"
	"time"
)

func initHandlers() {
	route.Group(pkts.KindPing).Handle("ping", nil, handlePingResp)

	route.Group(pkts.KindNotify).Use(handleMsgDecode).
		Handle("1000", Message{}, handleNotifyMsg)

	route.Group(pkts.KindSubscribe).
		Use(handleReqDecode).
		UseAfter(handleResp).
		Handle("1000", Message{}, handleSubMsg)
}

type Message struct {
	Id      int
	Content string
}

func handlePingResp(c *route.Context) {
	if !c.Pkt.Flags().Has(pkts.FlagPong) {
		c.Conn.Send(pkts.NewPongPkt())
	}
}

func handleMsgDecode(c *route.Context) {
	pkt := c.Pkt.(pkts.DataPkt)
	// The request packet has the application message.
	err := json.Unmarshal(pkt.Body(), c.Msg)
	if err == nil {
		return
	}
	fmt.Printf("unable to decode message: %v, error: %v\n", pkt.Desc(), err)
	c.Abort()
}

type ErrResp struct {
	Code int
	Msg  string
}

func handleReqDecode(c *route.Context) {
	pkt := c.Pkt.(pkts.DataPkt)
	err := json.Unmarshal(pkt.Body(), c.Msg)
	if err != nil {
		fmt.Printf("unable to decode message: %v, error: %v\n", pkt.Desc(), err)
		resp := &ErrResp{Code: 111, Msg: "server is busy"}
		// To ensure there must be a response to the client, we ignore the encoding result.
		c.Msg, _ = json.Marshal(resp)
		c.Pkt.(pkts.DataPkt).SetCmd(110)
		handleResp(c)
		c.Abort()
		return
	}
}

func handleResp(c *route.Context) {
	c.Pkt.SetBody(c.Msg.([]byte))
	c.Conn.Send(c.Pkt)
}

func handleNotifyMsg(c *route.Context) {
	msg := c.Msg.(*Message)
	fmt.Println("receive notification: ", msg.Content)
}

func handleSubMsg(c *route.Context) {
	// Send a response without body in the next handler.
	c.Msg = []byte{}
	// Set the topic property for the wresponse.
	mockToipc:="1000:2"
	c.Pkt.(*pkts.SubPkt).Props().WithStr(pkts.PropTopic, mockToipc)

	// Mock to publish later messages.
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			msg := &Message{Id: i, Content: fmt.Sprint("latest message ",i)}
			data, _ := json.Marshal(msg)
			pkt := pkts.NewEasyPubPkt(mockToipc,1000, data)
			c.Conn.Send(pkt)
		}

		panic("panic actively to test the closed connection")
	}()
}
