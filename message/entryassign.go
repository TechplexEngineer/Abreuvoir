package message

import (
	"io"

	"github.com/HowardStark/abreuvoir/entry"
)

// EntryAssign message
type EntryAssign struct {
	Base
	entry entry.IEntry
}

// EntryAssignFromReader builds an EntryAssign message and its entry from a reader
func EntryAssignFromReader(reader io.Reader) (*EntryAssign, error) {
	tempEntry, err := entry.BuildFromReader(reader)
	if err != nil {
		return nil, err
	}
	return &EntryAssign{
		entry: tempEntry,
		Base: Base{
			mType: TypeEntryAssign,
			mData: tempEntry.CompressToBytes(),
		},
	}, nil
}

// EntryAssignFromEntry builds an EntryAssign message from an entry
func EntryAssignFromEntry(newEntry entry.IEntry) *EntryAssign {
	return &EntryAssign{
		entry: newEntry,
		Base: Base{
			mType: TypeEntryAssign,
			mData: newEntry.CompressToBytes(),
		},
	}
}

// EntryAssignFromBytes builds an EntryAssign message and its entry from a byte slice
func EntryAssignFromBytes(data []byte) (*EntryAssign, error) {
	tempEntry, err := entry.BuildFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &EntryAssign{
		entry: tempEntry,
		Base: Base{
			mType: TypeEntryAssign,
			mData: data,
		},
	}, nil
}

// GetEntry returns the entry associated with this EntryAssign message
func (entryAssign *EntryAssign) GetEntry() entry.IEntry {
	return entryAssign.entry
}

// CompressToBytes returns the message in its byte array form
func (entryAssign *EntryAssign) CompressToBytes() []byte {
	return entryAssign.Base.compressToBytes()
}

// GetType returns the message's type
func (entryAssign *EntryAssign) GetType() MessageType {
	return TypeEntryAssign
}
