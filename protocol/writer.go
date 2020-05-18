package protocol

import (
	"io"
)

type Writer struct {
	// srcW is the underlying io Writer.
	srcW  io.Writer
	order ByteOrder
	// off is the start off for next write.
	off       int
	buf       []byte
}

func NewWriter(srcW io.Writer, order ByteOrder) *Writer {
	w := &Writer{
		srcW:  srcW,
		order: order,
	}
	return w
}

// Reset resets the current buffer to p,
// and resets w to write any future buffered data to srcW.
func (w *Writer) Reset(srcW io.Writer, p []byte) {
	w.srcW = srcW
	w.ResetBuf(p)
}

// Reset resets the current buffer to p.
func (w *Writer) ResetBuf(p []byte) {
	w.buf =p
	w.off = 0
}

// Available returns how many bytes are unused in the buffer.
func (w *Writer) Available() int { return len(w.buf) - w.off }

// Buffered returns the number of bytes that have been written into the current buffer.
func (w *Writer) Buffered() int { return w.off }

// Write implements io.Writer.
// Write writes the contents of p into the buffer.
// The return value n is the length of p; err is always nil.
func (w *Writer) Write(p []byte) (n int, err error) {
	w.PutBytes(p)
	return len(p), nil
}

// Flush writes the current buffered data to the underlying io.Writer.
func (w *Writer) Flush() (int, error) {
	return w.srcW.Write(w.buf[:])
}

func (w *Writer) PutByte(v byte) {
	w.buf[w.off] = v
	w.off++
}

func (w *Writer) PutString(v string) {
	old := w.off
	w.off += len(v)
	copy(w.buf[old:w.off], v)
}

func (w *Writer) PutBytes(v []byte) {
	old := w.off
	w.off += len(v)
	copy(w.buf[old:w.off], v)
}

func (w *Writer) PutUint16(v uint16) {
	old := w.off
	w.off += 2
	w.order.PutUint16(w.buf[old:w.off], v)
}

func (w *Writer) PutUint24(v uint32) {
	old := w.off
	w.off += 3
	w.order.PutUint24(w.buf[old:w.off], v)
}

func (w *Writer) PutUint32(v uint32) {
	old := w.off
	w.off += 4
	w.order.PutUint32(w.buf[old:w.off], v)
}

func (w *Writer) PutUint40(v uint64) {
	old := w.off
	w.off += 5
	w.order.PutUint40(w.buf[old:w.off], v)
}

func (w *Writer) PutUint48(v uint64) {
	old := w.off
	w.off += 6
	w.order.PutUint48(w.buf[old:w.off], v)
}

func (w *Writer) PutUint56(v uint64) {
	old := w.off
	w.off += 7
	w.order.PutUint56(w.buf[old:w.off], v)
}

func (w *Writer) PutUint64(v uint64) {
	old := w.off
	w.off += 8
	w.order.PutUint64(w.buf[old:w.off], v)
}

func (w *Writer) PutFloat32(v float32) {
	old := w.off
	w.off += 4
	w.order.PutFloat32(w.buf[old:w.off], v)
}

func (w *Writer) PutFloat64(v float64) {
	old := w.off
	w.off += 8
	w.order.PutFloat64(w.buf[old:w.off], v)
}
