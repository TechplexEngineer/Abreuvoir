package entry

import (
	"encoding/binary"
	"io"
	//"github.com/HowardStark/abreuvoir/util"
)

// DoubleArr Entry
type DoubleArr struct {
	Base
	trueValue    []float64
	isPersistent bool
}

// DoubleArrFromReader builds a DoubleArr entry using the provided parameters
func DoubleArrFromReader(name string, id [2]byte, sequence [2]byte, persist byte, reader io.Reader) (*DoubleArr, error) {
	var tempValSize [1]byte
	_, sizeErr := io.ReadFull(reader, tempValSize[:])
	if sizeErr != nil {
		return nil, sizeErr
	}
	valSize := int(tempValSize[0])
	value := make([]byte, valSize*8)
	_, valErr := io.ReadFull(reader, value[:])
	if valErr != nil {
		return nil, valErr
	}
	return DoubleArrFromItems(name, id, sequence, persist, value), nil
}

// DoubleArrFromItems builds a DoubleArr entry using the provided parameters
func DoubleArrFromItems(name string, id [2]byte, sequence [2]byte, persist byte, value []byte) *DoubleArr {
	//valSize := int(value[0])
	var val []float64
	//for counter := 0; (counter)/8 < valSize; counter += 8 {
	//	tempVal := util.BytesToFloat64(value[counter : counter+8])
	//	val = append(val, tempVal)
	//}
	persistant := (persist == flagPersist)
	return &DoubleArr{
		trueValue:    val,
		isPersistent: persistant,
		Base: Base{
			eName:  name,
			eType:  TypeDoubleArr,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}
}

// GetValue returns the value of the DoubleArr
func (o *DoubleArr) GetValue() interface{} {
	return o.trueValue
}

// GetValueAtIndex returns the value at the specified index
func (o *DoubleArr) GetValueAtIndex(index int) float64 {
	return o.trueValue[index]
}

// IsPersistant returns whether or not the entry should persist beyond restarts.
func (o *DoubleArr) IsPersistant() bool {
	return o.isPersistent
}

// Clone returns an identical entry
func (o *DoubleArr) Clone() *DoubleArr {
	return &DoubleArr{
		trueValue:    o.trueValue,
		isPersistent: o.isPersistent,
		Base:         o.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the DoubleArr entry
func (o *DoubleArr) CompressToBytes() []byte {
	return o.Base.compressToBytes()
}

func (o DoubleArr) GetName() string {
	return o.Base.eName
}
func (o DoubleArr) GetID() uint16 {
	return binary.LittleEndian.Uint16(o.eID[:])
}
func (DoubleArr) GetType() EntryType {
	return TypeDoubleArr
}
func (o *DoubleArr) SetValue(newValue interface{}) {

}
