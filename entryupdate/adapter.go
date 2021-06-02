package entryupdate

import "github.com/HowardStark/abreuvoir/entry"

// IEntryUpdate is the entry update interface
type IEntryUpdate interface {
	CompressToBytes() []byte
	GetType() entry.EntryType //@todo change to entryType
	GetID() uint16
	GetValueUnsafe() interface{}
}
