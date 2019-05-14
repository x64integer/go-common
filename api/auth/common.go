package auth

import (
	"encoding/json"
	"fmt"
)

func toBytes(content interface{}) []byte {
	b, err := json.Marshal(content)
	if err != nil {
		return []byte(fmt.Sprintf("failed to marshal Response: %v, err: %s", content, err.Error()))
	}

	return b
}
