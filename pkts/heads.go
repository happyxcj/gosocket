package pkts

import (
	"github.com/happyxcj/gosocket/protocol"
)

const (
	CodecBits = 3
)

type NotifyHead struct {
	cmd     uint16
	version byte
	codec   byte
	props   *Props
}

func (h *NotifyHead) InitProps() {
	h.props = NewProps()
}

func (h *NotifyHead) HeadSize() int {
	// 2Bytes(Cmd)+1Byte(Version<<3|Codec)+2Bytes(props size)+xBytes(props)
	return 5 + h.props.Size()
}

func (h *NotifyHead) EncodeHead(w *protocol.Writer) {
	w.PutUint16(h.cmd)
	w.PutByte(h.version<<CodecBits | h.codec)
	w.PutUint16(uint16(h.props.Size()))
	h.props.Encode(w)
}

func (h *NotifyHead) DecodeHead(r *protocol.Reader) error {
	if !r.HasSize(5) {
		return protocol.ErrDecodeBadPacket
	}
	h.cmd = r.Uint16()
	vc := r.Byte()
	h.version = vc >> CodecBits
	h.codec = vc & 0x08
	size := int(r.Uint16())
	h.props.Reset()
	return h.props.Decode(r, size)
}

func (h *NotifyHead) Cmd() uint16 {
	return h.cmd
}

func (h *NotifyHead) SetCmd(cmd uint16) {
	h.cmd = cmd
}

func (h *NotifyHead) Version() byte {
	return h.version
}

func (h *NotifyHead) SetVersion(version byte) {
	h.version = version
}

func (h *NotifyHead) Codec() byte {
	return h.codec
}

func (h *NotifyHead) SetCodec(codec byte) {
	h.codec = codec
}

func (h *NotifyHead) Props() *Props {
	return h.props
}

func (h *NotifyHead) SetProps(props *Props) {
	h.props = props
}

type HttpHead struct {
	seqId uint16
	NotifyHead
}

func (h *HttpHead) HeadSize() int {
	// 2Bytes(seqId)+ ...
	return 2 + h.NotifyHead.HeadSize()
}

func (h *HttpHead) EncodeHead(w *protocol.Writer) {
	w.PutUint16(h.seqId)
	h.NotifyHead.EncodeHead(w)
}

func (h *HttpHead) DecodeHead(r *protocol.Reader) error {
	if !r.HasSize(2) {
		return protocol.ErrDecodeBadPacket
	}
	h.seqId = r.Uint16()
	return h.NotifyHead.DecodeHead(r)
}

func (h *HttpHead) SeqId() uint16 {
	return h.seqId
}

func (h *HttpHead) SetSeqId(seqId uint16) {
	h.seqId = seqId
}

type PubHead struct {
	topic string
	NotifyHead
}

func (h *PubHead) HeadSize() int {
	// 2Bytes(topic size)+xBytes(topic)+ ...
	return 2 + len(h.topic) + h.NotifyHead.HeadSize()
}

func (h *PubHead) EncodeHead(w *protocol.Writer) {
	w.PutUint16(uint16(len(h.topic)))
	w.PutString(h.topic)
	h.NotifyHead.EncodeHead(w)
}

func (h *PubHead) DecodeHead(r *protocol.Reader) error {
	if !r.HasSize(2) {
		return protocol.ErrDecodeBadPacket
	}
	size := int(r.Uint16())
	if !r.HasSize(size) {
		return protocol.ErrDecodeBadPacket
	}
	h.topic = r.String(size)
	return h.NotifyHead.DecodeHead(r)
}

func (h *PubHead) Topic() string {
	return h.topic
}

func (h *PubHead) SetTopic(topic string) {
	h.topic = topic
}
