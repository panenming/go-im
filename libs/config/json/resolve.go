package json

import (
	"bytes"
	"encoding/json"
)

func Unmarshal(f []byte, result interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(f))
	decoder.UseNumber() // 此处可以保证bigint的精度
	return decoder.Decode(result)
}
