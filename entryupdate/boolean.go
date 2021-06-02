package entryupdate

import (
	"encoding/binary"
	"github.com/TechplexEngineer/FRC-NetworkTables-Go/entry"
	"io"
)

// Boolean Entry
type Boolean struct {
	Base
	trueValue bool
}

// BooleanFromReader builds a boolean entry using the provided parameters
func BooleanFromReader(id [2]byte, sequence [2]byte, etype byte, reader io.Reader) (*Boolean, error) {
	var value [1]byte
	_, err := io.ReadFull(reader, value[:])
	if err != nil {
		return nil, err
	}
	return BooleanFromItems(id, sequence, etype, value[:]), nil
}

// BooleanFromItems builds a boolean entry using the provided parameters
func BooleanFromItems(id [2]byte, sequence [2]byte, etype byte, value []byte) *Boolean {
	val := (value[0] == boolTrue)
	return &Boolean{
		trueValue: val,
		Base: Base{
			ID:    id,
			Seq:   sequence,
			Type:  entry.TypeBoolean,
			Value: value,
		},
	}
}

// GetValue returns the value of the Boolean
func (boolean *Boolean) GetValueUnsafe() interface{} {
	return boolean.trueValue
}

func (boolean *Boolean) GetValue() bool {
	return boolean.trueValue
}

// Clone returns an identical entry
func (boolean *Boolean) Clone() *Boolean {
	return &Boolean{
		trueValue: boolean.trueValue,
		Base:      boolean.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the Boolean entry
func (boolean *Boolean) CompressToBytes() []byte {
	return boolean.Base.compressToBytes()
}

func (Boolean) GetType() entry.EntryType {
	return entry.TypeBoolean
}

func (o Boolean) GetID() uint16 {
	return binary.LittleEndian.Uint16(o.ID[:])
}
