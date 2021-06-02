package entryupdate

import (
	"bytes"
	"encoding/binary"
	"github.com/techplexengineer/frc-networktables-go/entry"
	"io"

	"github.com/techplexengineer/frc-networktables-go/util"
)

// String Entry
type String struct {
	Base
	trueValue string
}

// StringFromReader builds a string entry using the provided parameters
func StringFromReader(id [2]byte, sequence [2]byte, etype byte, reader io.Reader) (*String, error) {
	valLen, sizeData := util.PeekULeb128(reader)
	valData := make([]byte, valLen)
	_, err := io.ReadFull(reader, valData[:])
	if err != nil {
		return nil, err
	}
	val := string(valData[:])
	value := append(sizeData, valData[:]...)
	return &String{
		trueValue: val,
		Base: Base{
			ID:    id,
			Seq:   sequence,
			Type:  entry.TypeString,
			Value: value,
		},
	}, nil
}

// StringFromItems builds a string entry using the provided parameters
func StringFromItems(id [2]byte, sequence [2]byte, etype byte, value []byte) *String {
	valLen, sizeLen := util.ReadULeb128(bytes.NewReader(value))
	val := string(value[sizeLen : valLen-1])
	return &String{
		trueValue: val,
		Base: Base{
			ID:    id,
			Seq:   sequence,
			Type:  entry.TypeString,
			Value: value,
		},
	}
}

// GetValue returns the value of the String
func (stringEntry *String) GetValue() string {
	return stringEntry.trueValue
}

func (stringEntry *String) GetValueUnsafe() interface{} {
	return stringEntry.trueValue
}

// Clone returns an identical entry
func (stringEntry *String) Clone() *String {
	return &String{
		trueValue: stringEntry.trueValue,
		Base:      stringEntry.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the String entry
func (stringEntry *String) CompressToBytes() []byte {
	return stringEntry.Base.compressToBytes()
}

func (String) GetType() entry.EntryType {
	return entry.TypeString
}
func (o String) GetID() uint16 {
	return binary.LittleEndian.Uint16(o.ID[:])
}
