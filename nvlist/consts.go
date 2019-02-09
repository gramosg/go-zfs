package nvlist

type nvtype uint32

const (
	typeUnknown nvtype = iota
	typeBoolean
	typeByte
	typeInt16
	typeUint16
	typeInt32
	typeUint32
	typeInt64
	typeUint64
	typeString
	typeByteArray
	typeInt16Array
	typeUint16Array
	typeInt32Array
	typeUint32Array
	typeInt64Array
	typeUint64Array
	typeStringArray
	typeHrtime
	typeNvlist
	typeNvlistArray
	typeBooleanValue
	typeInt8
	typeUint8
	typeBooleanArray
	typeInt8Array
	typeUint8Array
	typeDouble
)

func (nvt nvtype) decode(s nvPairReader, arraySize uint32) (interface{}, error) {
	switch nvt {
	case typeUnknown:
		return nil, ErrInvalidData
	case typeBoolean:
		return true, nil
	case typeByte:
		return s.ReadByte()
	case typeInt16:
		var val int16
		err := s.readInt(&val)
		return val, err
	case typeUint16:
		var val uint16
		err := s.readInt(&val)
		return val, err
	case typeInt32:
		var val int32
		err := s.readInt(&val)
		return val, err
	case typeUint32:
		var val uint32
		err := s.readInt(&val)
		return val, err
	case typeInt64:
		var val int64
		err := s.readInt(&val)
		return val, err
	case typeUint64:
		var val uint64
		err := s.readInt(&val)
		return val, err
	case typeString:
		data, err := s.ReadBytes(0x00)
		if err != nil {
			return nil, err
		}
		return string(data[:len(data)-1]), nil
	case typeByteArray:
		return s.readN(int(arraySize))
	case typeInt16Array:
		val := make([]int16, arraySize)
		for i := uint32(0); i < arraySize; i++ {
			if err := s.readInt(&val[i]); err != nil {
				return nil, err
			}
		}
		return val, nil
	case typeUint16Array:
		val := make([]uint16, arraySize)
		for i := uint32(0); i < arraySize; i++ {
			if err := s.readInt(&val[i]); err != nil {
				return nil, err
			}
		}
		return val, nil
	case typeInt32Array:
		val := make([]int32, arraySize)
		for i := uint32(0); i < arraySize; i++ {
			if err := s.readInt(&val[i]); err != nil {
				return nil, err
			}
		}
		return val, nil
	case typeUint32Array:
		val := make([]uint32, arraySize)
		for i := uint32(0); i < arraySize; i++ {
			if err := s.readInt(&val[i]); err != nil {
				return nil, err
			}
		}
		return val, nil
	case typeInt64Array:
		val := make([]int64, arraySize)
		for i := uint32(0); i < arraySize; i++ {
			if err := s.readInt(&val[i]); err != nil {
				return nil, err
			}
		}
		return val, nil
	case typeUint64Array:
		val := make([]uint64, arraySize)
		for i := uint32(0); i < arraySize; i++ {
			if err := s.readInt(&val[i]); err != nil {
				return nil, err
			}
		}
		return val, nil
	case typeStringArray:
		val := make([]string, arraySize)
		s.skipN(int(8 * arraySize)) // Skip pointers
		// Pointers are always aligned
		for i := uint32(0); i < arraySize; i++ {
			data, err := s.ReadBytes(0x00)
			if err != nil {
				return nil, err
			}
			val[i] = string(data[:len(data)-1])
		}
		return val, nil
	case typeHrtime:
		return nil, ErrUnsupportedType
	case typeNvlist:
		val := make(map[string]interface{})
		err := s.nvlist.readPairs(val)
		return val, err
	case typeNvlistArray:
		val := make([]map[string]interface{}, arraySize)
		// Drop unused data (nvlist header @ 8 bytes + 64 bit pointer @ 8 bytes)
		s.skipN(int((8 + 8) * arraySize))
		for i := uint32(0); i < arraySize; i++ {
			val[i] = make(map[string]interface{})
			err := s.nvlist.readPairs(val[i])
			if err != nil {
				return nil, err
			}
		}
		return val, nil
	case typeBooleanValue:
		var tmp int32
		err := s.readInt(&tmp)
		if err != nil {
			return nil, err
		}
		switch tmp {
		case 0:
			return false, nil
		case 1:
			return true, nil
		default:
			return nil, ErrInvalidData
		}
	case typeInt8:
		var val int8
		err := s.readInt(&val)
		return val, err
	case typeUint8:
		var val uint8
		err := s.readInt(&val)
		return val, err
	case typeBooleanArray:
		var tmp int32
		val := make([]bool, arraySize)
		for i := uint32(0); i < arraySize; i++ {
			if err := s.readInt(&tmp); err != nil {
				return nil, err
			}
			switch tmp {
			case 0:
				val[i] = false
			case 1:
				val[i] = true
			default:
				return nil, ErrInvalidData
			}
		}
		return val, nil
	case typeInt8Array:
		val := make([]int8, arraySize)
		for i := uint32(0); i < arraySize; i++ {
			if err := s.readInt(&val[i]); err != nil {
				return nil, err
			}
		}
		return val, nil
	case typeUint8Array:
		val := make([]uint8, arraySize)
		for i := uint32(0); i < arraySize; i++ {
			if err := s.readInt(&val[i]); err != nil {
				return nil, err
			}
		}
		return val, nil
	case typeDouble:
		var val float64
		err := s.readInt(&val)
		return val, err
	default:
		return nil, ErrInvalidData
	}
}
