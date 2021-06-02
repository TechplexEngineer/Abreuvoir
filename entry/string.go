package entry

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/TechplexEngineer/FRC-NetworkTables-Go/util"
)

// String Entry
type String struct {
	Base
	trueValue    string
	isPersistent bool
}

// StringFromReader builds a string entry using the provided parameters
func StringFromReader(name string, id [2]byte, sequence [2]byte, persist byte, reader io.Reader) (*String, error) {
	valLen, sizeData := util.PeekULeb128(reader)
	valData := make([]byte, valLen)
	_, err := io.ReadFull(reader, valData[:])
	if err != nil {
		return nil, err
	}
	val := string(valData[:])
	persistent := (persist == flagPersist)
	value := append(sizeData, valData[:]...)
	return &String{
		trueValue:    val,
		isPersistent: persistent,
		Base: Base{
			eName:  name,
			eType:  TypeString,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}, nil
}

// StringFromItems builds a string entry using the provided parameters
func StringFromItems(name string, id [2]byte, sequence [2]byte, persist byte, value []byte) *String {
	valLen, sizeLen := util.ReadULeb128(bytes.NewReader(value))
	val := string(value[sizeLen : valLen-1])
	persistent := (persist == flagPersist)
	return &String{
		trueValue:    val,
		isPersistent: persistent,
		Base: Base{
			eName:  name,
			eType:  TypeString,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}
}

// GetValue returns the value of the String
func (o *String) GetValue() interface{} {
	return o.trueValue
}

// IsPersistent returns whether or not the entry should persist beyond restarts.
func (o *String) IsPersistent() bool {
	return o.isPersistent
}

// Clone returns an identical entry
func (o *String) Clone() *String {
	return &String{
		trueValue:    o.trueValue,
		isPersistent: o.isPersistent,
		Base:         o.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the String entry
func (o *String) CompressToBytes() []byte {
	return o.Base.compressToBytes()
}

func (o String) GetName() string {
	return o.Base.eName
}
func (o String) GetID() uint16 {
	return binary.LittleEndian.Uint16(o.eID[:])
}
func (String) GetType() EntryType {
	return TypeString
}
func (o *String) SetValue(newValue interface{}) {

}
