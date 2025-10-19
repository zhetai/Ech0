package util

import "encoding/json"

// JSONMarshal JSON序列化
func JSONMarshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// JSONUnmarshal JSON反序列化
func JSONUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
