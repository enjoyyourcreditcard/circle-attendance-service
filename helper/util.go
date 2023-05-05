package helper

import (
	"circle/domain"
	"encoding/json"
)

func MarshalJSON(input []domain.AssignmentResp) ([]byte, error) {
	if input == nil {
		return json.Marshal(make([]string, 0))
	}
	return json.Marshal(input)
}
