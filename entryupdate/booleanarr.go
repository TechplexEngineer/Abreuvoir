package entryupdate

import (
	"encoding/binary"
	"github.com/TechplexEngineer/FRC-NetworkTables-Go/entry"
	"io"
)

// BooleanArr Entry
type BooleanArr struct {
	Base
	trueValue []bool
}

// BooleanArrFromReader builds a BooleanArr entry using the provided parameters
func BooleanArrFromReader(id [2]byte, sequence [2]byte, etype byte, reader io.Reader) (*BooleanArr, error) {
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
	return BooleanArrFromItems(id, sequence, etype, value), nil
}

// BooleanArrFromItems builds a BooleanArr entry using the provided parameters
func BooleanArrFromItems(id [2]byte, sequence [2]byte, etype byte, value []byte) *BooleanArr {
	valSize := int(value[0])
	var val []bool
	for counter := 1; counter-1 < valSize; counter++ {
		tempVal := (value[counter] == boolTrue)
		val = append(val, tempVal)
	}
	return &BooleanArr{
		trueValue: val,
		Base: Base{
			ID:    id,
			Seq:   sequence,
			Type:  entry.TypeBooleanArr,
			Value: value,
		},
	}
}

// GetValue returns the trueValue
func (booleanArr *BooleanArr) GetValue() []bool {
	return booleanArr.trueValue
}

func (booleanArr *BooleanArr) GetValueUnsafe() interface{} {
	return booleanArr.trueValue
}

// GetValueAtIndex returns the value at the specified index
func (booleanArr *BooleanArr) GetValueAtIndex(index int) bool {
	return booleanArr.trueValue[index]
}

// Clone returns an identical entry
func (booleanArr *BooleanArr) Clone() *BooleanArr {
	return &BooleanArr{
		trueValue: booleanArr.trueValue,
		Base:      booleanArr.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the BooleanArr entry
func (booleanArr *BooleanArr) CompressToBytes() []byte {
	return booleanArr.Base.compressToBytes()
}

func (BooleanArr) GetType() entry.EntryType {
	return entry.TypeDoubleArr
}
func (o BooleanArr) GetID() uint16 {
	return binary.LittleEndian.Uint16(o.ID[:])
}
