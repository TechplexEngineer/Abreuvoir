package entry

// IEntry is the interface for Entries
type IEntry interface {
	GetName() string
	GetValue() interface{}
	SetValue(interface{})
	CompressToBytes() []byte
	GetID() uint16
	GetType() EntryType
}
