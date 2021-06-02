package message

import (
	"io"

	"github.com/HowardStark/abreuvoir/entryupdate"
)

// EntryUpdate message
type EntryUpdate struct {
	Base
	Update entryupdate.IEntryUpdate
}

// EntryUpdateFromReader meme
func EntryUpdateFromReader(reader io.Reader) (*EntryUpdate, error) {
	tempUpdate, err := entryupdate.BuildFromReader(reader)
	if err != nil {
		return nil, err
	}
	return &EntryUpdate{
		Update: tempUpdate,
		Base: Base{
			mType: TypeEntryUpdate,
			mData: tempUpdate.CompressToBytes(),
		},
	}, nil
}

// EntryUpdateFromUpdate builds an EntryUpdate message from an update
func EntryUpdateFromUpdate(entryUpdate entryupdate.IEntryUpdate) *EntryUpdate {
	return &EntryUpdate{
		Update: entryUpdate,
		Base: Base{
			mType: TypeEntryUpdate,
			mData: entryUpdate.CompressToBytes(),
		},
	}
}

// GetUpdate returns the Update associated with this EntryUpdate
func (entryUpdate *EntryUpdate) GetUpdate() entryupdate.IEntryUpdate {
	return entryUpdate.Update
}

// CompressToBytes returns the message in its byte array form
func (entryUpdate *EntryUpdate) CompressToBytes() []byte {
	return entryUpdate.Base.compressToBytes()
}

// GetType returns the message's type
func (entryUpdate *EntryUpdate) GetType() MessageType {
	return TypeEntryUpdate
}
