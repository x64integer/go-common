package domain

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
