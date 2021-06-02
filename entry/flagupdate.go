package entry

import (
	"encoding/binary"
	"io"
)

// FlagUpdate entry is a partial entry containing only certain fields of an actual entry
type FlagUpdate struct {
	ID           [2]byte
	IsPersistant bool
	flags        byte
}

// FlagUpdateFromReader head hurts
func FlagUpdateFromReader(reader io.Reader) (*FlagUpdate, error) {
	var dID [2]byte
	_, idErr := io.ReadFull(reader, dID[:])
	if idErr != nil {
		return nil, idErr
	}
	var dFlags [1]byte
	_, flagErr := io.ReadFull(reader, dFlags[:])
	if flagErr != nil {
		return nil, flagErr
	}
	dPersist := (dFlags[0] == flagPersist)
	return &FlagUpdate{
		ID:           dID,
		IsPersistant: dPersist,
		flags:        dFlags[0],
	}, nil
}

// FlagUpdateFromBytes builds an FlagUpdate from a byte slice
func FlagUpdateFromBytes(data []byte) *FlagUpdate {
	dID := [2]byte{data[0], data[1]}
	dFlags := data[2]
	dPersist := (dFlags == flagPersist)
	return &FlagUpdate{
		ID:           dID,
		IsPersistant: dPersist,
		flags:        dFlags,
	}
}

// FlagUpdateFromItems builds an FlagUpdate using the provided parameters
func FlagUpdateFromItems(dID [2]byte, dFlags byte) *FlagUpdate {
	dPersist := (dFlags == flagPersist)
	return &FlagUpdate{
		ID:           dID,
		IsPersistant: dPersist,
		flags:        dFlags,
	}
}

// CompressToBytes returns a byte slice representing the FlagUpdate entry
func (o *FlagUpdate) CompressToBytes() []byte {
	compressed := []byte{}
	compressed = append(compressed, o.ID[:]...)
	compressed = append(compressed, o.flags)
	return compressed
}

//func (o FlagUpdate) GetName() string {
//	return o.Base.eName
//}
func (o FlagUpdate) GetID() uint16 {
	return binary.LittleEndian.Uint16(o.ID[:])
}
