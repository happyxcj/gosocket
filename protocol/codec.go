package protocol

const (
	fixedHeadLen = 3

	flagsBits = 4

	maxPktSize = 2<<16 - 1
)


// Codec is used to write and read packets.
// Be careful that it does not support concurrent 'Write' and 'Read'.
type Codec struct {
	w          *Writer
	wBuf       []byte
	wBufGetter func(size int) []byte

	r          *Reader
	rBuf       []byte
	rBufGetter func(size int) []byte
	// headBuf only used for buffer the fixed head of a packet.
	headBuf [fixedHeadLen]byte

	// maxPktSize is the maximum payload (variable head and body) size of the packet.
	// It can't be greater than "2<<16".
	maxPktSize int
}

type CodecOpt func(*Codec)

// MaxPktSize returns a CodecOpt to set the maximum payload size of a packet.
func MaxPktSize(size int) CodecOpt {
	return func(c *Codec) {
		if size > maxPktSize || size <= 0 {
			c.maxPktSize = maxPktSize
		}
		c.maxPktSize = size
	}
}

// WBufGetter returns a CodecOpt to set write buffer getter.
func WBufGetter(getter func(size int) []byte) CodecOpt {
	return func(c *Codec) {
		c.wBufGetter = getter
	}
}

// RBufGetter returns a CodecOpt to set read buffer getter.
func RBufGetter(getter func(size int) []byte) CodecOpt {
	return func(c *Codec) {
		c.wBufGetter = getter
	}
}

func NewCodec(w *Writer, r *Reader,opts ...CodecOpt) *Codec {
	c := &Codec{
		w: w,
		r: r,
	}
	c.maxPktSize = maxPktSize
	c.wBufGetter = func(size int) []byte {
		if cap(c.wBuf) >= size {
			return c.wBuf[:size]
		}
		c.wBuf = make([]byte, size)
		return c.wBuf
	}
	c.rBufGetter = func(size int) []byte {
		if cap(c.rBuf) >= size {
			return c.rBuf[:size]
		}
		c.rBuf = make([]byte, size)
		return c.rBuf
	}
	for _,opt:=range opts{
		opt(c)
	}
	return c
}

func (c *Codec) Write(p Packet) error {
	size := p.HeadSize() + len(p.Body())
	if size > c.maxPktSize {
		return ErrPacketTooLarge
	}
	c.w.ResetBuf(c.wBufGetter(size + fixedHeadLen))
	// 1. Write fixed head.
	c.w.PutByte(byte(p.Kind())<<flagsBits | byte(p.Flags()))
	c.w.PutUint16(uint16(size))
	// 2. Write variable head.
	p.EncodeHead(c.w)
	// 3. Write application message.
	if len(p.Body()) > 0 {
		c.w.PutBytes(p.Body())
	}
	_, err := c.w.Flush()
	return err
}

func (c *Codec) Read() (Packet, error) {
	var p Packet
	var err error
	// Read fixed head.
	if _, err = c.r.ReadFull(c.headBuf[:]); err != nil {
		return nil, err
	}
	// Decode fixed head.
	kindFlags := c.r.Byte()
	kind := PktKind(kindFlags >> flagsBits)
	flags := PktFlags(kindFlags & 0x0f)
	if p, err = FindPacket(kind, flags); err != nil {
		return nil, err
	}
	remainingSize := int(c.r.Uint16())
	if remainingSize > c.maxPktSize {
		return nil, ErrPacketTooLarge
	}
	// Read remaining data, it maybe contains variable head and application message.
	if _, err = c.r.ReadFull(c.rBufGetter(remainingSize)); err != nil {
		return nil, err
	}
	// Decode variable head.
	if err = p.DecodeHead(c.r); err != nil {
		return nil, err
	}
	// Decode application message.
	p.SetBody(c.r.RemainBytes())
	return p, nil
}
