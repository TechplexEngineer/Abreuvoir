package entry

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/HowardStark/abreuvoir/util"
)

// StringArr Entry
type StringArr struct {
	Base
	trueValue    []string
	isPersistant bool
}

// StringArrFromReader builds a StringArr entry using the provided parameters
func StringArrFromReader(name string, id [2]byte, sequence [2]byte, persist byte, reader io.Reader) (*StringArr, error) {
	var value []byte
	var tempValSize [1]byte
	_, sizeErr := io.ReadFull(reader, tempValSize[:])
	if sizeErr != nil {
		return nil, sizeErr
	}
	value = append(value, tempValSize[0])
	valSize := int(tempValSize[0])
	var val []string
	for counter := 0; counter < valSize; counter++ {
		strLen, sizeData := util.PeekULeb128(reader)
		value = append(value, sizeData...)
		strData := make([]byte, strLen)
		_, strErr := io.ReadFull(reader, strData[:])
		if strErr != nil {
			return nil, strErr
		}
		value = append(value, strData[:]...)
		val = append(val, string(strData[:]))
	}
	persistant := (persist == flagPersist)
	return &StringArr{
		trueValue:    val,
		isPersistant: persistant,
		Base: Base{
			eName:  name,
			eType:  TypeStringArr,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}, nil
}

// StringArrFromItems builds a StringArr entry using the provided parameters
func StringArrFromItems(name string, id [2]byte, sequence [2]byte, persist byte, value []byte) *StringArr {
	valSize := int(value[0])
	var val []string
	var previousPos uint32 = 1
	for counter := 0; counter < valSize; counter++ {
		strPos, sizePos := util.ReadULeb128(bytes.NewReader(value[previousPos:]))
		strPos += previousPos
		sizePos += previousPos
		tempVal := string(value[sizePos : strPos-1])
		val = append(val, tempVal)
		previousPos = strPos - 1
	}
	persistant := (persist == flagPersist)
	return &StringArr{
		trueValue:    val,
		isPersistant: persistant,
		Base: Base{
			eName:  name,
			eType:  TypeStringArr,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}
}

// GetValue returns the value of the StringArr
func (o *StringArr) GetValue() interface{} {
	return o.trueValue
}

// GetValueAtIndex returns the value at the specified index
func (o *StringArr) GetValueAtIndex(index int) string {
	return o.trueValue[index]
}

// IsPersistant returns whether or not the entry should persist beyond restarts.
func (o *StringArr) IsPersistant() bool {
	return o.isPersistant
}

// Clone returns an identical entry
func (o *StringArr) Clone() *StringArr {
	return &StringArr{
		trueValue:    o.trueValue,
		isPersistant: o.isPersistant,
		Base:         o.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the StringArr entry
func (o *StringArr) CompressToBytes() []byte {
	return o.Base.compressToBytes()
}

func (o StringArr) GetName() string {
	return o.Base.eName
}
func (o StringArr) GetID() uint16 {
	return binary.LittleEndian.Uint16(o.eID[:])
}
func (StringArr) GetType() EntryType {
	return TypeStringArr
}

func (o *StringArr) SetValue(newValue interface{}) {

}
