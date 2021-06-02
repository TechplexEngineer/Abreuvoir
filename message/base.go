package message

import (
	"errors"
	"io"
)

type MessageType byte

func (m MessageType) String() string {
	switch m {
	case TypeKeepAlive:
		return "KeepAlive"
	case TypeClientHello:
		return "ClientHello"
	case TypeProtoUnsupported:
		return "ProtoUnsupported"
	case TypeServerHelloComplete:
		return "ServerHelloComplete"
	case TypeServerHello:
		return "ServerHello"
	case TypeClientHelloComplete:
		return "ClientHelloComplete"
	case TypeEntryAssign:
		return "EntryAssign"
	case TypeEntryUpdate:
		return "EntryUpdate"
	case TypeEntryFlagUpdate:
		return "EntryFlagUpdate"
	case TypeEntryDelete:
		return "EntryDelete"
	case TypeClearAllEntries:
		return "ClearAllEntries"
	case TypeRPCExec:
		return "RPCExec"
	case TypeRPCResponse:
		return "RPCResponse"
	default:
		return "UNKNOWN"
	}
}

const (
	// TypeKeepAlive is the message type for the Keep Alive message
	TypeKeepAlive MessageType = 0x00
	// TypeClientHello is the message type for the Client Hello message
	TypeClientHello MessageType = 0x01
	// TypeProtoUnsupported is the message type for the Protocol Unsupported message
	TypeProtoUnsupported MessageType = 0x02
	// TypeServerHelloComplete is the message type for the Server Hello Complete message
	TypeServerHelloComplete MessageType = 0x03
	// TypeServerHello is the message type for the Server Hello message
	TypeServerHello MessageType = 0x04
	// TypeClientHelloComplete is the message type for the Client Hello Complete message
	TypeClientHelloComplete MessageType = 0x05
	// TypeEntryAssign is the message type for the Entry Assign message
	TypeEntryAssign MessageType = 0x10
	// TypeEntryUpdate is the message type for the Entry Update message
	TypeEntryUpdate MessageType = 0x11
	// TypeEntryFlagUpdate is the message type for the Entry Flag Update message
	TypeEntryFlagUpdate MessageType = 0x12
	// TypeEntryDelete is the message type for the Entry Delete message
	TypeEntryDelete MessageType = 0x13
	// TypeClearAllEntries is the message type for the Clear All Entries message
	TypeClearAllEntries MessageType = 0x14
	// TypeRPCExec is the message type for the Remote Procedure Call Execute message
	TypeRPCExec MessageType = 0x20
	// TypeRPCResponse is the message type for the RPC Response message
	TypeRPCResponse MessageType = 0x21

	lsbFirstConnect byte = 0x00
	lsbReconnect    byte = 0x01
)

// Base is the base struct for Messages
type Base struct {
	mType MessageType
	mData []byte
}

// BuildFromReader identifies and builds the message with the same type as the
// message type passed in
func BuildFromReader(messageType MessageType, reader io.Reader) (IMessage, error) {
	switch messageType {
	case TypeKeepAlive:
		return KeepAliveFromReader(), nil
	case TypeClientHello:
		return ClientHelloFromReader(reader)
	case TypeProtoUnsupported:
		return ProtoUnsupportedFromReader(reader)
	case TypeServerHelloComplete:
		return ServerHelloCompleteFromReader(), nil
	case TypeServerHello:
		return ServerHelloFromReader(reader)
	case TypeClientHelloComplete:
		return ClientHelloCompleteFromReader(), nil
	case TypeEntryAssign:
		return EntryAssignFromReader(reader)
	case TypeEntryUpdate:
		return EntryUpdateFromReader(reader)
	case TypeEntryFlagUpdate:
		return EntryFlagUpdateFromReader(reader)
	case TypeEntryDelete:
		return EntryDeleteFromReader(reader)
	case TypeClearAllEntries:
		//fallthrough
	case TypeRPCExec:
		//fallthrough
	case TypeRPCResponse:
		//fallthrough
	default:
		return nil, errors.New("message: Unknown message type")
	}
	return nil, errors.New("message: Unknown message type")
}

func (m MessageType) Byte() byte {
	return byte(m)
}

// compressToBytes remakes the original byte slice to represent this entry
func (base *Base) compressToBytes() []byte {
	output := []byte{}
	output = append(output, base.mType.Byte())
	output = append(output, base.mData...)
	return output
}
