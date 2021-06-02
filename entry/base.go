package entry

import (
	"bytes"
	"errors"
	"io"

	"github.com/TechplexEngineer/FRC-NetworkTables-Go/util"
)

type EntryType byte

func (t EntryType) String() string {
	switch t {
	case TypeBoolean:
		return "Boolean"
	case TypeDouble:
		return "Double"
	case TypeString:
		return "String"
	case TypeRaw:
		return "Raw"
	case TypeBooleanArr:
		return "BooleanArr"
	case TypeDoubleArr:
		return "DoubleArr"
	case TypeStringArr:
		return "StringArr"
	case TypeRPCDef:
		return "RPCDef"
	default:
		return "UNKNOWN ENTRY TYPE"
	}
}

func (t EntryType) Byte() byte {
	return byte(t)
}

const (
	TypeBoolean    EntryType = 0x00
	TypeDouble     EntryType = 0x01
	TypeString     EntryType = 0x02
	TypeRaw        EntryType = 0x03
	TypeBooleanArr EntryType = 0x10
	TypeDoubleArr  EntryType = 0x11
	TypeStringArr  EntryType = 0x12
	TypeRPCDef     EntryType = 0x20

	flagTemporary byte = 0x00
	flagPersist   byte = 0x01
	flagReserved  byte = 0xFE

	boolFalse byte = 0x00
	boolTrue  byte = 0x01
)

var (
	// idSent is the required ID for an entry that is being created/sent from the client
	idSent = [2]byte{0xFF, 0xFF}
)

// Base is the base struct for entries.
type Base struct {
	eName  string
	eType  EntryType
	eID    [2]byte
	eSeq   [2]byte
	eFlag  byte
	eValue []byte
}

// BuildFromReader creates an entry using the reader passed in
func BuildFromReader(reader io.Reader) (IEntry, error) {
	nameLen, _ := util.ReadULeb128(reader)
	nameData := make([]byte, nameLen)
	_, nameErr := io.ReadFull(reader, nameData[:])
	if nameErr != nil {
		return nil, nameErr
	}
	name := string(nameData[:])
	var typeData [1]byte
	_, typeErr := io.ReadFull(reader, typeData[:])
	if typeErr != nil {
		return nil, typeErr
	}
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
	var flagData [1]byte
	_, flagErr := io.ReadFull(reader, flagData[:])
	if flagErr != nil {
		return nil, flagErr
	}
	_, _ = nameData, name
	switch EntryType(typeData[0]) {
	case TypeBoolean:
		return BooleanFromReader(name, idData, seqData, flagData[0], reader)
	case TypeDouble:
		return DoubleFromReader(name, idData, seqData, flagData[0], reader)
	case TypeString:
		return StringFromReader(name, idData, seqData, flagData[0], reader)
	case TypeRaw:
		return RawFromReader(name, idData, seqData, flagData[0], reader)
	case TypeBooleanArr:
		return BooleanArrFromReader(name, idData, seqData, flagData[0], reader)
	case TypeDoubleArr:
		return DoubleArrFromReader(name, idData, seqData, flagData[0], reader)
	case TypeStringArr:
		return StringArrFromReader(name, idData, seqData, flagData[0], reader)
	default:
		return nil, errors.New("entry: Unknown entry type")
	}
}

// BuildFromBytes creates an entry using the data passed in.
func BuildFromBytes(data []byte) (IEntry, error) {
	nameLen, sizeLen := util.ReadULeb128(bytes.NewReader(data))
	dName := string(data[sizeLen : nameLen-1])
	dType := EntryType(data[nameLen])
	dID := [2]byte{data[nameLen+1], data[nameLen+2]}
	dSeq := [2]byte{data[nameLen+3], data[nameLen+4]}
	dFlag := data[nameLen+5]
	dValue := data[nameLen+6:]
	switch dType {
	case TypeBoolean:
		return BooleanFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case TypeDouble:
		return DoubleFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case TypeString:
		return StringFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case TypeRaw:
		return RawFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case TypeBooleanArr:
		return BooleanArrFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case TypeDoubleArr:
		return DoubleArrFromItems(dName, dID, dSeq, dFlag, dValue), nil
	case TypeStringArr:
		return StringArrFromItems(dName, dID, dSeq, dFlag, dValue), nil
	default:
		return nil, errors.New("entry: Unknown entry type")
	}
}

func (base *Base) clone() Base {
	return *base
}

// CompressToBytes remakes the original byte slice to represent this entry
func (base *Base) compressToBytes() []byte {
	var output []byte
	nameBytes := []byte(base.eName)
	nameLen := util.EncodeULeb128(uint32(len(nameBytes)))
	output = append(output, nameLen...)
	output = append(output, nameBytes...)
	output = append(output, base.eType.Byte())
	output = append(output, base.eID[:]...)
	output = append(output, base.eSeq[:]...)
	output = append(output, base.eFlag)
	output = append(output, base.eValue...)
	return output
}
