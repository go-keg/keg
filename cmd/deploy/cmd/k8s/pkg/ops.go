package pkg

import "encoding/json"

type OpItem struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value"`
}

type Ops []OpItem

func (o Ops) JSON() []byte {
	marshal, err := json.Marshal(o)
	if err != nil {
		return nil
	}
	return marshal
}
