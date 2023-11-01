package util

import "encoding/json"

func MarshalPeer(str []string) ([]byte, error) {
	data, err := json.Marshal(str)
	if err != nil {
		return nil, err
	}
	return data, nil
}
