/* TODO: Waiting To be completed and tested.
 *
 *
 */
package gosocket
//
//import (
//	"time"
//	"net"
//	"errors"
//	"unsafe"
//	"github.com/happyxcj/gosocket/pkts"
//	"github.com/happyxcj/gosocket/protocol"
//	"sync/atomic"
//)
//
//const (
//	defaultRespTimeout = 10 * time.Second
//	defaultDialTimeout = 5 * time.Second
//	//defaultWriteTimeout = 15 * time.Second
//	//defaultReadTimeout  = 20 * time.Second
//	//defaultSendChSize   = 200
//	maxConcurrentSeqId = 1 << 16
//	respChSize         = 1
//)
//
//var (
//	ErrTimeout     = errors.New("request timed out")
//	ErrPktDisorder = errors.New("the packet responded by the server is disordered")
//)
//
//// Client represents a client connection to an specified server.
//type Client struct {
//	// opts contains the options to dial a server.
//	opts      *DialOpts
//	autoSeqId uint64
//	seqIdCh   chan uint16
//	respChs   [maxConcurrentSeqId]*unsafe.Pointer
//	pool      ClientConnPool
//}
//
//// baseResp is used to deliver a packet response or an error response
//// to the corresponding request.
//type baseResp struct {
//	pkt pkts.ReqRespPkt
//	err error
//}
//
//type DialOpts struct {
//	ecOpts      []EasyConnOpt
//	qcOpts      []QueueConnOpt
//	dialTimeout time.Duration
//	// respTimeout specifies the timeout to wait for a server's response.
//	// It's default value is "10*time.second".
//	respTimeout time.Duration
//	dialer      MyDialer
//}
//
//// DialOpt specifies an option for a connection.
//type DialOpt func(*DialOpts)
//
//// ECOption returns a DialOpt to add a option to the internal ecOpts.
//func ECOption(opt EasyConnOpt) DialOpt {
//	return func(o *DialOpts) {
//		o.ecOpts = append(o.ecOpts, opt)
//	}
//}
//
//// QCOption returns a DialOpt to add a option to the internal qcOpts.
//func QCOption(opt QueueConnOpt) DialOpt {
//	return func(o *DialOpts) {
//		o.qcOpts = append(o.qcOpts, opt)
//	}
//}
//
//func WrapPktHandler()  DialOpt{
//
//}
//
//// RespTimeout returns a DialOpt to set the timeout to wait for a server's response.
//func RespTimeout(t time.Duration) DialOpt {
//	return func(o *DialOpts) {
//		o.dialTimeout = t
//	}
//}
//
//// DialTimeout returns a DialOpt to set the dial timeout for connecting to the server.
//func DialTimeout(t time.Duration) DialOpt {
//	return func(o *DialOpts) {
//		o.dialTimeout = t
//	}
//}
//
//// Dialer returns a DialOpt to set the dialer used to connect to the server.
//func Dialer(dialer MyDialer) DialOpt {
//	return func(o *DialOpts) {
//		o.dialer = dialer
//	}
//}
//
//func NewClient(addr string, opts ... DialOpt) *Client {
//	c := &Client{
//		seqIdCh: make(chan uint16, 10),
//	}
//	c.opts = &DialOpts{
//		respTimeout: defaultRespTimeout,
//		dialTimeout: defaultDialTimeout,
//		dialer:      &net.Dialer{Timeout: defaultDialTimeout},
//	}
//	for _, opt := range opts {
//		opt(c.opts)
//	}
//	for i := 0; i < maxConcurrentSeqId; i++ {
//		var addr *uint16
//		c.respChs[i] = (*unsafe.Pointer)(unsafe.Pointer(&addr))
//	}
//	connCreator := func(nc net.Conn) Conn {
//		inner := NewEasyConn(nc, c.opts.ecOpts...)
//		conn := NewQueueConn(inner, c.opts.qcOpts...)
//		return conn
//	}
//	c.pool = newClientConnPool(addr, c.opts.dialer, connCreator)
//	return c
//}
//
//// NewAndInitClient returns a Client and initializes to connect to the server.
//func NewAndInitClient(addr string, opts ... DialOpt) (*Client, error) {
//	c := NewClient(addr, opts...)
//	_, err := c.pool.GetConn()
//	if err != nil {
//		return nil, err
//	}
//	return c, nil
//}
//
//func WrapPktHandler(h func(p protocol.Packet)) func(p protocol.Packet) {
//
//}
//
//func (c *Client) handlePkt(p protocol.Packet) {
//	pkt, ok := p.(pkts.ReqRespPkt)
//	if !ok {
//		if c.opts.pktHandler != nil {
//			c.opts.pktHandler(p)
//		}
//		return
//	}
//	respCh, ok := c.getRespCh(pkt.SeqId())
//	if !ok {
//		return
//	}
//	// Release the sequence id for the possible waiters.
//	c.releaseSeqId(pkt.SeqId())
//	// Deliver the response to the waiter.
//	respCh <- &baseResp{pkt: pkt, err: nil}
//}
//
//// Send sends a packet to the server.
//func (c *Client) Send(p protocol.Packet) error {
//	conn, err := c.pool.GetConn()
//	if err != nil {
//		return err
//	}
//	return conn.Send(p)
//}
//
//// Http will send a 'http' request.
//// It returns the corresponding response or an error synchronously.
//func (c *Client) Http(pkt *pkts.HttpPkt) (*pkts.HttpPkt, error) {
//	resp, err := c.doRequest(pkt)
//	if err != nil {
//		return nil, err
//	}
//	pkt, ok := resp.(*pkts.HttpPkt)
//	if !ok {
//		return nil, ErrPktDisorder
//	}
//	return pkt, nil
//}
//
//// Subscribe will send a 'subscribe' request.
//// It returns the corresponding response or an error synchronously.
//func (c *Client) Subscribe(pkt *pkts.SubPkt) (*pkts.SubPkt, error) {
//	resp, err := c.doRequest(pkt)
//	if err != nil {
//		return nil, err
//	}
//	pkt, ok := resp.(*pkts.SubPkt)
//	if !ok {
//		return nil, ErrPktDisorder
//	}
//	return pkt, nil
//}
//
//// Unsubscribe will send a 'unsubscribe' request.
//// It returns the corresponding response or an error synchronously.
//func (c *Client) Unsubscribe(pkt *pkts.UnsubPkt) (*pkts.UnsubPkt, error) {
//	resp, err := c.doRequest(pkt)
//	if err != nil {
//		return nil, err
//	}
//	pkt, ok := resp.(*pkts.UnsubPkt)
//	if !ok {
//		return nil, ErrPktDisorder
//	}
//	return pkt, nil
//}
//
//// Request will send a common request.
//// It returns the corresponding response or an error synchronously.
//func (c *Client) Request(pkt pkts.ReqRespPkt) (pkts.ReqRespPkt, error) {
//	return c.doRequest(pkt)
//}
//
//// doRequest is the internal common request function that is used to
//// send a request and deliver the corresponding response.
//func (c *Client) doRequest(pkt pkts.ReqRespPkt) (pkts.ReqRespPkt, error) {
//	conn, err := c.pool.GetConn()
//	if err != nil {
//		return nil, err
//	}
//	t := acquireTimer(c.opts.respTimeout)
//	defer releaseTimer(t)
//	seqId, respCh, err := c.getSeqIdAndRespCh(t.C)
//	if err != nil {
//		return nil, err
//	}
//	// Replace the seqId to make it unique among the remote server.
//	pkt.SetSeqId(seqId)
//	err = conn.Send(pkt)
//	if err != nil {
//		return nil, err
//	}
//	select {
//	case resp, ok := <-respCh:
//		if !ok {
//			return nil, ErrConnClosed
//		}
//		if resp.pkt != nil {
//			// Reset the original sequence id for the response.
//			resp.pkt.SetSeqId(pkt.SeqId())
//		}
//		return resp.pkt, resp.err
//	case <-t.C:
//		c.releaseSeqId(seqId)
//		return nil, ErrTimeout
//	}
//}
//
//func (c *Client) getSeqIdAndRespCh(timeoutCh <-chan time.Time) (uint16, chan *baseResp, error) {
//	seqId := uint16(atomic.AddUint64(&c.autoSeqId, 1))
//	respCh := make(chan *baseResp, respChSize)
//	//// We simply think that the sequence id will not be reused in the range of 0~65535
//	//// to obtain the target sequence id efficiently.
//	//// TODO: the sequence id in the range of 0~65535 may be in use when
//	//// TODO: the auto-incrementing id overflows and start a new round of auto-increment.
//	//if autoId < maxConcurrentSeqId {
//	//	c.saveRespCh(seqId, respCh)
//	//	return seqId, respCh, nil
//	//}
//	if c.swapRespCh(seqId, respCh) {
//		// The sequence id is not in use.
//		return seqId, respCh, nil
//	}
//	// The sequence id is in use, wait for any sequence id in use to be released.
//	// TODO: Is it necessary to expand the maximum concurrent number of sequence idc.
//	select {
//	case seqId := <-c.seqIdCh:
//		// Replace the expired old value with the new value.
//		c.saveRespCh(seqId, respCh)
//		return seqId, respCh, nil
//	case <-timeoutCh:
//		return 0, nil, ErrTimeout
//	}
//}
//
//// releaseSeqId releases the sequence id for the possible waiters.
//func (c *Client) releaseSeqId(seqId uint16) {
//	select {
//	case c.seqIdCh <- seqId:
//	default:
//		// If there are no waiters waiting for any sequence id to be released,
//		// reset the value corresponding to the sequence id.
//		c.deleteRespCh(seqId)
//	}
//}
//
//// getRespCh returns the response channel for the given sequence id (i.e., (respCh, true)).
//// If the response channel does not exist it returns (nil,false).
//func (c *Client) getRespCh(seqId uint16) (chan *baseResp, bool) {
//	val := atomic.LoadPointer(c.respChs[seqId])
//	if val == nil {
//		return nil, false
//	}
//	respCh := *(*chan *baseResp)(val)
//	return respCh, true
//}
//
//// saveRespCh stores the response channel for the sequence id.
//func (c *Client) saveRespCh(seqId uint16, respCh chan *baseResp) {
//	atomic.StorePointer(c.respChs[seqId], unsafe.Pointer(&respCh))
//}
//
//// saveRespCh swaps the response channel for the sequence id.
//// It returns true if the old response channel is nil, otherwise false.
//func (c *Client) swapRespCh(seqId uint16, respCh chan *baseResp) bool {
//	return atomic.CompareAndSwapPointer(c.respChs[seqId], nil, unsafe.Pointer(&respCh))
//}
//
//// deleteRespCh deletes the response channel for the sequence id.
//func (c *Client) deleteRespCh(seqId uint16) {
//	atomic.StorePointer(c.respChs[seqId], nil)
//}
