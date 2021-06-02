package entryupdate

import (
	"errors"
	"github.com/techplexengineer/frc-networktables-go/entry"
	"io"
)

const (
	//typeBoolean    byte = 0x00
	//typeDouble     byte = 0x01
	//typeString     byte = 0x02
	//typeRaw        byte = 0x03
	//typeBooleanArr byte = 0x10
	//typeDoubleArr  byte = 0x11
	//typeStringArr  byte = 0x12
	//typeRPCDef     byte = 0x20

	flagTemporary byte = 0x00
	flagPersist   byte = 0x01
	flagReserved  byte = 0xFE

	boolFalse byte = 0x00
	boolTrue  byte = 0x01
)

// Base is the base struct for Entry updates
type Base struct {
	ID    [2]byte
	Seq   [2]byte
	Type  entry.EntryType
	Value []byte
}

// BuildFromReader sleep
func BuildFromReader(reader io.Reader) (IEntryUpdate, error) {
	var idData [2]byte
	_, idErr := io.ReadFull(reader, idData[:])
	if idErr != nil {
		return nil, idErr
	}
	var seqData [2]byte
	_, seqErr := io.ReadFull(reader, seqData[:])
	if seqErr != nil {
		return nil, seqErr
	}
	var typeData [1]byte
	_, typeErr := io.ReadFull(reader, typeData[:])
	if typeErr != nil {
		return nil, typeErr
	}
	switch entry.EntryType(typeData[0]) {
	case entry.TypeBoolean:
		return BooleanFromReader(idData, seqData, typeData[0], reader)
	case entry.TypeDouble:
		return DoubleFromReader(idData, seqData, typeData[0], reader)
	case entry.TypeString:
		return StringFromReader(idData, seqData, typeData[0], reader)
	case entry.TypeRaw:
		return RawFromReader(idData, seqData, typeData[0], reader)
	case entry.TypeBooleanArr:
		return BooleanArrFromReader(idData, seqData, typeData[0], reader)
	case entry.TypeDoubleArr:
		return DoubleArrFromReader(idData, seqData, typeData[0], reader)
	case entry.TypeStringArr:
		return StringArrFromReader(idData, seqData, typeData[0], reader)
	default:
		return nil, errors.New("entry: Unknown entry type")
	}
}

func (base *Base) clone() Base {
	return *base
}

// CompressToBytes returns a byte slice representing the Update entry
func (base *Base) compressToBytes() []byte {
	var compressed []byte
	compressed = append(compressed, base.ID[:]...)
	compressed = append(compressed, base.Seq[:]...)
	compressed = append(compressed, base.Type.Byte())
	compressed = append(compressed, base.Value...)
	return compressed
}
