package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// decodeFromReader is helper function to decode entity from io.Reader
func decodeFromReader(entity interface{}, body io.Reader) error {
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&entity); err != nil {
		return errors.New(fmt.Sprint("failed to decode body from io.Reader, : ", entity, err.Error()))
	}

	return nil
}

// toBytes is helper function to marshal content to []byte
func toBytes(content interface{}) []byte {
	b, err := json.Marshal(content)
	if err != nil {
		return []byte(fmt.Sprintf("failed to marshal Response: %v, err: %s", content, err.Error()))
	}

	return b
}
