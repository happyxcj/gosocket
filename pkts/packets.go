package pkts

import (
	"github.com/happyxcj/gosocket/protocol"
	"fmt"
)

// packet kinds
const (
	KindPing        protocol.PktKind = iota + 1
	KindNotify
	KindHttp
	KindSubscribe
	KindUnsubscribe
	KindPublish
)

// packet flags
const (
	FlagNo   protocol.PktFlags = 0
	FlagPong protocol.PktFlags = 0x01
)


func init() {
	protocol.RegisterPktCreator(KindPing, func(b *protocol.PktBase) protocol.Packet {
		pkt := NewPingPkt()
		pkt.PktBase = b
		return pkt
	})
	protocol.RegisterPktCreator(KindNotify, func(b *protocol.PktBase) protocol.Packet {
		return NewNotifyPkt(b)
	})
	protocol.RegisterPktCreator(KindHttp, func(b *protocol.PktBase) protocol.Packet {
		return NewHttpPkt(b)
	})
	protocol.RegisterPktCreator(KindSubscribe, func(b *protocol.PktBase) protocol.Packet {
		return NewSubPkt(b)
	})
	protocol.RegisterPktCreator(KindUnsubscribe, func(b *protocol.PktBase) protocol.Packet {
		return NewUnsubPkt(b)
	})
	protocol.RegisterPktCreator(KindPublish, func(b *protocol.PktBase) protocol.Packet {
		return NewPubPkt(b)
	})
}

// DataPkt represents a packet that has the application message.
type DataPkt interface {
	protocol.Packet
	// Cmd returns the unique message identifier of the data.
	// It decides the body structure together with the version and codec.
	Cmd() uint16
	// SetCmd sets the unique message identifier of the data.
	SetCmd(cmd uint16)
	// Version returns the version of the data.
	Version() byte
	// SetCmd sets the the version of the data.
	SetVersion(version byte)
	// Codec returns the codec mode of the data.
	Codec() byte
	// SetCodec sets the codec mode of the data.
	SetCodec(codec byte)
	// Props returns the optional properties.
	Props() *Props
	// Props sets the optional properties.
	SetProps(props *Props)
}

// ReqRespPkt represents the type of data packet:
// For a client request packet, the server must have a corresponding response packet,
// Both the request packet and response packet have the same unique sequence id.
type ReqRespPkt interface {
	DataPkt
	// SeqId returns the unique sequence id of the packet.
	SeqId() uint16
	// SetSeqId sets the unique sequence id for the packet.
	SetSeqId(seqId uint16)
}

var _ protocol.Packet = (*PingPkt)(nil)

type PingPkt struct {
	*protocol.PktBase
}

func NewPingPkt() *PingPkt {
	return &PingPkt{PktBase: protocol.NewPktBase(KindPing, FlagNo)}
}

func NewPongPkt() *PingPkt {
	return &PingPkt{PktBase: protocol.NewPktBase(KindPing, FlagPong)}
}

func (p *PingPkt) Desc() string {
	if p.Flags().Has(FlagPong) {
		return "Pong"
	}
	return "Ping"
}

func (p *PingPkt) HeadSize() int {
	return 0
}

func (p *PingPkt) EncodeHead(w *protocol.Writer) {
}

func (p *PingPkt) DecodeHead(r *protocol.Reader) error {
	return nil
}

var _ DataPkt = (*NotifyPkt)(nil)

type NotifyPkt struct {
	*protocol.PktBase
	NotifyHead
}

func NewNotifyPkt(b *protocol.PktBase) *NotifyPkt {
	p := &NotifyPkt{PktBase: b}
	p.InitProps()
	return p
}

func NewEasyNotifyPkt(cmd uint16, body []byte) *NotifyPkt {
	p := NewNotifyPkt(protocol.NewPktBase(KindNotify, FlagNo))
	p.cmd = cmd
	p.SetBody(body)
	return p
}

func NewFullNotifyPkt(cmd uint16, version byte, codec byte, body []byte) *NotifyPkt {
	p := NewEasyNotifyPkt(cmd, body)
	p.version = version
	p.codec = codec
	return p
}

func (p *NotifyPkt) Desc() string {
	// format: "Notify:cmd:version"
	return fmt.Sprintf("Notify:%v:%v", p.cmd, p.version)
}

var _ ReqRespPkt = (*HttpPkt)(nil)

type HttpPkt struct {
	*protocol.PktBase
	HttpHead
}

func NewHttpPkt(b *protocol.PktBase) *HttpPkt {
	p := &HttpPkt{PktBase: b}
	p.InitProps()
	return p
}

func NewEasyHttpPkt(seqId, cmd uint16, body []byte) *HttpPkt {
	p := NewHttpPkt(protocol.NewPktBase(KindHttp, FlagNo))
	p.seqId = seqId
	p.cmd = cmd
	p.SetBody(body)
	return p
}

func NewFullHttpPkt(seqId, cmd uint16, version, codec byte, body []byte) *HttpPkt {
	p := NewEasyHttpPkt(seqId, cmd, body)
	p.version = version
	p.codec = codec
	return p
}

func (p *HttpPkt) Desc() string {
	return fmt.Sprintf("Http:%v:%v", p.cmd, p.version)
}

var _ ReqRespPkt = (*SubPkt)(nil)

type SubPkt struct {
	*protocol.PktBase
	HttpHead
}

func NewSubPkt(b *protocol.PktBase) *SubPkt {
	p := &SubPkt{PktBase: b}
	p.InitProps()
	return p
}

func NewEasySubPkt(seqId, cmd uint16, body []byte) *SubPkt {
	p := NewSubPkt(protocol.NewPktBase(KindSubscribe, FlagNo))
	p.seqId = seqId
	p.cmd = cmd
	p.SetBody(body)
	return p
}

func NewFullSubPkt(seqId, cmd uint16, version, codec byte, body []byte) *SubPkt {
	p := NewEasySubPkt(seqId, cmd, body)
	p.version = version
	p.codec = codec
	return p
}

func (p *SubPkt) Desc() string {
	return fmt.Sprintf("Subscribe:%v:%v", p.cmd, p.version)
}

var _ ReqRespPkt = (*UnsubPkt)(nil)

type UnsubPkt struct {
	*protocol.PktBase
	HttpHead
}

func NewUnsubPkt(b *protocol.PktBase) *UnsubPkt {
	p := &UnsubPkt{PktBase: b}
	p.InitProps()
	return p
}

func NewEasyUnsubPkt(seqId, cmd uint16, body []byte) *UnsubPkt {
	p := NewUnsubPkt(protocol.NewPktBase(KindUnsubscribe, FlagNo))
	p.seqId = seqId
	p.cmd = cmd
	p.SetBody(body)
	return p
}

func NewFullUnsubPkt(seqId, cmd uint16, version, codec byte, body []byte) *UnsubPkt {
	p := NewEasyUnsubPkt(seqId, cmd, body)
	p.version = version
	p.codec = codec
	return p
}

func (p *UnsubPkt) Desc() string {
	return fmt.Sprintf("Unsubscribe:%v:%v", p.cmd, p.version)
}

var _ DataPkt = (*PubPkt)(nil)

type PubPkt struct {
	*protocol.PktBase
	PubHead
}

func NewPubPkt(b *protocol.PktBase) *PubPkt {
	p := &PubPkt{PktBase: b}
	p.InitProps()
	return p
}

func NewEasyPubPkt(topic string, cmd uint16, body []byte) *PubPkt {
	p := NewPubPkt(protocol.NewPktBase(KindPublish, FlagNo))
	p.topic = topic
	p.cmd = cmd
	p.SetBody(body)
	return p
}

func NewFullPubPkt(topic string, cmd uint16, version, codec byte, body []byte) *PubPkt {
	p := NewEasyPubPkt(topic, cmd, body)
	p.version = version
	p.codec = codec
	return p
}

func (p *PubPkt) Desc() string {
	return fmt.Sprintf("Publish:%v:%v", p.cmd, p.version)
}
