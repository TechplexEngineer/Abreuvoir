package entryupdate

import "github.com/TechplexEngineer/FRC-NetworkTables-Go/entry"

// IEntryUpdate is the entry update interface
type IEntryUpdate interface {
	CompressToBytes() []byte
	GetType() entry.EntryType //@todo change to entryType
	GetID() uint16
	GetValueUnsafe() interface{}
}
