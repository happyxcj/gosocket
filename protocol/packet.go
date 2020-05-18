package protocol

import "fmt"


var pktCreators = make(map[PktKind]PktCreator)

type PktCreator func(b *PktBase) Packet

// RegisterPktCreator registers a specified packet creator based on the kind.
// It will panic if the given kind had exist.
func RegisterPktCreator(kind PktKind, creator PktCreator) {
	if _, ok := pktCreators[kind]; ok {
		panic(fmt.Sprintf("the given kind '%v' had exist", kind))
	}
	pktCreators[kind] = creator
}

// DeletePktCreators deletes the packet creators based on the given kinds
// or all packet creators if len(kinds) is zero.
func DeletePktCreators(kinds ...PktKind) {
	if len(kinds) == 0 {
		pktCreators = make(map[PktKind]PktCreator)
		return
	}
	for _, kind := range kinds {
		delete(pktCreators, kind)
	}
}

// CreatePacket returns a new packet based on the packet kind and flags,
// if the packet kind is invalid, it returns an error "ErrInvalidPktKind".
func FindPacket(kind PktKind, flags PktFlags) (Packet, error) {
	if creator, ok := pktCreators[kind]; ok {
		return creator(NewPktBase(kind, flags)), nil
	}
	return nil, ErrInvalidPktKind
}

// PktKind is an enum of packet kind.
type PktKind byte

type PktFlags byte

// Has indicates whether f contains specified (0 or more) flags in v.
func (f PktFlags) Has(v PktFlags) bool {
	return (f & v) == v
}

type Packet interface {
	// Desc returns a string description of the packet.
	Desc() string
	// Kind returns the kind of the packet.
	Kind() PktKind
	// Kind returns the flags of the packet.
	Flags() PktFlags
	// Size returns the size of the variable head.
	HeadSize() int
	// EncodeHead writes the variable head to w.
	EncodeHead(w *Writer)
	// Decode reads the variable head form the r.
	DecodeHead(r *Reader) error
	// Body returns the application message.
	Body() []byte
	// SetBody sets the application message.
	SetBody(body []byte)
}

type PktBase struct {
	kind  PktKind
	flags PktFlags
	body  []byte
}

func NewPktBase(kind PktKind, flags PktFlags) *PktBase {
	return &PktBase{
		kind:  kind,
		flags: flags,
	}
}

func (b *PktBase) Desc() string {
	return fmt.Sprintf("kind=%v",b.kind)
}

func (b *PktBase) Kind() PktKind {
	return b.kind
}

func (b *PktBase) Flags() PktFlags {
	return b.flags
}

func (b *PktBase) Body() []byte {
	return b.body
}

func (b *PktBase) SetBody(body []byte) {
	b.body = body
}

// AddFlags adds the given flags to the h.
func (b *PktBase) AddFlags(flags ...PktFlags) {
	for _, f := range flags {
		b.flags |= f
	}
}
