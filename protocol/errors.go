package protocol

import "errors"

//
var (

	// ErrInvalidPacketKind signals that the remote peer send an invalid packet kind.
	ErrInvalidPktKind = errors.New("invalid packet kind")

	// ErrDecodeBadPacket signals that the remote peer send a incorrect packet.
	ErrDecodeBadPacket   = errors.New("decode bad packet")

	// ErrDecodeBadPacket signals that the remote peer send a packet that is too large.
	ErrPacketTooLarge    = errors.New("packet is too large")

)