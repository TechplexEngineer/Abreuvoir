package message

// IMessage is the Message interface
type IMessage interface {
	CompressToBytes() []byte
	GetType() MessageType
}
