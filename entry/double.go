package entry

import (
	"encoding/binary"
	"io"

	"github.com/techplexengineer/frc-networktables-go/util"
)

// Double Entry
type Double struct {
	Base
	trueValue    float64
	isPersistent bool
}

// DoubleFromReader builds a double entry using the provided parameters
func DoubleFromReader(name string, id [2]byte, sequence [2]byte, persist byte, reader io.Reader) (*Double, error) {
	var value [8]byte
	_, err := io.ReadFull(reader, value[:])
	if err != nil {
		return nil, err
	}
	return DoubleFromItems(name, id, sequence, persist, value[:]), nil
}

// DoubleFromItems builds a double entry using the provided parameters
func DoubleFromItems(name string, id [2]byte, sequence [2]byte, persist byte, value []byte) *Double {
	val := util.BytesToFloat64(value[:8])
	persistent := (persist == flagPersist)
	return &Double{
		trueValue:    val,
		isPersistent: persistent,
		Base: Base{
			eName:  name,
			eType:  TypeDouble,
			eID:    id,
			eSeq:   sequence,
			eFlag:  persist,
			eValue: value,
		},
	}
}

// GetValue returns the value of the Double
func (o *Double) GetValue() interface{} {
	return o.trueValue
}

// IsPersistent returns whether or not the entry should persist beyond restarts.
func (o *Double) IsPersistent() bool {
	return o.isPersistent
}

// Clone returns an identical entry
func (o *Double) Clone() *Double {
	return &Double{
		trueValue:    o.trueValue,
		isPersistent: o.isPersistent,
		Base:         o.Base.clone(),
	}
}

// CompressToBytes returns a byte slice representing the Double entry
func (o *Double) CompressToBytes() []byte {
	return o.Base.compressToBytes()
}
func (o Double) GetName() string {
	return o.Base.eName
}
func (o Double) GetID() uint16 {
	return binary.LittleEndian.Uint16(o.eID[:])
}

func (Double) GetType() EntryType {
	return TypeDouble
}

func (o *Double) SetValue(newValue interface{}) {

}
