package pkts

import (
	"errors"
	"github.com/happyxcj/gosocket/protocol"
	"fmt"
)

// packet property ids
const (
	PropCreatedTime = 1
	PropSubId       = 2
	PropTopic       = 3
)

var (
	ErrInvalidPropId = errors.New("invalid property id")

	propCreators = make(map[PropID]PropCreator)
)

type PropID byte

type PropCreator func() Prop

func init() {
	RegisterPropCreator(PropCreatedTime, func() Prop { return new(Uint64Prop) })
	RegisterPropCreator(PropSubId, func() Prop { return new(Uint16Prop) })
	RegisterPropCreator(PropTopic, func() Prop { return new(StringProp) })
}

// RegisterPropCreator registers a specified property creator based on the id.
// It will panic if the given id had exist.
func RegisterPropCreator(id PropID, creator func() Prop) {
	if _, ok := propCreators[id]; ok {
		panic(fmt.Sprintf("the given id '%v' had exist", id))
	}
	propCreators[id] = creator
}

// DeletePropCreators unregisters the property creators based on the given ids
// or all property creators if len(ids) is zero.
func DeletePropCreators(ids ...PropID) {
	if len(ids) == 0 {
		propCreators = make(map[PropID]PropCreator)
		return
	}
	for _, id := range ids {
		delete(propCreators, id)
	}
}

// FindProp returns a new property based on the id,
// if the id is invalid, it returns an error "ErrInvalidPropId".
func FindProp(id PropID) (Prop, error) {
	if creator, ok := propCreators[id]; ok {
		return creator(), nil
	}
	return nil, ErrInvalidPropId
}

type Props struct {
	props map[PropID]Prop
	len   int
}

func NewProps() *Props {
	return &Props{
		props: make(map[PropID]Prop),
	}
}

func (p *Props) Reset() {
	p.props = make(map[PropID]Prop)
	p.len = 0
}

func (p *Props) ForEach(f func(PropID, Prop)) {
	for id, prop := range p.props {
		f(id, prop)
	}
}

// With adds the given prop based on the given id.
//
// If a prop with the same id already exists in the props,
// replace it with the given prop.
func (p *Props) With(id PropID, prop Prop) *Props {
	// Check to remove the expired prop.
	p.Remove(id)
	p.props[id] = prop
	p.len += 1 + prop.Size()
	return p
}

// Get gets the specified prop based on the given id if it exists.
func (p *Props) Get(id PropID) (Prop, bool) {
	v, ok := p.props[id]
	return v, ok
}

// Remove removes the specified prop based on the given id.
func (p *Props) Remove(id PropID) *Props {
	if prop, ok := p.props[id]; ok {
		delete(p.props, id)
		p.len -= 1 + prop.Size()
	}
	return p
}

// Has returns a boolean that indicates whether there is the specified prop
// based on the given id in the p.
func (p *Props) Has(id PropID) bool {
	_, ok := p.props[id]
	return ok
}

// Size returns the size of all property key-value pairs.
// It does not include the size of properties length prefix.
func (p *Props) Size() int {
	return p.len
}

func (p *Props) Encode(w *protocol.Writer) {
	for id, prop := range p.props {
		w.PutByte(byte(id))
		prop.Encode(w)
	}
}

func (p *Props) Decode(r *protocol.Reader, size int) error {
	if size == 0 {
		// Without properties.
		return nil
	}
	// Make sure that there are enough bytes in the buffer to read all property ids
	// in the next for loop.
	if !r.HasSize(size) {
		return protocol.ErrDecodeBadPacket
	}
	var id PropID
	var prop Prop
	var err error
	tmp := r.AlreadyRead()
	end := tmp + size
	for r.AlreadyRead() < end {
		id = PropID(r.Byte())
		if prop, err = FindProp(id); err != nil {
			return err
		}
		err = prop.Decode(r)
		if err != nil {
			return err
		}
		p.props[id] = prop
	}
	// real len
	p.len = r.AlreadyRead() - tmp
	return nil
}

// WithInt8 adds the given value based on the given id.
func (p *Props) WithInt8(id PropID, val int8) *Props {
	return p.WithUint8(id, uint8(val))
}

// GetInt8 gets the specified value based on the given id if it exists.
func (p *Props) GetInt8(id PropID) (int8, bool) {
	v, ok := p.GetUint8(id)
	return int8(v), ok
}

// WithInt16 adds the given value based on the given id.
func (p *Props) WithInt16(id PropID, val int16) *Props {
	return p.WithUint16(id, uint16(val))
}

// GetInt16 gets the specified value based on the given id if it exists.
func (p *Props) GetInt16(id PropID) (int16, bool) {
	v, ok := p.GetUint16(id)
	return int16(v), ok
}

// WithInt32 adds the given value based on the given id.
func (p *Props) WithInt32(id PropID, val int32) *Props {
	return p.WithUint32(id, uint32(val))
}

// GetInt8 gets the specified value based on the given id if it exists.
func (p *Props) GetInt32(id PropID) (int32, bool) {
	v, ok := p.GetUint32(id)
	return int32(v), ok
}

// WithInt64 adds the given value based on the given id.
func (p *Props) WithInt64(id PropID, val int64) *Props {
	return p.WithUint64(id, uint64(val))
}

// GetInt64 gets the specified value based on the given id if it exists.
func (p *Props) GetInt64(id PropID) (int64, bool) {
	v, ok := p.GetUint64(id)
	return int64(v), ok
}

// WithByte adds the given value based on the given id.
func (p *Props) WithByte(id PropID, val byte) *Props {
	return p.WithUint8(id, val)
}

// GetByte gets the specified value based on the given id if it exists.
func (p *Props) GetByte(id PropID) (byte, bool) {
	return p.GetUint8(id)
}

// WithUint8 adds the given value based on the given id.
func (p *Props) WithUint8(id PropID, val uint8) *Props {
	p.With(id, &Uint8Prop{val})
	return p
}

// GetUint8 gets the specified value based on the given id if it exists.
func (p *Props) GetUint8(id PropID) (uint8, bool) {
	v, ok := p.Get(id)
	if !ok {
		return 0, false
	}
	return v.(*Uint8Prop).Val, true
}

// WithUint16 adds the given value based on the given id.
func (p *Props) WithUint16(id PropID, val uint16) *Props {
	p.With(id, &Uint16Prop{val})
	return p
}

// GetUint16 gets the specified value based on the given id if it exists.
func (p *Props) GetUint16(id PropID) (uint16, bool) {
	v, ok := p.Get(id)
	if !ok {
		return 0, false
	}
	return v.(*Uint16Prop).Val, true
}

// WithUint32 adds the given value based on the given id.
func (p *Props) WithUint32(id PropID, val uint32) *Props {
	p.With(id, &Uint32Prop{val})
	return p
}

// GetUint32 gets the specified value based on the given id if it exists.
func (p *Props) GetUint32(id PropID) (uint32, bool) {
	v, ok := p.Get(id)
	if !ok {
		return 0, false
	}
	return v.(*Uint32Prop).Val, true
}

// WithUint64 adds the given value based on the given id.
func (p *Props) WithUint64(id PropID, val uint64) *Props {
	p.With(id, &Uint64Prop{val})
	return p
}

// GetUint64 gets the specified value based on the given id if it exists.
func (p *Props) GetUint64(id PropID) (uint64, bool) {
	v, ok := p.Get(id)
	if !ok {
		return 0, false
	}
	return v.(*Uint64Prop).Val, true
}

// WithFloat32 adds the given value based on the given id.
func (p *Props) WithFloat32(id PropID, val float32) *Props {
	p.With(id, &Float32Prop{val})
	return p
}

// GetFloat32 gets the specified value based on the given id if it exists.
func (p *Props) GetFloat32(id PropID) (float32, bool) {
	v, ok := p.Get(id)
	if !ok {
		return 0, false
	}
	return v.(*Float32Prop).Val, true
}

// WithFloat64 adds the given value based on the given id.
func (p *Props) WithFloat64(id PropID, val float64) *Props {
	p.With(id, &Float64Prop{val})
	return p
}

// GetFloat64 gets the specified value based on the given id if it exists.
func (p *Props) GetFloat64(id PropID) (float64, bool) {
	v, ok := p.Get(id)
	if !ok {
		return 0, false
	}
	return v.(*Float64Prop).Val, true
}

// WithStr adds the given value based on the given id.
func (p *Props) WithStr(id PropID, val string) *Props {
	p.With(id, &StringProp{val})
	return p
}

// GetStr gets the specified value based on the given id if it exists.
func (p *Props) GetStr(id PropID) (string, bool) {
	v, ok := p.Get(id)
	if !ok {
		return "", false
	}
	return v.(*StringProp).Val, true
}

type Prop interface {
	// Size returns the size of the property.
	Size() int
	// Encode writes the property to w.
	Encode(w *protocol.Writer)
	// Decode reads the property form the r.
	Decode(r *protocol.Reader) error
}

type Uint8Prop struct {
	Val uint8
}

func (p *Uint8Prop) Size() int {
	return 1
}

func (p *Uint8Prop) Encode(w *protocol.Writer) {
	w.PutByte(p.Val)
}

func (p *Uint8Prop) Decode(r *protocol.Reader) (err error) {
	if !r.HasSize(1) {
		return protocol.ErrDecodeBadPacket
	}
	p.Val = r.Byte()
	return nil
}

type Uint16Prop struct {
	Val uint16
}

func (p *Uint16Prop) Size() int {
	return 2
}

func (p *Uint16Prop) Encode(w *protocol.Writer) {
	w.PutUint16(p.Val)
}

func (p *Uint16Prop) Decode(r *protocol.Reader) error {
	if !r.HasSize(2) {
		return protocol.ErrDecodeBadPacket
	}
	p.Val = r.Uint16()
	return nil
}

type Uint32Prop struct {
	Val uint32
}

func (p *Uint32Prop) Size() int {
	return 4
}

func (p *Uint32Prop) Encode(w *protocol.Writer) {
	w.PutUint32(p.Val)
}

func (p *Uint32Prop) Decode(r *protocol.Reader) error {
	if !r.HasSize(4) {
		return protocol.ErrDecodeBadPacket
	}
	p.Val = r.Uint32()
	return nil
}

type Uint64Prop struct {
	Val uint64
}

func (p *Uint64Prop) Size() int {
	return 8
}

func (p *Uint64Prop) Encode(w *protocol.Writer) {
	w.PutUint64(p.Val)
}

func (p *Uint64Prop) Decode(r *protocol.Reader) (err error) {
	if !r.HasSize(8) {
		return protocol.ErrDecodeBadPacket
	}
	p.Val = r.Uint64()
	return nil
}

type Float32Prop struct {
	Val float32
}

func (p *Float32Prop) Size() int {
	return 4
}

func (p *Float32Prop) Encode(w *protocol.Writer) {
	w.PutFloat32(p.Val)
}

func (p *Float32Prop) Decode(r *protocol.Reader) error {
	if !r.HasSize(4) {
		return protocol.ErrDecodeBadPacket
	}
	p.Val = r.Float32()
	return nil
}

type Float64Prop struct {
	Val float64
}

func (p *Float64Prop) Size() int {
	return 8
}

func (p *Float64Prop) Encode(w *protocol.Writer) {
	w.PutFloat64(p.Val)
}

func (p *Float64Prop) Decode(r *protocol.Reader) (err error) {
	if !r.HasSize(8) {
		return protocol.ErrDecodeBadPacket
	}
	p.Val = r.Float64()
	return nil
}

type StringProp struct {
	Val string
}

func (p *StringProp) Size() int {
	return 2 + len(p.Val)
}

func (p *StringProp) Encode(w *protocol.Writer) {
	w.PutUint16(uint16(len(p.Val)))
	w.PutString(p.Val)
}

func (p *StringProp) Decode(r *protocol.Reader) (err error) {
	if !r.HasSize(2) {
		return protocol.ErrDecodeBadPacket
	}
	size := int(r.Uint16())
	if !r.HasSize(size) {
		return protocol.ErrDecodeBadPacket
	}
	p.Val = r.String(size)
	return nil
}
