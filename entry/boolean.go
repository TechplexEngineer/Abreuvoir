package entry

import (
	"encoding/binary"
	"io"
)

// Boolean Entry
type Boolean struct {
	Base
	trueValue    bool
	isPersistent bool
}

// BooleanFromReader builds a boolean entry using the provided parameters
func BooleanFromReader(name string, id [2]byte, sequence [2]byte, persist byte, reader io.Reader) (*Boolean, error) {
	var value [1]byte
	_, err := io.ReadFull(reader, value[:])
	if err != nil {
		return nil, err
	}
	return BooleanFromItems(name, id, sequence, persist, value[:]), nil
}

// BooleanFromItems builds a boolean entry using the provided parameters
func BooleanFromItems(name string, id [2]byte, sequence [2]byte, persist byte, value []byte) *Boolean {
	val := (value[0] == boolTrue)
	persistent := (persist == flagPersist)
	return &Boolean{
		trueValue:    val,
		isPersistent: persistent,
		Base: Base{
			eName:  name,
			eType:  TypeBoolean,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}
}

// GetValue returns the value of the Boolean
func (o *Boolean) GetValue() interface{} {
	return o.trueValue
}

// IsPersistent returns whether or not the entry should persist beyond restarts.
func (o *Boolean) IsPersistent() bool {
	return o.isPersistent
}

// Clone returns an identical entry
func (o *Boolean) Clone() *Boolean {
	return &Boolean{
		trueValue:    o.trueValue,
		isPersistent: o.isPersistent,
		Base:         o.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the Boolean entry
func (o *Boolean) CompressToBytes() []byte {
	return o.Base.compressToBytes()
}

func (o Boolean) GetName() string {
	return o.Base.eName
}

func (o Boolean) GetID() uint16 {
	return binary.LittleEndian.Uint16(o.eID[:])
}

func (Boolean) GetType() EntryType {
	return TypeBoolean
}

func (o *Boolean) SetValue(newValue interface{}) {

}
