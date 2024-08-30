package common

import (
	"errors"
	"fmt"
)

// SerializeString serializes a string into a byte buffer
func SerializeString(value string) []byte {
	data := []byte{0x73}
	data = append(data, []byte(value)...)
	data = append(data, 0x00)
	return data
}

// DeserializeString deserializes a string from a byte buffer
func DeserializeString(data []byte) (*string, error) {
	if data[0] != 0x73 {
		return nil, errors.New("missing delimiter")
	}

	data = data[1:]
	value := ""

	// Read string until null byte
	for data[0] != 0x00 && len(data) > 1 {
		value += string(data[0])
		data = data[1:]
	}

	return &value, nil
}

// SerializeBinary serializes binary data into a byte buffer
func SerializeBinary(value []byte) []byte {
	data := []byte{0x62}
	size := len(value)
	sizeBytes := []byte{
		byte(size>>24) & 0xFF,
		byte(size>>16) & 0xFF,
		byte(size>>8) & 0xFF,
		byte(size) & 0xFF,
	}
	data = append(data, sizeBytes...)
	data = append(data, value...)
	return data
}

// DeserializeBinary deserializes binary data from a byte buffer
func DeserializeBinary(data []byte) ([]byte, error) {
	if data[0] != 0x62 {
		return nil, errors.New("missing delimiter")
	}

	data = data[1:]
	size := (int(data[0]) << 24) + (int(data[1]) << 16) + (int(data[2]) << 8) + int(data[3])
	data = data[4:]

	if len(data) < size {
		return nil, errors.New("missing data")
	}

	return data[:size], nil
}

// SerializeDataListInner serializes a list of data into a byte buffer with outer brackets
func SerializeDataListInner(data []interface{}) ([]byte, error) {
	serializedData := []byte{0x5B} // Opening bracket '['

	for _, item := range data {
		switch v := item.(type) {
		case string:
			serializedData = append(serializedData, SerializeString(v)...)
		case []byte:
			serializedData = append(serializedData, SerializeBinary(v)...)
		case []interface{}:
			nestedData, err := SerializeDataListInner(v)
			if err != nil {
				return nil, err
			}
			serializedData = append(serializedData, nestedData...)
		default:
			return nil, fmt.Errorf("unsupported type %T serialized in list", v)
		}
	}

	serializedData = append(serializedData, 0x5D) // Closing bracket ']'
	return serializedData, nil
}

// SerializeDataList serializes a list of data into a byte buffer without outer brackets
func SerializeDataList(data []interface{}) ([]byte, error) {
	serializedData, err := SerializeDataListInner(data)
	if err != nil {
		return nil, err
	}

	// Remove outer and inner brackets
	return serializedData[1 : len(serializedData)-1], nil
}

// DeserializeDataList deserializes a list of data from a byte buffer
func DeserializeDataList(data []byte) ([]interface{}, error) {
	var result []interface{}

	for len(data) > 0 {
		switch data[0] {
		case 0x5D:
			// End of list ']'
			return result, nil
		case 0x73:
			// 's' for string
			str, err := DeserializeString(data)
			if err != nil {
				return nil, err
			}
			result = append(result, *str)
			data = data[len(*str)+2:] // Skip string and delimiter
		case 0x62:
			// 'b' for binary
			bin, err := DeserializeBinary(data)
			if err != nil {
				return nil, err
			}
			result = append(result, bin)
			data = data[len(bin)+5:] // Skip binary data and size
		case 0x5B:
			// Nested list '['
			data = data[1:] // Skip opening bracket
			nestedList, err := DeserializeDataList(data)
			if err != nil {
				return nil, err
			}
			result = append(result, nestedList)
			serializedNestedList, _ := SerializeDataList(nestedList)
			// Adjust the data slice by the size of the serialized nested list
			data = data[len(serializedNestedList)+1:]
		default:
			return nil, errors.New("corrupted buffer or unknown type delimiter")
		}
	}

	return result, nil
}

func GetStringListItem(data []interface{}, index int) (string, error) {
	if len(data) <= index {
		return "", errors.New("index out of range")
	}

	if str, ok := data[index].(string); ok {
		return str, nil
	}

	return "", errors.New("item is not a string")
}

func GetBinaryListItem(data []interface{}, index int) ([]byte, error) {
	if len(data) <= index {
		return nil, errors.New("index out of range")
	}

	if bin, ok := data[index].([]byte); ok {
		return bin, nil
	}

	return nil, errors.New("item is not binary")
}

func GetListItem(data []interface{}, index int) ([]interface{}, error) {
	if len(data) <= index {
		return nil, errors.New("index out of range")
	}

	if list, ok := data[index].([]interface{}); ok {
		return list, nil
	}

	return nil, errors.New("item is not a list")
}
