package entry

import (
	"encoding/binary"
	"io"
)

// BooleanArr Entry
type BooleanArr struct {
	Base
	trueValue    []bool
	isPersistent bool
}

// BooleanArrFromReader builds a BooleanArr entry using the provided parameters
func BooleanArrFromReader(name string, id [2]byte, sequence [2]byte, persist byte, reader io.Reader) (*BooleanArr, error) {
	var tempValSize [1]byte
	_, sizeErr := io.ReadFull(reader, tempValSize[:])
	if sizeErr != nil {
		return nil, sizeErr
	}
	valSize := int(tempValSize[0])
	value := make([]byte, valSize)
	_, valErr := io.ReadFull(reader, value[:])
	if valErr != nil {
		return nil, valErr
	}
	return BooleanArrFromItems(name, id, sequence, persist, value), nil
}

// BooleanArrFromItems builds a BooleanArr entry using the provided parameters
func BooleanArrFromItems(name string, id [2]byte, sequence [2]byte, persist byte, value []byte) *BooleanArr {
	valSize := int(value[0])
	var val []bool
	for counter := 1; counter-1 < valSize; counter++ {
		tempVal := (value[counter] == boolTrue)
		val = append(val, tempVal)
	}
	persistant := (persist == flagPersist)
	return &BooleanArr{
		trueValue:    val,
		isPersistent: persistant,
		Base: Base{
			eName:  name,
			eType:  TypeBooleanArr,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}
}

// GetValue returns the trueValue
func (o *BooleanArr) GetValue() interface{} {
	return o.trueValue
}

// GetValueAtIndex returns the value at the specified index
func (o *BooleanArr) GetValueAtIndex(index int) bool {
	return o.trueValue[index]
}

// IsPersistent returns whether or not the entry should persist beyond restarts.
func (o *BooleanArr) IsPersistent() bool {
	return o.isPersistent
}

// Clone returns an identical entry
func (o *BooleanArr) Clone() *BooleanArr {
	return &BooleanArr{
		trueValue:    o.trueValue,
		isPersistent: o.isPersistent,
		Base:         o.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the BooleanArr entry
func (o *BooleanArr) CompressToBytes() []byte {
	return o.Base.compressToBytes()
}
func (o BooleanArr) GetName() string {
	return o.Base.eName
}
func (o BooleanArr) GetID() uint16 {
	return binary.LittleEndian.Uint16(o.eID[:])
}
func (BooleanArr) GetType() EntryType {
	return TypeBooleanArr
}
func (o *BooleanArr) SetValue(newValue interface{}) {

}
