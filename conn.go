package gosocket

import (
	"github.com/happyxcj/gosocket/protocol"
)

// Conn represents a custom socket connection,
// it describes how to send or receive a Packet.
type Conn interface {

	// Send writes the packet to the connection.
	Send(p protocol.Packet) error

	// Receive reads a Packet from the connection.
	Receive() (protocol.Packet, error)

	// Close closes the connection.
	// Any blocked Read or Write operations will be unblocked and return errors.
	Close() error
}

