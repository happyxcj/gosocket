package protocol

import (
	"io"
)

type Reader struct {
	// srcR is the underlying io Reader.
	srcR  io.Reader
	order ByteOrder
	buf   []byte
	off   int
}

func NewReader(srcR io.Reader, order ByteOrder) *Reader {
	r := &Reader{
		srcR:  srcR,
		order: order,
	}
	return r
}

// Reset discards any unread buffered data,
// and resets w to read any future buffered data from srcR.
func (r *Reader) Reset(srcR io.Reader, buf []byte) {
	r.srcR = srcR
	r.ResetBuf(buf)
}

// Reset resets the current buffer.
func (r *Reader) ResetBuf(buf []byte) {
	r.buf = buf
	r.off = 0
}

// Unread returns the size of bytes that have been read in the buffer.
func (r *Reader) AlreadyRead() int { return r.off }

// Unread returns how many bytes are unread in the buffer.
func (r *Reader) Unread() int { return len(r.buf) - r.off }

// Read implements io.Reader.
// Read resets the buffer to p and reads data form the underlying io.Reader into the buffer.
// It returns the number of bytes read into the buffer.
func (r *Reader) Read(p []byte) (n int, err error) {
	r.ResetBuf(p)
	n, err = r.srcR.Read(r.buf)
	return
}

// ReadFull resets the buffer to p and
// reads exactly len(buf) bytes form the underlying io.Reader into buffer.
// It returns the number of bytes copied and an error if fewer bytes were read.
func (r *Reader) ReadFull(p []byte) (n int, err error) {
	r.ResetBuf(p)
	n, err = io.ReadFull(r.srcR, r.buf)
	return
}

// Discard skips the next n bytes, returning the number of bytes discarded.
//
// If Discard skips fewer than n bytes, it also returns an error.
// If 0 <= n <= r.Unread(), Discard is guaranteed to succeed without
// reading from the underlying io.Reader.
func (r *Reader) Discard(n int) (int, error) {
	if n <= 0 {
		return 0, nil
	}
	unread := r.Unread()
	if n > unread {
		return unread, io.ErrUnexpectedEOF
	}
	r.off += n
	return n, nil
}

// HasSize returns a boolean that indicates whether there is
// no less than given size bytes unread in the buffer.
func (r *Reader) HasSize(size int) bool {
	return r.Unread() >= size
}

func (r *Reader) ReadBytes(size int) ([]byte, error) {
	if !r.HasSize(size) {
		return nil, io.ErrUnexpectedEOF
	}
	return r.Bytes(size), nil
}

func (r *Reader) ReadString(size int) (string, error) {
	if !r.HasSize(size) {
		return "", io.ErrUnexpectedEOF
	}
	return r.String(size), nil
}

func (r *Reader) ReadByte() (byte, error) {
	if !r.HasSize(1) {
		return 0, io.ErrUnexpectedEOF
	}
	return r.Byte(), nil
}

func (r *Reader) ReadUint16() (uint16, error) {
	if !r.HasSize(2) {
		return 0, io.ErrUnexpectedEOF
	}
	return r.Uint16(), nil
}
func (r *Reader) ReadUint24() (uint32, error) {
	if !r.HasSize(3) {
		return 0, io.ErrUnexpectedEOF
	}
	return r.Uint24(), nil
}

func (r *Reader) ReadUint32() (uint32, error) {
	if !r.HasSize(4) {
		return 0, io.ErrUnexpectedEOF
	}
	return r.Uint32(), nil
}

func (r *Reader) ReadUint40() (uint64, error) {
	if !r.HasSize(5) {
		return 0, io.ErrUnexpectedEOF
	}
	return r.Uint40(), nil
}

func (r *Reader) ReadUint48() (uint64, error) {
	if !r.HasSize(6) {
		return 0, io.ErrUnexpectedEOF
	}
	return r.Uint48(), nil
}

func (r *Reader) ReadUint56() (uint64, error) {
	if !r.HasSize(7) {
		return 0, io.ErrUnexpectedEOF
	}
	return r.Uint56(), nil
}

func (r *Reader) ReadUint64() (uint64, error) {
	if !r.HasSize(8) {
		return 0, io.ErrUnexpectedEOF
	}
	return r.Uint64(), nil
}

func (r *Reader) ReadFloat32() (float32, error) {
	if !r.HasSize(4) {
		return 0, io.ErrUnexpectedEOF
	}
	return r.Float32(), nil
}

func (r *Reader) ReadFloat64() (float64, error) {
	if !r.HasSize(8) {
		return 0, io.ErrUnexpectedEOF
	}
	return r.Float64(), nil
}

func (r *Reader) RemainBytes() []byte {
	old := r.off
	r.off = len(r.buf)
	// Make a copy of the buffer.
	v := make([]byte, r.off-old)
	copy(v, r.buf[old:])
	return v
}

func (r *Reader) Bytes(size int) []byte {
	old := r.off
	r.off += size
	// Make a copy of the buffer.
	v := make([]byte, size)
	copy(v, r.buf[old:r.off])
	return v
}

func (r *Reader) String(size int) string {
	old := r.off
	r.off += size
	v := string(r.buf[old:r.off])
	return v
}

func (r *Reader) Byte() byte {
	v := r.buf[r.off]
	r.off += 1
	return v
}

func (r *Reader) Uint16() uint16 {
	old := r.off
	r.off += 2
	v := r.order.Uint16(r.buf[old:r.off])
	return v
}

func (r *Reader) Uint24() uint32 {
	old := r.off
	r.off += 3
	v := r.order.Uint24(r.buf[old:r.off])
	return v
}

func (r *Reader) Uint32() uint32 {
	old := r.off
	r.off += 4
	v := r.order.Uint32(r.buf[old:r.off])
	return v
}

func (r *Reader) Uint40() uint64 {
	old := r.off
	r.off += 5
	v := r.order.Uint40(r.buf[old:r.off])
	return v
}

func (r *Reader) Uint48() uint64 {
	old := r.off
	r.off += 6
	v := r.order.Uint48(r.buf[old:r.off])
	return v
}

func (r *Reader) Uint56() uint64 {
	old := r.off
	r.off += 7
	v := r.order.Uint56(r.buf[old:r.off])
	return v
}

func (r *Reader) Uint64() uint64 {
	old := r.off
	r.off += 8
	v := r.order.Uint64(r.buf[old:r.off])
	return v
}

func (r *Reader) Float32() float32 {
	old := r.off
	r.off += 4
	v := r.order.Float32(r.buf[old:r.off])
	return v
}

func (r *Reader) Float64() float64 {
	old := r.off
	r.off += 8
	v := r.order.Float64(r.buf[old:r.off])
	return v
}
